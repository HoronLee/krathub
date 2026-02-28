<!-- Parent: ../AGENTS.md -->
# Redis 客户端封装 (pkg/redis)

**最后更新时间**: 2026-02-09

## 模块目的
对 `github.com/redis/go-redis/v9` 进行二次封装，提供连接测试（Ping）、日志集成及简化的常用操作方法。

## 关键文件
- `redis.go`: 客户端初始化、配置解析及核心操作方法。

## 使用示例

### 1. 从配置初始化
```go
import "github.com/horonlee/servora/pkg/redis"

cfg := &redis.Config{Addr: "localhost:6379", DB: 0}
client, cleanup, err := redis.NewClient(cfg, logger)
defer cleanup()
```

### 2. 基础操作
```go
ctx := context.Background()
err := client.Set(ctx, "key", "value", time.Hour)
val, err := client.Get(ctx, "key")
```

## 测试指南
运行单元测试（需要本地 Redis）：
```bash
go test -v ./pkg/redis/...
```
🤖 Generated with [Claude Code](https://claude.com/claude-code)
