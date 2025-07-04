package data

import (
	"context"
	"time"

	"krathub/internal/conf"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis 服务
type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

// NewRedisClient 依赖注入形式的 RedisClient 构造函数
func NewRedisClient(cfg *conf.Data) *RedisClient {
	rds := &RedisClient{
		Context: context.Background(),
	}
	rds.Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.UserName,
		Password: cfg.Redis.Password,
		DB:       int(cfg.Redis.Db),
	})
	if err := rds.Ping(); err != nil {
		return nil
	}
	return rds
}

// Ping 用以测试 redis 连接是否正常
func (rds *RedisClient) Ping() error {
	_, err := rds.Client.Ping(rds.Context).Result()
	return err
}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds *RedisClient) Set(key string, value any, expiration time.Duration) error {
	return rds.Client.Set(rds.Context, key, value, expiration).Err()
}

// Get 获取 key 对应的 value
func (rds *RedisClient) Get(key string) (string, error) {
	result, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// Has 判断一个 key 是否存在，内部错误和 redis.Nil 都返回 false
func (rds *RedisClient) Has(key string) (bool, error) {
	_, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Del 删除存储在 redis 里的数据，支持多个 key 传参
func (rds *RedisClient) Del(keys ...string) error {
	return rds.Client.Del(rds.Context, keys...).Err()
}

// FlushDB 清空当前 redis db 里的所有数据
func (rds *RedisClient) FlushDB() error {
	return rds.Client.FlushDB(rds.Context).Err()
}

// Increment 自增
func (rds *RedisClient) Increment(key string, value ...int64) error {
	if len(value) == 0 {
		return rds.Client.Incr(rds.Context, key).Err()
	} else {
		return rds.Client.IncrBy(rds.Context, key, value[0]).Err()
	}
}

// Decrement 自减
func (rds *RedisClient) Decrement(key string, value ...int64) error {
	if len(value) == 0 {
		return rds.Client.Decr(rds.Context, key).Err()
	} else {
		return rds.Client.DecrBy(rds.Context, key, value[0]).Err()
	}
}
