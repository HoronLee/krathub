package data

import (
	"errors"
	"strings"
	"time"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	dao "github.com/horonlee/krathub/app/krathub/service/internal/data/dao"
	pkglogger "github.com/horonlee/krathub/pkg/logger"
	"github.com/horonlee/krathub/pkg/redis"
	"github.com/horonlee/krathub/pkg/transport/client"

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

	var dialector gorm.Dialector
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		dialector = mysql.Open(cfg.Database.GetSource())
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.GetSource())
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.Database.GetSource())
	default:
		return nil, errors.New("connect db fail: unsupported db driver")
	}

	var db *gorm.DB
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(dialector, &gorm.Config{
			Logger: gormLogger,
		})
		if err == nil {
			return db, nil
		}
		if i < maxRetries-1 {
			delay := time.Duration(1<<uint(i)) * time.Second
			log.NewHelper(l).Warnf("database connection failed (attempt %d/%d), retrying in %v: %v", i+1, maxRetries, delay, err)
			time.Sleep(delay)
		}
	}
	return nil, err
}

func NewRedis(cfg *conf.Data, logger log.Logger) (*redis.Client, func(), error) {
	redisConfig := redis.NewConfigFromProto(cfg.Redis)
	if redisConfig == nil {
		return nil, nil, errors.New("redis configuration is required")
	}

	return redis.NewClient(redisConfig, pkglogger.WithModule(logger, "redis/data/krathub-service"))
}
