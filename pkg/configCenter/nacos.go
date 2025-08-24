package configCenter

import (
	"fmt"
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

	timeoutMs := uint64(5000)
	if c.Timeout != nil && c.Timeout.AsDuration() > 0 {
		timeoutMs = uint64(c.Timeout.AsDuration().Milliseconds())
	}

	cc := &constant.ClientConfig{
		NamespaceId:         c.Namespace,
		Username:            c.Username,
		Password:            c.Password,
		TimeoutMs:           timeoutMs,
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
		panic(fmt.Sprintf("failed to create nacos client: %v", err))
	}

	// 从Nacos获取dataId
	dataID := c.DataId
	if dataID == "" {
		dataID = "config.yaml"
	}

	// 从Nacos获取group
	group := c.Group
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	return nacosCfg.NewConfigSource(
		client,
		nacosCfg.WithGroup(group),
		nacosCfg.WithDataID(dataID),
	)
}
