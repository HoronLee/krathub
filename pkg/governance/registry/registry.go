package registry

import (
	"fmt"
	"time"

	"github.com/horonlee/micro-forge/api/gen/go/conf/v1"

	"github.com/go-kratos/kratos/v2/registry"
)

func NewRegistrar(cfg *conf.Registry) registry.Registrar {
	if cfg == nil {
		return nil
	}

	switch c := cfg.Registry.(type) {
	case *conf.Registry_Consul:
		return NewConsulRegistry(c.Consul)
	case *conf.Registry_Etcd:
		var opts []Option
		if c.Etcd.Namespace != "" {
			opts = append(opts, Namespace(c.Etcd.Namespace))
		}
		opts = append(opts, RegisterTTL(15*time.Second), MaxRetry(5))
		registrar, err := NewEtcdRegistry(c.Etcd, opts...)
		if err != nil {
			panic(fmt.Sprintf("failed to create etcd registry: %v", err))
		}
		return registrar
	case *conf.Registry_Nacos:
		return NewNacosRegistry(c.Nacos)
	case *conf.Registry_Kubernetes:
		return NewKubernetesRegistry(c.Kubernetes)
	default:
		return nil
	}
}
