package data

import (
	"errors"
	"fmt"
	"krathub/internal/client"
	"krathub/internal/conf"
	"krathub/internal/data/query"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery, NewData, NewDB, NewAuthRepo, NewUserRepo)

// Data .
type Data struct {
	query         *query.Query
	log           *log.Helper
	clientFactory client.GrpcClientFactory
}

// NewData .
func NewData(db *gorm.DB, c *conf.Data, logger log.Logger, clientFactory client.GrpcClientFactory) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)
	return &Data{
		query:         query.Q,
		log:           log.NewHelper(logger),
		clientFactory: clientFactory,
	}, cleanup, nil
}

func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.GetSource()))
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.GetSource()))
	}
	return nil, errors.New("connect db fail: unsupported db driver")
}

// NewDiscovery 根据配置创建服务发现客户端
func NewDiscovery(cfg *conf.Registry) registry.Discovery {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Registry.(type) {
	case *conf.Registry_Consul_:
		return NewConsulDiscovery(c.Consul)
	case *conf.Registry_Etcd_:
		// TODO: 实现 Etcd 服务发现
		return nil
	case *conf.Registry_Nacos_:
		// TODO: 实现 Nacos 服务发现
		return nil
	default:
		return nil
	}
}

// NewConsulDiscovery 创建 Consul 服务发现客户端
func NewConsulDiscovery(c *conf.Registry_Consul) registry.Discovery {
	if c == nil {
		return nil
	}
	// 创建 Consul 客户端配置
	consulConfig := api.DefaultConfig()
	consulConfig.Address = c.Addr
	if c.Scheme != "" {
		consulConfig.Scheme = c.Scheme
	}
	if c.Token != "" {
		consulConfig.Token = c.Token
	}
	if c.Datacenter != "" {
		consulConfig.Datacenter = c.Datacenter
	}
	if c.Timeout != nil {
		consulConfig.WaitTime = c.Timeout.AsDuration()
	} else {
		consulConfig.WaitTime = 5 * time.Second // 默认超时时间
	}
	// 创建 Consul 客户端
	client, err := api.NewClient(consulConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create consul client: %v", err))
	}
	r := consul.New(client)
	return r
}
