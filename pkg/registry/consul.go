package registry

import (
	"fmt"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
	"google.golang.org/protobuf/types/known/durationpb"
)

// ConsulConfig Consul 注册中心配置
type ConsulConfig struct {
	Addr       string
	Scheme     string
	Token      string
	Datacenter string
	Timeout    *durationpb.Duration
	Tags       []string
}

// NewConsulRegistry 创建 Consul 统一注册中心客户端（支持注册和发现）
func NewConsulRegistry(c *ConsulConfig) registry.Registrar {
	if c == nil {
		return nil
	}

	// 创建 Consul 客户端配置
	consulConfig := api.DefaultConfig()

	// 设置基本配置项，Consul API 内部会处理空值
	consulConfig.Address = c.Addr
	consulConfig.Scheme = c.Scheme
	consulConfig.Token = c.Token
	consulConfig.Datacenter = c.Datacenter

	// 超时时间仍需要设置默认值
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

	// 定义 Consul 注册中心的选项
	opts := []consul.Option{
		consul.WithHealthCheck(true),
	}
	if len(c.Tags) > 0 {
		opts = append(opts, consul.WithTags(c.Tags))
	}

	// 创建 Consul 注册中心
	r := consul.New(client, opts...)
	return r
}

// NewConsulDiscovery 创建 Consul 服务发现客户端
func NewConsulDiscovery(c *ConsulConfig) registry.Discovery {
	if c == nil {
		return nil
	}

	// 创建 Consul 客户端配置
	consulConfig := api.DefaultConfig()

	// 设置基本配置项，Consul API 内部会处理空值
	consulConfig.Address = c.Addr
	consulConfig.Scheme = c.Scheme
	consulConfig.Token = c.Token
	consulConfig.Datacenter = c.Datacenter

	// 超时时间仍需要设置默认值
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

// NewConsulRegistrar 创建 Consul 注册中心客户端
// Deprecated: 使用 NewConsulRegistry 替代
func NewConsulRegistrar(c *ConsulConfig) registry.Registrar {
	return NewConsulRegistry(c)
}
