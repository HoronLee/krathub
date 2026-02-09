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

### 2. 添加模块信息
```go
helper := log.NewHelper(logger.WithModule(l, "auth/biz"))
helper.Info("service started")
```

## 测试指南
```bash
go test -v ./pkg/logger/...
```
🤖 Generated with [Claude Code](https://claude.com/claude-code)
