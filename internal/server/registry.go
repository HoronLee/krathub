package server

import (
	"github.com/horonlee/krathub/internal/conf"
	"github.com/horonlee/krathub/pkg/registry"

	kr "github.com/go-kratos/kratos/v2/registry"
)

// NewRegistrar 根据配置创建注册中心客户端
func NewRegistrar(cfg *conf.Registry) kr.Registrar {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Registry.(type) {
	case *conf.Registry_Consul:
		return registry.NewConsulRegistrar(&registry.ConsulConfig{
			Addr:       c.Consul.Addr,
			Scheme:     c.Consul.Scheme,
			Token:      c.Consul.Token,
			Datacenter: c.Consul.Datacenter,
			Timeout:    c.Consul.Timeout,
			Tags:       c.Consul.Tags,
		})
	case *conf.Registry_Etcd:
		// TODO: 实现 Etcd 注册中心
		return nil
	case *conf.Registry_Nacos:
		return registry.NewNacosRegistrar(&registry.NacosConfig{
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
