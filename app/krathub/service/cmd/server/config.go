package main

import (
	"strings"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	cC "github.com/horonlee/krathub/pkg/governance/configCenter"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
)

func loadConfig() (*conf.Bootstrap, config.Config, error) {
	sources := []config.Source{
		file.NewSource(flagconf),
	}

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

	var configCenterSource config.Source
	if configCfg := bc.Config; configCfg != nil {
		switch cT := configCfg.Config.(type) {
		case *conf.Config_Nacos:
			configCenterSource = cC.NewNacosConfigSource(cT.Nacos)
		case *conf.Config_Consul:
			configCenterSource = cC.NewConsulConfigSource(cT.Consul)
		case *conf.Config_Etcd:
			configCenterSource = cC.NewEtcdConfigSource(cT.Etcd)
		}
	}

	tempConfig.Close()

	finalSources := []config.Source{
		file.NewSource(flagconf),
	}

	if configCenterSource != nil {
		finalSources = append(finalSources, configCenterSource)
	}

	// 从 Name 生成环境变量前缀: "krathub.service" -> "KRATHUB_"
	envPrefix := strings.ToUpper(strings.TrimSuffix(Name, ".service")) + "_"
	finalSources = append(finalSources, env.NewSource(envPrefix))

	c := config.New(
		config.WithSource(finalSources...),
		config.WithResolveActualTypes(true),
	)

	if err := c.Load(); err != nil {
		return nil, nil, err
	}

	if err := c.Scan(&bc); err != nil {
		return nil, nil, err
	}

	return &bc, c, nil
}
