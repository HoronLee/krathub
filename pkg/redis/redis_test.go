package redis

import (
	"context"
	"testing"
	"time"

	"github.com/horonlee/krathub/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis 设置测试用的Redis客户端
func setupTestRedis(t *testing.T) (*Client, func()) {
	cfg := &Config{
		Addr:         "localhost:6379",
		Password:     "114514",
		DB:           1, // 使用测试数据库
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	testLogger := logger.NewLogger(&logger.Config{
		Env:   "test",
		Level: 1, // info level
	})
	client, cleanup, err := NewClient(cfg, testLogger)
	require.NoError(t, err, "Failed to create Redis client for testing")

	// 清理测试数据库
	ctx := context.Background()
	client.rdb.FlushDB(ctx)

	return client, cleanup
}

func TestNewClient_Success(t *testing.T) {
	cfg := &Config{
		Addr:         "localhost:6379",
		Password:     "114514",
		DB:           1,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	testLogger := logger.NewLogger(&logger.Config{
		Env:   "test",
		Level: 1, // info level
	})
	client, cleanup, err := NewClient(cfg, testLogger)

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, cleanup)

	if cleanup != nil {
		cleanup()
	}
}

func TestNewClient_ConnectionFailure(t *testing.T) {
	cfg := &Config{
		Addr:         "localhost:9999", // 无效端口
		Password:     "wrong",
		DB:           0,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	testLogger := logger.NewLogger(&logger.Config{
		Env:   "test",
		Level: 1, // info level
	})
	client, cleanup, err := NewClient(cfg, testLogger)

	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Nil(t, cleanup)
}

func TestClient_Set_Get(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test_key"
	value := "test_value"

	// 测试Set操作
	err := client.Set(ctx, key, value, time.Hour)
	assert.NoError(t, err)

	// 测试Get操作
	result, err := client.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestClient_Get_NotFound(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "nonexistent_key"

	result, err := client.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Empty(t, result)
}

func TestClient_Del(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test_key"
	value := "test_value"

	// 先设置一个键
	err := client.Set(ctx, key, value, time.Hour)
	require.NoError(t, err)

	// 验证键存在
	result, err := client.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value, result)

	// 删除键
	err = client.Del(ctx, key)
	assert.NoError(t, err)

	// 验证键已被删除
	_, err = client.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestClient_Has(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test_key"
	value := "test_value"

	// 测试不存在的键
	exists := client.Has(ctx, key)
	assert.False(t, exists)

	// 设置键
	err := client.Set(ctx, key, value, time.Hour)
	require.NoError(t, err)

	// 测试存在的键
	exists = client.Has(ctx, key)
	assert.True(t, exists)
}

func TestClient_Keys(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	// 设置多个键
	keys := []string{"test:key1", "test:key2", "other:key3"}
	for _, key := range keys {
		err := client.Set(ctx, key, "value", time.Hour)
		require.NoError(t, err)
	}

	// 测试模式匹配
	result, err := client.Keys(ctx, "test:*")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "test:key1")
	assert.Contains(t, result, "test:key2")
	assert.NotContains(t, result, "other:key3")
}

func TestClient_SAdd_SMembers(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test_set"
	members := []interface{}{"member1", "member2", "member3"}

	// 测试SAdd操作
	err := client.SAdd(ctx, key, members...)
	assert.NoError(t, err)

	// 测试SMembers操作
	result, err := client.SMembers(ctx, key)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Contains(t, result, "member1")
	assert.Contains(t, result, "member2")
	assert.Contains(t, result, "member3")
}

func TestClient_Expire(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test_expire_key"
	value := "test_value"

	// 设置键
	err := client.Set(ctx, key, value, time.Hour)
	require.NoError(t, err)

	// 设置过期时间
	err = client.Expire(ctx, key, 2*time.Second)
	assert.NoError(t, err)

	// 验证键仍然存在
	exists := client.Has(ctx, key)
	assert.True(t, exists)

	// 等待过期
	time.Sleep(3 * time.Second)

	// 验证键已过期
	exists = client.Has(ctx, key)
	assert.False(t, exists)
}

func TestClient_SetWithExpiration(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test_expire_key"
	value := "test_value"

	// 设置带过期时间的键
	err := client.Set(ctx, key, value, 2*time.Second)
	assert.NoError(t, err)

	// 验证键存在
	result, err := client.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// 等待过期
	time.Sleep(3 * time.Second)

	// 验证键已过期
	_, err = client.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestClient_DelMultipleKeys(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	keys := []string{"key1", "key2", "key3"}

	// 设置多个键
	for _, key := range keys {
		err := client.Set(ctx, key, "value", time.Hour)
		require.NoError(t, err)
	}

	// 删除多个键
	err := client.Del(ctx, keys...)
	assert.NoError(t, err)

	// 验证所有键都已删除
	for _, key := range keys {
		exists := client.Has(ctx, key)
		assert.False(t, exists)
	}
}

func TestClient_SAddDuplicateMembers(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	key := "test_set"

	// 添加重复成员
	err := client.SAdd(ctx, key, "member1", "member2", "member1")
	assert.NoError(t, err)

	// 验证集合只包含唯一成员
	result, err := client.SMembers(ctx, key)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "member1")
	assert.Contains(t, result, "member2")
}
