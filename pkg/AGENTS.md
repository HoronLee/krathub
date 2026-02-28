# AGENTS.md - pkg/ 共享库层

<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-02-09 | Updated: 2026-02-09 -->

## 目录概述

`pkg/` 是 servora 项目的共享库层，提供跨服务的通用基础设施组件和工具函数。所有 `app/` 下的微服务都可以通过 `github.com/horonlee/servora/pkg/{module}` 导入这些模块。

**设计原则**：
- **无状态**：pkg 模块应该是纯工具库，不依赖特定服务的业务逻辑
- **可复用**：代码应该通用化，避免硬编码服务特定的配置
- **独立测试**：每个模块必须包含单元测试，可独立运行
- **最小依赖**：减少对外部库的依赖，避免引入过重的依赖树

## 子目录结构

### `jwt/` - JWT 认证工具
通用的 JWT 令牌生成和解析库，使用 Go 泛型支持自定义 Claims。

**核心文件**：
- `jwt.go` - 泛型 JWT 服务实现

**主要功能**：
- `JWT[T]` - 泛型 JWT 结构体，T 为自定义 Claims 类型
- `GenerateToken()` - 生成 JWT 令牌
- `ParseToken()` - 解析和验证 JWT 令牌
- `NewContext()` / `FromContext()` - 在 context 中存储/提取 Claims

**使用示例**：
```go
type MyClaims struct {
    UserID uint
    jwt.RegisteredClaims
}

jwtService := jwt.NewJWT[MyClaims](&jwt.Config{SecretKey: "secret"})
token, _ := jwtService.GenerateToken(&MyClaims{UserID: 123})
claims, _ := jwtService.ParseToken(token)
```

**依赖**：
- `github.com/golang-jwt/jwt/v5`

### `redis/` - Redis 客户端封装
提供 Redis 连接管理和常用操作的封装，支持从 Protobuf 配置创建客户端。

**核心文件**：
- `redis.go` - Redis 客户端实现
- `redis_test.go` - 单元测试

**主要功能**：
- `NewClient()` - 创建 Redis 客户端（带连接测试和清理函数）
- `NewConfigFromProto()` - 从 Protobuf 配置创建 Config
- 基础操作：`Set()`, `Get()`, `Del()`, `Has()`, `Keys()`, `Expire()`
- 集合操作：`SAdd()`, `SMembers()`

**配置**：
- 默认超时：DialTimeout 5s, ReadTimeout 3s, WriteTimeout 3s
- 支持用户名/密码认证、DB 选择

**依赖**：
- `github.com/redis/go-redis/v9`
- `github.com/horonlee/servora/api/gen/go/conf/v1` (配置 Proto)

### `logger/` - 日志封装
基于 Zap 的日志库封装，适配 Kratos 日志接口，支持多环境配置和日志切割。

**核心文件**：
- `log.go` - Zap 日志适配器
- `gorm_log.go` - GORM 日志适配器

**主要功能**：
- `NewLogger()` - 创建 Zap 日志实例，实现 `log.Logger` 接口
- `WithModule()` - 为日志添加 module 标签（命名规范：`[组件]/[层]/[服务名]`）
- `GetGormLogger()` - 获取 GORM 日志适配器

**环境模式**：
- `dev` - 终端彩色输出，不输出到文件
- `prod` - 终端非彩色 + 文件 JSON 输出（带日志切割）
- `test` - 不输出日志（NopCore）

**日志切割**：
- 使用 `lumberjack` 实现自动切割
- 默认配置：MaxSize 10MB, MaxBackups 5, MaxAge 30 天

**依赖**：
- `go.uber.org/zap`
- `gopkg.in/natefinch/lumberjack.v2`

### `middleware/` - 通用中间件
HTTP 中间件集合，用于跨域、白名单等通用功能。

**子目录**：
- `cors/` - CORS 跨域中间件
  - `cors.go` - CORS 中间件实现
  - `cors_test.go` - 单元测试
- `whitelist.go` - 白名单中间件（根目录）

**CORS 中间件功能**：
- `Middleware()` - 创建 CORS 中间件（接受 `conf.CORS` 配置）
- 支持动态启用/禁用（通过 `Enable` 字段）
- 支持通配符域名（如 `*.example.com`）
- 预检请求处理（OPTIONS）
- 默认配置：允许所有源、常用 HTTP 方法、标准头部

