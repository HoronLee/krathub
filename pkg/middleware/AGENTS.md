<!-- Parent: ../AGENTS.md -->
# 中间件 (pkg/middleware)

**最后更新时间**: 2026-02-09

## 模块目的
提供通用的 HTTP/gRPC 中间件，包括跨域资源共享 (CORS) 处理和基于 IP 的白名单访问控制。

## 关键文件
- `cors/cors.go`: 处理跨域请求的中间件，支持通配符和自定义配置。
- `whitelist.go`: 提供简单的 IP 白名单过滤。

## 使用示例

### CORS 中间件集成 (HTTP)
```go
import "github.com/horonlee/servora/pkg/middleware/cors"

httpSrv := http.NewServer(
    http.Middleware(
        // 将 CORS 配置传入
        cors.Middleware(conf.CORS),
    ),
)
```

## 测试指南
```bash
go test -v ./pkg/middleware/cors/...
```
🤖 Generated with [Claude Code](https://claude.com/claude-code)
