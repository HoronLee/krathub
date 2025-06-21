package data

import (
	"context"
	"errors"
	hellov1 "krathub/api/hello/v1"
	"krathub/internal/conf"
	"krathub/internal/data/query"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewHelloGrpcClient, NewData, NewDB, NewAuthRepo, NewUserRepo)

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

// NewHelloGrpcClient 创建一个 gRPC 客户端连接到 hello 服务
func NewHelloGrpcClient(cfg *conf.Data) (hellov1.HelloServiceClient, error) {
	var endpoint string
	for _, c := range cfg.Client.GetGrpc() {
		if c.ServiceName == "hello" { // 这里 hello 为所需要访问的服务名
			endpoint = c.Endpoint
			break
		}
	}
	if endpoint == "" {
		return nil, errors.New("no grpc client config found for service: hello")
	}
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithTimeout(3600*time.Second),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	)
	if err != nil {
		return nil, err
	}
	return hellov1.NewHelloServiceClient(conn), nil
}