**使用示例**：
```go
corsMiddleware := cors.Middleware(corsConfig) // corsConfig 为 *conf.CORS
httpSrv := http.NewServer(
    http.Middleware(corsMiddleware),
)
```

**依赖**：
- `github.com/horonlee/servora/api/gen/go/conf/v1`

### `governance/` - 服务治理
服务注册发现和配置中心的抽象和实现，支持多种服务治理组件。

**子目录结构**：
- `registry/` - 服务注册与发现
  - `etcd.go` / `etcd_test.go` - etcd 注册中心实现
  - `consul.go` - Consul 注册中心实现
  - `nacos.go` - Nacos 注册中心实现
  - `kubernetes.go` / `kubernetes_test.go` - Kubernetes 服务发现
  - `etcd_watcher.go` - etcd 服务监听器
  - `etcd_error_test.go` - 错误处理测试
- `configCenter/` - 配置中心
  - `etcd.go` / `etcd_test.go` - etcd 配置中心
  - `consul.go` - Consul 配置中心
  - `nacos.go` - Nacos 配置中心

**支持的服务治理组件**：
- **etcd** - 分布式 KV 存储（推荐用于云原生环境）
- **Consul** - HashiCorp 服务网格
- **Nacos** - 阿里云微服务治理平台
- **Kubernetes** - K8s 原生服务发现

**使用场景**：
- 微服务注册与发现（服务间调用）
- 动态配置管理（配置热更新）
- 健康检查与负载均衡

### `k8s/` - Kubernetes 工具
Kubernetes 客户端封装，用于与 K8s API 交互。

**核心文件**：
- `client.go` - K8s 客户端实现
- `client_test.go` - 单元测试

**主要功能**：
- K8s 资源操作（ConfigMap、Secret、Service 等）
- 用于服务发现和配置管理

**依赖**：
- `k8s.io/client-go`
- `k8s.io/api`

### `helpers/` - 辅助函数
通用工具函数集合。

**子目录**：
- `hash/` - 密码哈希工具
  - `hash.go` - bcrypt 密码哈希实现

**主要功能**：
- 密码加密和验证
- 其他通用辅助函数

### `mapper/` - 数据映射工具
提供数据结构转换的工具函数，特别是 Protobuf 和内部模型之间的映射。

**核心文件**：
- `mapper.go` - 通用映射接口
- `proto_mapper.go` - Protobuf 映射工具
- `converter.go` - 数据转换工具

**主要功能**：
- DTO ↔ Entity 转换
- Protobuf ↔ Domain Model 转换
- 批量数据映射

### `transport/` - 传输层客户端
微服务间通信的客户端封装，主要用于 gRPC 连接管理。

**子目录**：
- `client/` - 客户端实现
  - `client.go` - 客户端接口和 ProviderSet
  - `connection.go` - 连接管理
  - `factory.go` - 客户端工厂
  - `grpc_conn.go` - gRPC 连接实现

**主要功能**：
- gRPC 连接池管理
- 服务发现集成
- 连接复用和健康检查

**连接类型**：
- `GRPC` - gRPC 连接（当前实现）
- `WebSocket` - WebSocket 连接（预留）
- `HTTP` - HTTP 连接（预留）

**依赖**：
- `github.com/google/wire` (依赖注入)

## AI Agent 工作指南

### 添加新的 pkg 模块

**标准流程**：
1. **创建目录** - 在 `pkg/` 下创建新模块目录
   ```bash
   mkdir -p pkg/newmodule
   ```

2. **编写代码** - 创建主文件和测试文件
   ```bash
   touch pkg/newmodule/newmodule.go
   touch pkg/newmodule/newmodule_test.go
   ```

3. **编写测试** - 必须包含单元测试，覆盖核心功能
   ```go
   // pkg/newmodule/newmodule_test.go
   package newmodule

   import "testing"

   func TestNewModule(t *testing.T) {
       // 测试代码
   }
   ```

4. **运行测试** - 确保测试通过
   ```bash
   go test -v ./pkg/newmodule/...
   ```

5. **更新文档** - 在本文件中添加模块说明

