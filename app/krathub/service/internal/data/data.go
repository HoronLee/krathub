package data

import (
	"errors"
	"strings"
	"time"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/horonlee/krathub/app/krathub/service/internal/client"
	dao "github.com/horonlee/krathub/app/krathub/service/internal/data/dao"
	"github.com/horonlee/krathub/pkg/logger"
	"github.com/horonlee/krathub/pkg/redis"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery, NewDB, NewRedis, NewData, NewAuthRepo, NewUserRepo)

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
	// 为GEN生成的query代码设置数据库连接对象
	dao.SetDefault(db)
	return &Data{
		query:  dao.Q,
		log:    log.NewHelper(logger),
		client: client,
		redis:  redisClient,
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

	// 设置超时时间
	if cfg.Redis.GetReadTimeout() != nil {
		redisConfig.ReadTimeout = cfg.Redis.GetReadTimeout().AsDuration()
	} else {
		redisConfig.ReadTimeout = 3 * time.Second // 默认读超时
	}

	if cfg.Redis.GetWriteTimeout() != nil {
		redisConfig.WriteTimeout = cfg.Redis.GetWriteTimeout().AsDuration()
	} else {
		redisConfig.WriteTimeout = 3 * time.Second // 默认写超时
	}

	return redis.NewClient(redisConfig, logger)
}
