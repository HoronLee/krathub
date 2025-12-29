package main

import (
	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	cC "github.com/horonlee/krathub/pkg/configCenter"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

// loadConfig 加载配置，支持从本地文件和配置中心加载
func loadConfig() (*conf.Bootstrap, config.Config, error) {
	// 准备所有配置源
	sources := []config.Source{
		env.NewSource("KRATHUB_"),
		file.NewSource(flagconf),
	}

	// 先加载基本配置以检查是否需要配置中心
	tempConfig := config.New(
		config.WithSource(sources...),
		config.WithResolveActualTypes(true),
	)
	if err := tempConfig.Load(); err != nil {
		return nil, nil, err
	}

	var bc conf.Bootstrap
	if err := tempConfig.Scan(&bc); err != nil {
		return nil, nil, err
	}

	// 检查是否配置了远程配置中心
	if configCfg := bc.Config; configCfg != nil {
		switch cT := configCfg.Config.(type) {
		case *conf.Config_Nacos:
			sources = append(sources, cC.NewNacosConfigSource(cT.Nacos))
		case *conf.Config_Consul:
			sources = append(sources, cC.NewConsulConfigSource(cT.Consul))
		case *conf.Config_Etcd:
			etcdSource := cC.NewEtcdConfigSource(cT.Etcd)
			if etcdSource != nil {
				sources = append(sources, etcdSource)
			}
		}
	}

	// 关闭临时配置
	tempConfig.Close()

	// 创建最终配置对象，包含所有配置源
	c := config.New(
		config.WithSource(sources...),
		config.WithResolveActualTypes(true),
	)

	// 加载配置
	if err := c.Load(); err != nil {
		return nil, nil, err
	}

	// 扫描配置到结构体
	if err := c.Scan(&bc); err != nil {
		return nil, nil, err
	}

	return &bc, c, nil
}
