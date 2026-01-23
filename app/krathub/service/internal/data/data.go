package data

import (
	"errors"
	"strings"
	"time"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	dao "github.com/horonlee/krathub/app/krathub/service/internal/data/dao"
	"github.com/horonlee/krathub/pkg/client"
	pkglogger "github.com/horonlee/krathub/pkg/logger"
	"github.com/horonlee/krathub/pkg/redis"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery, NewDB, NewRedis, NewData, NewAuthRepo, NewUserRepo, NewTestRepo)

// Data .
type Data struct {
	query  *dao.Query
	log    *log.Helper
	client client.Client
	redis  *redis.Client
}

// NewData .
func NewData(db *gorm.DB, c *conf.Data, logger log.Logger, client client.Client, redisClient *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	dao.SetDefault(db)
	return &Data{
		query:  dao.Q,
		log:    log.NewHelper(pkglogger.WithModule(logger, "data/data/krathub-service")),
		client: client,
		redis:  redisClient,
	}, cleanup, nil
}

func NewDB(cfg *conf.Data, l log.Logger) (*gorm.DB, error) {
	gormLogger := l.(*pkglogger.ZapLogger).GetGormLogger("gorm/data/krathub-service")
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.GetSource()), &gorm.Config{
			Logger: gormLogger,
		})
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.GetSource()), &gorm.Config{
			Logger: gormLogger,
		})
	case "postgres", "postgresql":
		return gorm.Open(postgres.Open(cfg.Database.GetSource()), &gorm.Config{
			Logger: gormLogger,
		})
	}
	return nil, errors.New("connect db fail: unsupported db driver")
}

// NewRedis 创建Redis客户端从配置初始化
func NewRedis(cfg *conf.Data, logger log.Logger) (*redis.Client, func(), error) {
	if cfg.Redis == nil {
		return nil, nil, errors.New("redis configuration is required")
	}

	redisConfig := &redis.Config{
		Addr:     cfg.Redis.GetAddr(),
		Username: cfg.Redis.GetUserName(),
		Password: cfg.Redis.GetPassword(),
		DB:       int(cfg.Redis.GetDb()),
	}

	if cfg.Redis.GetReadTimeout() != nil {
		redisConfig.ReadTimeout = cfg.Redis.GetReadTimeout().AsDuration()
	} else {
		redisConfig.ReadTimeout = 3 * time.Second
	}

	if cfg.Redis.GetWriteTimeout() != nil {
		redisConfig.WriteTimeout = cfg.Redis.GetWriteTimeout().AsDuration()
	} else {
		redisConfig.WriteTimeout = 3 * time.Second
	}

	return redis.NewClient(redisConfig, pkglogger.WithModule(logger, "redis/data/krathub-service"))
}
