package configCenter

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
	"google.golang.org/protobuf/types/known/durationpb"
)

// ConsulConfig Consul 配置中心配置
type ConsulConfig struct {
	Addr       string
	Scheme     string
	Token      string
	Datacenter string
	Key        string // Consul KV 存储的键名
	Timeout    *durationpb.Duration
}

// NewConsulConfigSource 创建 Consul 配置源
func NewConsulConfigSource(c *ConsulConfig) config.Source {
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

	// 设置超时时间
	if c.Timeout != nil {
		consulConfig.WaitTime = c.Timeout.AsDuration()
	} else {
		consulConfig.WaitTime = 5 * time.Second
	}

	// 创建 Consul 客户端
	client, err := api.NewClient(consulConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create consul client: %v", err))
	}

	// 设置配置键名，默认为 config
	key := "config"
	if c.Key != "" {
		key = c.Key
	}

	source, err := consul.New(client, consul.WithPath(key))
	if err != nil {
		panic(fmt.Sprintf("failed to create consul config source: %v", err))
	}

	return source
}
