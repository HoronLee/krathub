package data

import (
	"fmt"
	"time"

	"github.com/horonlee/micro-forge/api/gen/go/conf/v1"
	"github.com/horonlee/micro-forge/pkg/governance/registry"

	kratosRegistry "github.com/go-kratos/kratos/v2/registry"
)

// NewDiscovery 根据配置创建服务发现客户端
func NewDiscovery(cfg *conf.Discovery) kratosRegistry.Discovery {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Discovery.(type) {
	case *conf.Discovery_Consul:
		return registry.NewConsulDiscovery(c.Consul)
	case *conf.Discovery_Etcd:
		var opts []registry.Option
		if c.Etcd.Namespace != "" {
			opts = append(opts, registry.Namespace(c.Etcd.Namespace))
		}
		opts = append(opts, registry.RegisterTTL(15*time.Second), registry.MaxRetry(5))
		discovery, err := registry.NewEtcdDiscovery(c.Etcd, opts...)
		if err != nil {
			panic(fmt.Sprintf("failed to create etcd discovery: %v", err))
		}
		return discovery
	case *conf.Discovery_Nacos:
		return registry.NewNacosDiscovery(c.Nacos)
	case *conf.Discovery_Kubernetes:
		return registry.NewKubernetesDiscovery(c.Kubernetes)
	default:
		return nil
	}
}