**设计要求**：
- 模块必须是无状态的纯工具库
- 避免依赖特定服务的业务逻辑
- 优先使用依赖注入（通过 Wire）而非全局变量
- 提供清晰的构造函数（如 `NewXxx()`）
- 支持从 Protobuf 配置创建（如 `NewConfigFromProto()`）

### 测试要求

**必须测试的场景**：
1. **正常功能** - 核心功能的正常路径
2. **边界条件** - 空值、nil、极端值
3. **错误处理** - 错误输入、异常情况

**测试模式**：
- 使用表驱动测试（Table-Driven Tests）
- 对外部依赖（Redis、数据库等）使用 `t.Skipf()` 优雅跳过
- Mock 外部依赖（使用接口抽象）

**示例 - 表驱动测试**：
```go
func TestJWTGenerate(t *testing.T) {
    tests := []struct {
        name      string
        secretKey string
        claims    *MyClaims
        wantErr   bool
    }{
        {"valid", "secret123", &MyClaims{UserID: 1}, false},
        {"empty_secret", "", &MyClaims{UserID: 1}, false}, // JWT 允许空密钥
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            jwtService := NewJWT[MyClaims](&Config{SecretKey: tt.secretKey})
            token, err := jwtService.GenerateToken(tt.claims)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && token == "" {
                t.Error("GenerateToken() returned empty token")
            }
        })
    }
}
```

**示例 - 跳过外部依赖**：
```go
func TestRedisClient(t *testing.T) {
    cfg := &Config{Addr: "localhost:6379"}
    client, cleanup, err := NewClient(cfg, log.DefaultLogger)
    if err != nil {
        t.Skipf("redis not available: %v", err)
    }
    defer cleanup()

    // 测试逻辑
}
```

### 代码风格规范

**导入顺序**（遵循项目规范）：
```go
import (
    "context"           // 1. 标准库
    "time"

    "github.com/go-kratos/kratos/v2/log"  // 2. 第三方库
    "github.com/redis/go-redis/v9"

    conf "github.com/horonlee/servora/api/gen/go/conf/v1"  // 3. 项目内
)
```

**命名约定**：
- 包名：小写单数（`jwt`, `redis`, `logger`）
- 导出函数：PascalCase（`NewClient`, `GenerateToken`）
- 私有函数：camelCase（`parseConfig`, `isOriginAllowed`）
- 构造函数：统一使用 `New` 前缀

**错误处理**：
- 优先返回错误而非 panic
- 使用 `fmt.Errorf` 包装错误信息
- 对于库代码，不要使用 `log.Fatal`（会终止程序）

**配置管理**：
- 使用结构体传递配置（如 `Config`）
- 提供默认值（如 `DefaultDialTimeout`）
- 支持从 Protobuf 配置转换（如 `NewConfigFromProto`）

### 依赖注入集成

**Wire ProviderSet**：
如果模块需要在微服务中通过 Wire 注入，提供 `ProviderSet`：

```go
// pkg/newmodule/wire.go
package newmodule

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
    NewClient,
    NewConfigFromProto,
)
```

**在服务中使用**：
```go
// app/service/cmd/server/wire.go
var dataSet = wire.NewSet(
    newmodule.ProviderSet,
    // 其他依赖
)
```

### 常见任务

**添加新的中间件**：
1. 在 `pkg/middleware/` 下创建子目录（如 `ratelimit/`）
2. 实现中间件函数签名：`func(http.Handler) http.Handler`
3. 支持从 Protobuf 配置创建
4. 编写单元测试（测试中间件逻辑）
5. 在服务中通过 `http.Middleware()` 注册

**添加新的服务治理实现**：
1. 在 `pkg/governance/registry/` 或 `configCenter/` 下创建文件
2. 实现 Kratos 的 `registry.Registrar` 或 `config.Source` 接口
3. 提供构造函数和配置结构
4. 编写单元测试和集成测试
5. 更新文档说明支持的组件

**优化日志模块**：
1. 编辑 `pkg/logger/log.go`
2. 保持与 Kratos `log.Logger` 接口兼容
3. 更新 `GetGormLogger()` 以支持新功能
4. 测试多环境模式（dev/prod/test）
5. 更新 module 命名规范（如需要）

## 依赖关系

