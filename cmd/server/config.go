package main

import (
	"krathub/internal/conf"

	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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
	if configSrc := bc.Config; configSrc != nil {
		var newConfig config.Config
		var err error

		switch configType := configSrc.Config.(type) {
		case *conf.Config_Nacos_:
			newConfig, err = loadNacosConfig(c, configType.Nacos)
		case *conf.Config_Consul_:
			newConfig, err = loadConsulConfig(c, configType.Consul)
		case *conf.Config_Etcd_:
			newConfig, err = loadEtcdConfig(c, configType.Etcd)
		}

		if err != nil {
			return nil, nil, err
		}

		if newConfig != nil {
			// 替换配置对象
			c.Close()
			c = newConfig

			// 重新扫描配置到结构体
			if err := c.Scan(&bc); err != nil {
				return nil, nil, err
			}
		}
	}

	return &bc, c, nil
}

// loadNacosConfig 从Nacos配置中心加载配置
func loadNacosConfig(c config.Config, nacos *conf.Config_Nacos) (config.Config, error) {
	if nacos == nil {
		return nil, nil
	}

	// 创建Nacos配置客户端
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(nacos.Addr, nacos.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         nacos.Namespace,
		Username:            nacos.Username,
		Password:            nacos.Password,
		TimeoutMs:           uint64(nacos.Timeout.GetSeconds() * 1000),
		NotLoadCacheAtStart: true,
		LogDir:              "../../logs",
		CacheDir:            "../../cache",
		LogLevel:            "debug",
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		return nil, err
	}

	// 从Nacos获取dataId
	dataID := "krathub.yaml"
	if nacos.DataId != "" {
		dataID = nacos.DataId
	}

	// 从Nacos获取group
	group := "DEFAULT_GROUP"
	if nacos.Group != "" {
		group = nacos.Group
	}

	// 创建新配置对象，包含Nacos配置源
	newConfig := config.New(
		config.WithSource(
			env.NewSource("KRATHUB_"),
			file.NewSource(flagconf),
			knacos.NewConfigSource(
				client,
				knacos.WithGroup(group),
				knacos.WithDataID(dataID),
			),
		),
	)

	// 重新加载配置
	if err := newConfig.Load(); err != nil {
		return nil, err
	}

	return newConfig, nil
}

// loadConsulConfig 从Consul配置中心加载配置
func loadConsulConfig(c config.Config, consul *conf.Config_Consul) (config.Config, error) {
	// TODO: 实现Consul配置中心加载逻辑
	return nil, nil
}

// loadEtcdConfig 从Etcd配置中心加载配置
func loadEtcdConfig(c config.Config, etcd *conf.Config_Etcd) (config.Config, error) {
	// TODO: 实现Etcd配置中心加载逻辑
	return nil, nil
}
