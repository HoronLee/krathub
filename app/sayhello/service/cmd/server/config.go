package main

import (
	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	cC "github.com/horonlee/krathub/pkg/configCenter"

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

	finalSources = append(finalSources, env.NewSource("SAYHELLO_"))

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
