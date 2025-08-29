package data

import (
	"github.com/horonlee/krathub/internal/conf"
	"github.com/horonlee/krathub/pkg/discovery"

	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
)

// NewDiscovery 根据配置创建服务发现客户端
func NewDiscovery(cfg *conf.Discovery) kratosRegistry.Discovery {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Discovery.(type) {
	case *conf.Discovery_Consul_:
		return discovery.NewConsulDiscovery(&discovery.ConsulConfig{
			Addr:       c.Consul.Addr,
			Scheme:     c.Consul.Scheme,
			Token:      c.Consul.Token,
			Datacenter: c.Consul.Datacenter,
			Timeout:    c.Consul.Timeout,
		})
	case *conf.Discovery_Etcd_:
		// TODO: 实现 Etcd 服务发现
		return nil
	case *conf.Discovery_Nacos_:
		return discovery.NewNacosDiscovery(&discovery.NacosConfig{
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
