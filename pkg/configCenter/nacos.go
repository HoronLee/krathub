package configCenter

import (
	nacosCfg "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/protobuf/types/known/durationpb"
)

// NacosConfig Nacos 配置中心配置
type NacosConfig struct {
	Addr      string
	Port      uint64
	Namespace string
	Username  string
	Password  string
	Group     string
	DataId    string
	Timeout   *durationpb.Duration
}

// NewNacosConfigSource 创建 Nacos 配置源
func NewNacosConfigSource(c *NacosConfig) config.Source {
	if c == nil {
		return nil
	}

	// 创建Nacos配置客户端
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Addr, c.Port),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         c.Namespace,
		Username:            c.Username,
		Password:            c.Password,
		TimeoutMs:           uint64(c.Timeout.GetSeconds() * 1000),
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	// 从Nacos获取dataId
	dataID := "config.yaml"
	if c.DataId != "" {
		dataID = c.DataId
	}

	// 从Nacos获取group
	group := "DEFAULT_GROUP"
	if c.Group != "" {
		group = c.Group
	}

	return nacosCfg.NewConfigSource(
		client,
		nacosCfg.WithGroup(group),
		nacosCfg.WithDataID(dataID),
	)
}
