package data

import (
	"errors"
	"krathub/internal/conf"
	"krathub/internal/data/query"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAuthRepo, NewDB)

// Data .
type Data struct {
	query *query.Query
	log   *log.Helper
}

// NewData .
func NewData(db *gorm.DB, c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)
	return &Data{query: query.Q, log: log.NewHelper(logger)}, cleanup, nil
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
