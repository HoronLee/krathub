package configCenter

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/types/known/durationpb"
)

// EtcdConfig Etcd 配置中心配置
type EtcdConfig struct {
	Endpoints []string
	Username  string
	Password  string
	Key       string // Etcd 存储的键名
	Timeout   *durationpb.Duration
}

// etcdSource 实现 config.Source 接口的 Etcd 配置源
type etcdSource struct {
	client   *clientv3.Client
	key      string
	mu       sync.RWMutex
	watchers []config.Watcher
}

// NewEtcdConfigSource 创建 Etcd 配置源
func NewEtcdConfigSource(c *EtcdConfig) config.Source {
	if c == nil {
		return nil
	}

	// 创建 Etcd 客户端配置
	etcdConfig := clientv3.Config{
		Endpoints: c.Endpoints,
		Username:  c.Username,
		Password:  c.Password,
	}

	// 设置超时时间
	if c.Timeout != nil {
		etcdConfig.DialTimeout = c.Timeout.AsDuration()
	} else {
		etcdConfig.DialTimeout = 5 * time.Second
	}

	// 创建 Etcd 客户端
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create etcd client: %v", err))
	}

	// 设置配置键名，默认为 /config
	key := "/config"
	if c.Key != "" {
		key = c.Key
	}

	return &etcdSource{
		client: client,
		key:    key,
	}
}

// Load 实现 config.Source 接口
func (s *etcdSource) Load() ([]*config.KeyValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.client.Get(ctx, s.key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var kvs []*config.KeyValue
	for _, kv := range resp.Kvs {
		kvs = append(kvs, &config.KeyValue{
			Key:    string(kv.Key),
			Value:  kv.Value,
			Format: s.getFormat(string(kv.Key)),
		})
	}

	return kvs, nil
}

// Watch 实现 config.Source 接口
func (s *etcdSource) Watch() (config.Watcher, error) {
	watcher := &etcdWatcher{
		source: s,
		ctx:    context.Background(),
		cancel: func() {},
	}

	s.mu.Lock()
	s.watchers = append(s.watchers, watcher)
	s.mu.Unlock()

	return watcher, nil
}

// getFormat 根据文件扩展名推断配置格式
func (s *etcdSource) getFormat(key string) string {
	ext := filepath.Ext(key)
	switch ext {
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".toml":
		return "toml"
	default:
		return "yaml" // 默认格式
	}
}

// etcdWatcher 实现 config.Watcher 接口的 Etcd 监听器
type etcdWatcher struct {
	source *etcdSource
	ctx    context.Context
	cancel context.CancelFunc
}

// Next 实现 config.Watcher 接口
func (w *etcdWatcher) Next() ([]*config.KeyValue, error) {
	w.ctx, w.cancel = context.WithCancel(context.Background())

	watchCh := w.source.client.Watch(w.ctx, w.source.key, clientv3.WithPrefix())

	for watchResp := range watchCh {
		if watchResp.Err() != nil {
			return nil, watchResp.Err()
		}

		var kvs []*config.KeyValue
		for _, event := range watchResp.Events {
			kvs = append(kvs, &config.KeyValue{
				Key:    string(event.Kv.Key),
				Value:  event.Kv.Value,
				Format: w.source.getFormat(string(event.Kv.Key)),
			})
		}

		if len(kvs) > 0 {
			return kvs, nil
		}
	}

	return nil, fmt.Errorf("watch channel closed")
}

// Stop 实现 config.Watcher 接口
func (w *etcdWatcher) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}
