package server

import (
	"fmt"
	"krathub/internal/conf"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

// NewRegistrar 根据配置创建注册中心客户端
func NewRegistrar(cfg *conf.Registry) registry.Registrar {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Registry.(type) {
	case *conf.Registry_Consul_:
		return NewConsulRegistry(c.Consul)
	case *conf.Registry_Etcd_:
		// TODO: 实现 Etcd 注册中心
		return nil
	case *conf.Registry_Nacos_:
		// TODO: 实现 Nacos 注册中心
		return nil
	default:
		return nil
	}
}

// NewConsulRegistry 创建 Consul 注册中心客户端
func NewConsulRegistry(c *conf.Registry_Consul) registry.Registrar {
	if c == nil {
		return nil
	}
	// 创建 Consul 客户端配置
	cConfig := api.DefaultConfig()
	cConfig.Address = c.Addr
	if c.Scheme != "" {
		cConfig.Scheme = c.Scheme
	}
	if c.Token != "" {
		cConfig.Token = c.Token
	}
	if c.Datacenter != "" {
		cConfig.Datacenter = c.Datacenter
	}
	if c.Timeout != nil {
		cConfig.WaitTime = c.Timeout.AsDuration()
	} else {
		cConfig.WaitTime = 5 * time.Second // 默认超时时间
	}
	// 创建 Consul 客户端
	client, err := api.NewClient(cConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create consul client: %v", err))
	}
	// 创建 Consul 注册中心，使用健康检查
	r := consul.New(client, consul.WithHealthCheck(true))
	return r
}