### 核心依赖
- **Go 1.21+** - Go 运行时
- **Kratos v2** - 微服务框架（部分模块依赖）
  - `github.com/go-kratos/kratos/v2/log` (日志接口)
  - `github.com/go-kratos/kratos/v2/registry` (注册中心接口)
  - `github.com/go-kratos/kratos/v2/config` (配置接口)

### 第三方库依赖（按模块）
- **jwt/**
  - `github.com/golang-jwt/jwt/v5` - JWT 实现
- **redis/**
  - `github.com/redis/go-redis/v9` - Redis 客户端
- **logger/**
  - `go.uber.org/zap` - 高性能日志库
  - `gopkg.in/natefinch/lumberjack.v2` - 日志切割
- **governance/**
  - `go.etcd.io/etcd/client/v3` - etcd 客户端
  - `github.com/hashicorp/consul/api` - Consul 客户端
  - `github.com/nacos-group/nacos-sdk-go/v2` - Nacos 客户端
- **k8s/**
  - `k8s.io/client-go` - Kubernetes 客户端
  - `k8s.io/api` - Kubernetes API 类型
- **helpers/hash/**
  - `golang.org/x/crypto/bcrypt` - 密码哈希
- **transport/**
  - `github.com/google/wire` - 依赖注入

### 项目内部依赖
所有 pkg 模块可能依赖：
- `github.com/horonlee/servora/api/gen/go/conf/v1` - 配置 Protobuf 定义（生成的代码）

**依赖方向**：
- `app/` → `pkg/` (微服务依赖共享库)
- `pkg/` → `api/gen/go/` (共享库依赖生成的配置)
- `pkg/` 模块间应避免循环依赖

## 注意事项

### 测试注意事项
1. **外部服务依赖** - 测试需要 Redis、etcd 等外部服务时，使用 `t.Skipf()` 跳过
   ```go
   if err := client.Ping(ctx); err != nil {
       t.Skipf("redis not available: %v", err)
   }
   ```

2. **并发测试** - 对于并发安全的模块（如连接池），添加并发测试
   ```go
   t.Run("concurrent", func(t *testing.T) {
       t.Parallel()
       // 并发测试逻辑
   })
   ```

3. **测试覆盖率** - 运行测试覆盖率检查
   ```bash
   go test -v -cover ./pkg/...
   ```

### 配置管理
- 所有模块应支持从 Protobuf 配置创建（通过 `NewConfigFromProto()`）
- 配置结构体应提供合理的默认值
- 避免硬编码配置，使用配置参数

### 版本兼容性
- 保持 API 向后兼容，避免破坏性修改
- 重大修改应使用新的包名或版本号
- 在注释中说明最小依赖版本

### 性能考虑
- 避免在初始化时进行重量级操作
- 提供资源清理函数（cleanup function）
- 连接池、缓存等资源应可复用

### 安全注意事项
- 不要在日志中输出敏感信息（密码、Token、密钥）
- JWT 密钥应从配置读取，不要硬编码
- Redis 密码等敏感配置应通过环境变量或密钥管理系统注入

## 常见陷阱

1. **循环依赖** - pkg 模块间不应相互依赖，保持单向依赖图
2. **全局状态** - 避免使用全局变量，优先使用依赖注入
3. **panic 使用** - 库代码不应 panic，应返回错误
4. **测试污染** - 测试间应相互独立，避免共享状态
5. **版本不匹配** - 确保 `api/gen/go/` 代码是最新生成的（运行 `make gen`）

## 快速参考

**运行所有 pkg 测试**：
```bash
go test -v ./pkg/...
```

**测试单个模块**：
```bash
go test -v ./pkg/jwt/...
go test -v ./pkg/redis/...
```

**查看测试覆盖率**：
```bash
go test -v -cover ./pkg/...
```

**查看详细覆盖率报告**：
```bash
go test -coverprofile=coverage.out ./pkg/...
go tool cover -html=coverage.out
```

**在服务中导入 pkg 模块**：
```go
import (
    "github.com/horonlee/servora/pkg/jwt"
    "github.com/horonlee/servora/pkg/redis"
    "github.com/horonlee/servora/pkg/logger"
)
```

**添加 Wire 依赖**：
```go
// app/{service}/cmd/server/wire.go
var dataSet = wire.NewSet(
    redis.ProviderSet,
    // 其他 pkg 模块的 ProviderSet
)
```
