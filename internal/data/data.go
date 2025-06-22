package data

import (
	"context"
	"errors"
	"fmt"
	hellov1 "krathub/api/hello/v1"
	"krathub/internal/conf"
	"krathub/internal/data/query"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery, NewHelloGrpcClient, NewData, NewDB, NewAuthRepo, NewUserRepo)

// Data .
type Data struct {
	query *query.Query
	log   *log.Helper
	hc    hellov1.HelloServiceClient
}

// NewData .
func NewData(db *gorm.DB, c *conf.Data, logger log.Logger, hc hellov1.HelloServiceClient) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)
	return &Data{query: query.Q, log: log.NewHelper(logger), hc: hc}, cleanup, nil
}

func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.GetSource()))
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.GetSource()))
	}
	return nil, errors.New("connect db fail: unsupported db driver")
}

// NewHelloGrpcClient 创建 Hello 服务的 gRPC 客户端
func NewHelloGrpcClient(cfg *conf.Data, d registry.Discovery) (hellov1.HelloServiceClient, error) {
	var serviceName string
	for _, c := range cfg.Client.GetGrpc() {
		if c.ServiceName == "hello" { // 这里为所需要访问的服务名
			serviceName = c.ServiceName
			break
		}
	}
	if serviceName == "" {
		return nil, errors.New("no grpc client config found for hello service")
	}
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///"+serviceName),
		grpc.WithDiscovery(d),
		grpc.WithTimeout(3600*time.Second),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	return hellov1.NewHelloServiceClient(conn), nil
}

// NewDiscovery 根据配置创建服务发现客户端
func NewDiscovery(cfg *conf.Registry) registry.Discovery {
	if cfg == nil {
		return nil
	}
	switch c := cfg.Registry.(type) {
	case *conf.Registry_Consul_:
		return NewConsulDiscovery(c.Consul)
	case *conf.Registry_Etcd_:
		// TODO: 实现 Etcd 服务发现
		return nil
	case *conf.Registry_Nacos_:
		// TODO: 实现 Nacos 服务发现
		return nil
	default:
		return nil
	}
}

// NewConsulDiscovery 创建 Consul 服务发现客户端
func NewConsulDiscovery(c *conf.Registry_Consul) registry.Discovery {
	if c == nil {
		return nil
	}
	// 创建 Consul 客户端配置
	consulConfig := api.DefaultConfig()
	consulConfig.Address = c.Addr
	if c.Scheme != "" {
		consulConfig.Scheme = c.Scheme
	}
	if c.Token != "" {
		consulConfig.Token = c.Token
	}
	if c.Datacenter != "" {
		consulConfig.Datacenter = c.Datacenter
	}
	if c.Timeout != nil {
		consulConfig.WaitTime = c.Timeout.AsDuration()
	} else {
		consulConfig.WaitTime = 5 * time.Second // 默认超时时间
	}
	// 创建 Consul 客户端
	client, err := api.NewClient(consulConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create consul client: %v", err))
	}
	r := consul.New(client)
	return r
}
