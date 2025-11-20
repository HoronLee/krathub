package data

import (
	"fmt"
	"time"

	"github.com/horonlee/krathub/internal/conf"
	"github.com/horonlee/krathub/pkg/registry"

	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
)

// NewDiscovery 根据配置创建服务发现客户端
func NewDiscovery(cfg *conf.Discovery) kratosRegistry.Discovery {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Discovery.(type) {
	case *conf.Discovery_Consul:
		return registry.NewConsulDiscovery(&registry.ConsulConfig{
			Addr:       c.Consul.Addr,
			Scheme:     c.Consul.Scheme,
			Token:      c.Consul.Token,
			Datacenter: c.Consul.Datacenter,
			Timeout:    c.Consul.Timeout,
		})
	case *conf.Discovery_Etcd:
		namespace := "/krathub"
		if c.Etcd.Namespace != "" {
			namespace = c.Etcd.Namespace
		}
		discovery, err := registry.NewEtcdDiscovery(c.Etcd,
			registry.Namespace(namespace),
			registry.RegisterTTL(15*time.Second),
			registry.MaxRetry(5),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create etcd discovery: %v", err))
		}
		return discovery
	case *conf.Discovery_Nacos:
		return registry.NewNacosDiscovery(&registry.NacosConfig{
			Addr:      c.Nacos.Addr,
			Port:      c.Nacos.Port,
			Namespace: c.Nacos.Namespace,
			Username:  c.Nacos.Username,
			Password:  c.Nacos.Password,
			Group:     c.Nacos.Group,
			Timeout:   c.Nacos.Timeout,
		})
	default:
		return nil
	}
}
