package data

import (
	"fmt"
	"krathub/internal/conf"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NewDiscovery 根据配置创建服务发现客户端
func NewDiscovery(cfg *conf.Discovery) registry.Discovery {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Discovery.(type) {
	case *conf.Discovery_Consul_:
		return NewConsulDiscovery(c.Consul)
	case *conf.Discovery_Etcd_:
		// TODO: 实现 Etcd 服务发现
		return nil
	case *conf.Discovery_Nacos_:
		return NewNacosDiscovery(c.Nacos)
	default:
		return nil
	}
}

// NewConsulDiscovery 创建 Consul 服务发现客户端
func NewConsulDiscovery(c *conf.Discovery_Consul) registry.Discovery {
	if c == nil {
		return nil
	}
	// 创建 Consul 客户端配置
	consulConfig := consulApi.DefaultConfig()
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
	client, err := consulApi.NewClient(consulConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create consul client: %v", err))
	}
	r := consul.New(client)
	return r
}

// NewNacosDiscovery 创建 Nacos 服务发现客户端
func NewNacosDiscovery(c *conf.Discovery_Nacos) registry.Discovery {
	if c == nil {
		return nil
	}

	// 创建 Nacos 客户端配置
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Addr, c.Port),
	}

	cc := constant.ClientConfig{
		NamespaceId:         c.Namespace,
		Username:            c.Username,
		Password:            c.Password,
		TimeoutMs:           uint64(c.Timeout.GetSeconds() * 1000),
		NotLoadCacheAtStart: true,
		LogLevel:            "warn",
		LogDir:              "../../logs",
		CacheDir:            "../../cache",
	}

	// 创建命名客户端
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create nacos client: %v", err))
	}

	// 创建group参数，如果未设置则使用默认值
	group := c.Group
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	// 创建 Nacos 服务发现
	r := nacos.New(client, nacos.WithGroup(group))
	return r
}
