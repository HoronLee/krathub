<!-- Parent: ../AGENTS.md -->
# 日志封装 (pkg/logger)

**最后更新时间**: 2026-02-09

## 模块目的
基于 `go.uber.org/zap` 实现 Kratos 的 `log.Logger` 接口。支持分级日志、日志轮转（Lumberjack）、多环境适配及 GORM 日志集成。

## 关键文件
- `log.go`: Zap 适配器实现与初始化。
- `gorm_log.go`: GORM v2 日志接口适配。

## 使用示例

### 1. 初始化 Logger
```go
import "github.com/horonlee/krathub/pkg/logger"

l := logger.NewLogger(&logger.Config{
    Env:      "dev",
    Filename: "logs/app.log",
})
```

### 2. 使用 Option 模式添加字段
```go
// 仅添加 module
helper := log.NewHelper(logger.With(l, logger.WithModule("auth/biz")))
helper.Info("service started")

// 添加多个字段
helper := log.NewHelper(logger.With(l, 
    logger.WithModule("auth/biz"),
    logger.WithField("version", "v1.0.0"),
    logger.WithField("instance", "node-1"),
))
helper.Info("service started")
```

## 测试指南
```bash
go test -v ./pkg/logger/...
```
🤖 Generated with [Claude Code](https://claude.com/claude-code)
