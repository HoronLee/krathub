package server

import (
	"fmt"
	"time"

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
		return registry.NewConsulRegistry(&registry.ConsulConfig{
			Addr:       c.Consul.Addr,
			Scheme:     c.Consul.Scheme,
			Token:      c.Consul.Token,
			Datacenter: c.Consul.Datacenter,
			Timeout:    c.Consul.Timeout,
			Tags:       c.Consul.Tags,
		})
	case *conf.Registry_Etcd:
		var opts []registry.Option
		if c.Etcd.Namespace != "" {
			opts = append(opts, registry.Namespace(c.Etcd.Namespace))
		}
		opts = append(opts, registry.RegisterTTL(15*time.Second), registry.MaxRetry(5))
		registrar, err := registry.NewEtcdRegistry(c.Etcd, opts...)
		if err != nil {
			panic(fmt.Sprintf("failed to create etcd registry: %v", err))
		}
		return registrar
	case *conf.Registry_Nacos:
		return registry.NewNacosRegistry(&registry.NacosConfig{
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
