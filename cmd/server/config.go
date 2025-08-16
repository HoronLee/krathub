package main

import (
	"krathub/internal/conf"
	cC "krathub/pkg/configCenter"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

// loadConfig 加载配置，支持从本地文件和配置中心加载
func loadConfig() (*conf.Bootstrap, config.Config, error) {
	// 创建基本配置源
	c := config.New(
		config.WithSource(
			env.NewSource("KRATHUB_"),
			file.NewSource(flagconf),
		),
	)

	// 加载基本配置
	if err := c.Load(); err != nil {
		return nil, nil, err
	}

	// 扫描基本配置到结构体
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return nil, nil, err
	}

	// 检查是否配置了远程配置中心
	if configCfg := bc.Config; configCfg != nil {
		var configSource config.Source

		switch configType := configCfg.Config.(type) {
		case *conf.Config_Nacos_:
			configSource = cC.NewNacosConfigSource(&cC.NacosConfig{
				Addr:      configType.Nacos.Addr,
				Port:      configType.Nacos.Port,
				Namespace: configType.Nacos.Namespace,
				Username:  configType.Nacos.Username,
				Password:  configType.Nacos.Password,
				Group:     configType.Nacos.Group,
				DataId:    configType.Nacos.DataId,
				Timeout:   configType.Nacos.Timeout,
			})
		case *conf.Config_Consul_:
			configSource = cC.NewConsulConfigSource(&cC.ConsulConfig{
				Addr:       configType.Consul.Addr,
				Scheme:     configType.Consul.Scheme,
				Token:      configType.Consul.Token,
				Datacenter: configType.Consul.Datacenter,
				Key:        configType.Consul.Key,
				Timeout:    configType.Consul.Timeout,
			})
		case *conf.Config_Etcd_:
			configSource = cC.NewEtcdConfigSource(&cC.EtcdConfig{
				Endpoints: configType.Etcd.Endpoints,
				Username:  configType.Etcd.Username,
				Password:  configType.Etcd.Password,
				Key:       configType.Etcd.Key,
				Timeout:   configType.Etcd.Timeout,
			})
		}

		if configSource != nil {
			// 创建新的配置对象，包含远程配置源
			newConfig := config.New(
				config.WithSource(
					env.NewSource("KRATHUB_"),
					file.NewSource(flagconf),
					configSource,
				),
			)

			// 替换配置对象
			c.Close()
			c = newConfig

			// 重新加载配置
			if err := c.Load(); err != nil {
				return nil, nil, err
			}

			// 重新扫描配置到结构体
			if err := c.Scan(&bc); err != nil {
				return nil, nil, err
			}
		}
	}

	return &bc, c, nil
}
