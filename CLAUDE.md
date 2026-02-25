# CLAUDE.md - Krathub Development Guide

> **重要**: 永远使用中文回复

Instructions for AI assistants working in this project.

<!-- OPENSPEC:START -->
## OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

## Project Overview

Krathub is a Kratos v2 (Go) microservices project using Buf (Protobuf), Wire (DI), GORM + GORM GEN (ORM), and Vue 3 + Vite (frontend at `app/krathub/service/web/`).

## Build / Lint / Test Commands

### Root-Level
```bash
make init          # Install all dev tools
make gen           # Generate all code (ent + wire + api + openapi)
make build         # Build all services
make test          # Run all Go tests
make lint          # Run golangci-lint
```

### Service-Level (`app/{service}/service/`)
```bash
make run           # Run service
make wire          # Generate wire code
make gen.dao       # Generate GORM GEN PO/DAO
```

### Frontend (`app/krathub/service/web/`)
```bash
bun install && bun dev      # Dev server
bun test:unit               # Vitest unit tests
bun test:e2e                # Playwright E2E tests
bun lint                    # ESLint
```

### Running Single Tests
```bash
# Go
go test -v -run TestFunctionName ./path/to/package
go test -v ./pkg/redis/...

# Frontend
bun test:unit src/__tests__/example.spec.ts
bun test:e2e e2e/example.spec.ts --project=chromium
```

## Project Structure

```
krathub/
├── api/
│   ├── protos/           # Proto 定义（i_*.proto=HTTP, 其他=gRPC）
│   └── gen/go/           # 自动生成的代码（勿修改）
├── app/
│   └── {service}/service/
│       ├── cmd/          # 服务入口
│       ├── internal/     # DDD 三层架构
│       │   ├── biz/      # 业务逻辑层
│       │   ├── data/     # 数据访问层
│       │   └── service/  # 接口实现层
│       └── web/          # Vue 3 前端（仅 krathub）
├── pkg/                  # 共享库（jwt, redis, logger, middleware, governance）
└── openspec/             # OpenSpec 规范文档（可选）
```

## Code Style Guidelines

### Go Imports
```go
import (
    "context"                                              // 1. stdlib

    "github.com/go-kratos/kratos/v2/log"                   // 2. third-party

    authv1 "github.com/horonlee/krathub/api/gen/go/auth/service/v1"  // 3. project
)
```

### Naming
- Interfaces: `UserRepo`, `AuthRepo`
- Constructors: `NewUserUsecase`, `NewUserRepo`
- Private types: lowercase (`userRepo`)

### Error Handling
Use Kratos error types from generated protos:
```go
return userv1.ErrorUserNotFound("user not found: %v", err)
return authv1.ErrorUnauthorized("user not authenticated")
```

### DDD 分层架构
- **Service 层**: API 接口实现、参数验证、DTO 转换
- **Biz 层**: 业务逻辑、UseCase、领域模型、Repository 接口定义
- **Data 层**: Repository 实现、数据访问（GORM）、缓存（Redis）

**依赖规则**: Service → Biz → Data（单向依赖，严禁反向引用）

### TypeScript/Vue
- Use `<script setup lang="ts">` for components
- Never use `as any` or `@ts-ignore`
- Tests: Vitest (unit), Playwright (E2E)

## Testing Patterns

### Table-Driven Tests
```go
tests := []struct {
    name     string
    input    string
    expected bool
}{
    {"valid", "https://example.com", true},
    {"invalid", "https://bad.com", false},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        assert.Equal(t, tt.expected, isValid(tt.input))
    })
}
```

### Skip External Dependencies
```go
client, err := redis.NewClient(cfg)
if err != nil {
    t.Skipf("redis not available: %v", err)
}
```

## Development Workflow

1. Define API in `api/protos/` → 2. `make gen` → 3. Implement biz → data → service → 4. `make wire` → 5. `make test` → 6. `make run`

## ⚠️ 禁止事项

- ❌ **不要修改生成的代码**: `api/gen/go/`、`wire_gen.go`、`*.pb.go` 等
- ❌ **不要跳过代码生成**: 修改 proto 后必须 `make gen`，修改 DI 后必须 `make wire`
- ❌ **不要在 Go 中使用 `panic()`**: 使用 Kratos 错误类型（如 `userv1.ErrorUserNotFound`）
- ❌ **不要在 TypeScript 中使用 `as any` 或 `@ts-ignore`**
- ❌ **不要提交生成的文件**: 已在 `.gitignore` 中配置
- ❌ **不要跨层调用**: Service 层不能直接调用 Data 层，必须通过 Biz 层

## 📚 详细文档引用

遇到以下情况时，应主动查阅对应的 `AGENTS.md` 获取详细指导：

| 场景 | 查阅文档 |
|------|---------|
| 项目概览、开发工作流 | `AGENTS.md` (根目录) |
| 修改 API 定义 | `api/AGENTS.md`、`api/protos/AGENTS.md` |
| 实现业务逻辑（DDD） | `app/krathub/service/internal/AGENTS.md` |
| Wire 依赖注入 | `app/AGENTS.md` |
| 前端开发 | `app/krathub/service/web/AGENTS.md` |
| 修改共享库 | `pkg/AGENTS.md` 和对应子目录文档 |

**提示**: AGENTS.md 包含详细的代码示例、架构图、最佳实践和常见问题解答。

## Common Pitfalls

- Run `make gen` after modifying proto files
- Run `make wire` after changing DI
- Frontend E2E requires `npx playwright install` first
- Use `t.Skipf()` for tests needing external services
- Never commit generated files (in `.gitignore`)
