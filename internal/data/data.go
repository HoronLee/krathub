package data

import (
	"errors"
	"strings"

	"github.com/horonlee/krathub/internal/client"
	"github.com/horonlee/krathub/internal/conf"
	"github.com/horonlee/krathub/internal/data/query"
	"github.com/horonlee/krathub/pkg/logger"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery, NewDB, NewData, NewAuthRepo, NewUserRepo)

// Data .
type Data struct {
	query  *query.Query
	log    *log.Helper
	client client.Client
}

// NewData .
func NewData(db *gorm.DB, c *conf.Data, logger log.Logger, client client.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)
	return &Data{
		query:  query.Q,
		log:    log.NewHelper(logger),
		client: client,
	}, cleanup, nil
}

func NewDB(cfg *conf.Data, l log.Logger) (*gorm.DB, error) {
	gormLogger := l.(*logger.ZapLogger).GetGormLogger()
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.GetSource()), &gorm.Config{
			Logger: gormLogger,
		})
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.GetSource()), &gorm.Config{
			Logger: gormLogger,
		})
	}
	return nil, errors.New("connect db fail: unsupported db driver")
}
