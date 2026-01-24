# CLAUDE.md - Krathub Development Guide

Instructions for AI assistants working in this project.

无论何时，用中文回答

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
make genDao        # Generate GORM GEN PO/DAO
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
├── api/protos/              # Proto definitions (i_*.proto = HTTP, others = gRPC)
├── api/gen/go/              # Generated Go code
├── app/{service}/service/   # Microservices (cmd/, internal/biz|data|service|server/)
├── pkg/                     # Shared packages (jwt, redis, logger, hash)
└── app.mk                   # Common Makefile for services
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

### Layered Architecture
`service/` → `biz/` → `data/` (handlers → business logic → repository)

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

## Common Pitfalls

- Run `make gen` after modifying proto files
- Run `make wire` after changing DI
- Frontend E2E requires `npx playwright install` first
- Use `t.Skipf()` for tests needing external services
- Never commit generated files (in `.gitignore`)
