package data

import (
	"errors"
	"krathub/internal/client"
	"krathub/internal/conf"
	"krathub/internal/data/query"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery, NewRedisClient, NewDB, NewData, NewAuthDBRepo, NewAuthGrpcRepo, NewUserDBRepo)

// Data .
type Data struct {
	query         *query.Query
	log           *log.Helper
	clientFactory client.GrpcClientFactory
	rc            *RedisClient
}

// NewData .
func NewData(db *gorm.DB, c *conf.Data, logger log.Logger, clientFactory client.GrpcClientFactory) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)
	return &Data{
		query:         query.Q,
		log:           log.NewHelper(logger),
		clientFactory: clientFactory,
		rc:            NewRedisClient(c),
	}, cleanup, nil
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
