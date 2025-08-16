package registry

import (
	"fmt"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/protobuf/types/known/durationpb"
)

// NacosConfig Nacos 注册中心配置
type NacosConfig struct {
	Addr      string
	Port      uint64
	Namespace string
	Username  string
	Password  string
	Group     string
	Timeout   *durationpb.Duration
}

// NewNacosRegistrar 创建 Nacos 注册中心客户端
func NewNacosRegistrar(c *NacosConfig) registry.Registrar {
	if c == nil {
		return nil
	}

	// 创建 Nacos 服务端配置
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(c.Addr, c.Port),
	}

	// 创建 Nacos 客户端配置
	cc := constant.ClientConfig{
		NamespaceId:         c.Namespace,
		Username:            c.Username,
		Password:            c.Password,
		TimeoutMs:           uint64(c.Timeout.GetSeconds() * 1000),
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
		LogDir:              "./logs",
		CacheDir:            "./cache",
	}

	// 创建命名客户端
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create nacos client: %v", err))
	}

	// 创建group参数，如果未设置则使用默认值
	group := c.Group
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	// 创建 Nacos 注册中心
	r := nacos.New(client, nacos.WithGroup(group))
	return r
}
