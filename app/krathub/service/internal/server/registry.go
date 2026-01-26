package server

import (
	"fmt"
	"time"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/horonlee/krathub/pkg/governance/registry"

	kr "github.com/go-kratos/kratos/v2/registry"
)

// NewRegistrar 根据配置创建注册中心客户端
func NewRegistrar(cfg *conf.Registry) kr.Registrar {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Registry.(type) {
	case *conf.Registry_Consul:
		return registry.NewConsulRegistry(c.Consul)
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
		return registry.NewNacosRegistry(c.Nacos)
	default:
		return nil
	}
}
