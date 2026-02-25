# Krathub

[English](README.en-US.md) | 简体中文

Krathub 是一个基于 Go Kratos 框架的微服务项目示例。它展示了现代化微服务架构的最佳实践，采用 Buf 进行 Protobuf 管理，支持多服务架构，可作为学习参考或作为新项目的起点。

## ✨ 核心特性

- **现代化微服务架构**: 基于 Kratos v2 构建，支持多服务独立开发和部署
- **Buf 工具链**: 使用 Buf 进行 Protobuf 管理，提供更好的依赖管理和代码生成体验
- **双协议支持**: 同时提供 gRPC 和 HTTP 接口，满足不同场景需求
- **统一 OpenAPI 生成**: 自动为每个服务生成独立的 OpenAPI 文档
- **服务治理**: 集成 Consul、Nacos 和 etcd，提供开箱即用的服务注册与发现能力
- **配置中心**: 支持通过 Consul、Nacos 或 etcd 进行动态配置管理
- **数据库集成**: 采用 GORM 作为 ORM，支持 MySQL、PostgreSQL 和 SQLite
- **依赖注入**: 使用 Wire 进行依赖注入，清晰化对象之间的依赖关系
- **通用工具链**: 提供 `app.mk` 通用 Makefile，所有服务共享统一的构建流程
- **认证鉴权**: 内置 JWT 中间件，方便实现用户认证
- **容器化**: 提供 `Dockerfile` 和 `docker-compose.yml`，轻松实现容器化部署
- **可观测性**: 已集成 `Metrics` (Prometheus) 和 `Trace` (Jaeger) 的基础配置

## 📖 项目文档

- **API 文档**: [在线文档](https://jqovxjvrtw.apifox.cn) - 查看和测试所有 API 接口
- **项目配置**: `configs/config.yaml` - 完整的配置选项说明

## 🚀 快速开始

### 前置要求

- Go 1.21 或更高版本
- Buf CLI (用于 Protobuf 管理)
- Wire (用于依赖注入)
- Make 工具

### 克隆项目

```shell
# 克隆项目
git clone https://github.com/horonlee/krathub.git
cd krathub
```

### 安装开发工具
```
# 通过Make初始化开发环境
make init
```

### 配置项目

```shell
# 复制配置示例文件
cp api/protos/conf/v1/config-example.yaml app/krathub/service/configs/config.yaml

# 根据需要修改配置
vim app/krathub/service/configs/config.yaml
```

### 启动依赖服务

Krathub 需要以下依赖服务：

```shell
# 使用 Docker 启动 Redis
docker run -d -p 6379:6379 redis:latest

# 使用 Docker 启动 MySQL（可选，也可以使用 SQLite）
docker run -d -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=krathub \
  mysql:8.0
```

然后修改配置文件中的数据库和 Redis 连接信息。
```shell
vim app/krathub/service/configs/config.yaml
```

### 生成代码并启动

```shell
# 生成所有代码（protobuf、wire、openapi）
make gen

# 构建并运行 krathub 服务
cd app/krathub/service
make run
```

### 使用 Compose + Air 开发（推荐）

```shell
# 先在宿主机生成代码（避免 Air 触发生成导致重启循环）
make gen

# 构建 Air 开发镜像
make compose.dev.build

# 启动双服务热重载开发环境（krathub + sayhello）
make compose.dev.up

# 查看实时日志
make compose.dev.logs

# 停止开发环境
make compose.dev.down
```

服务启动后，HTTP 服务将监听在 `0.0.0.0:8000`，gRPC 服务将监听在 `0.0.0.0:8001` (以默认配置为例)。

## 📁 项目结构

```
.
├── api/                                # Protobuf API 定义与代码生成相关配置
│   ├── buf.gen.yaml                    # Buf 代码生成配置（Go）
│   ├── buf.work.yaml                   # Buf workspace 配置
│   ├── buf.{service}.openapi.gen.yaml  # 各服务 OpenAPI 生成配置
│   ├── gen/                            # 生成的代码
│   │   └── go/                         # 生成的 Go protobuf 代码
│   └── protos/                         # Proto 源文件
│       ├── buf.yaml                    # Proto 依赖配置
│       ├── conf/v1/                    # 配置定义（proto）与配置示例
│       ├── krathub/service/v1/         # Krathub HTTP 接口（i_*.proto）
│       ├── auth/service/v1/            # Auth gRPC 服务
│       ├── user/service/v1/            # User gRPC 服务
│       ├── test/service/v1/            # Test gRPC 服务
│       └── sayhello/service/v1/        # SayHello 独立微服务
│
├── app/                                # 微服务应用目录
│   ├── krathub/service/                # Krathub 主服务
│   │   ├── cmd/server/                 # 服务启动入口
│   │   ├── internal/                   # 内部实现（不对外暴露）
│   │   │   ├── biz/                    # 业务逻辑层
│   │   │   ├── data/                   # 数据访问层
│   │   │   ├── server/                 # gRPC/HTTP 服务器配置
│   │   │   └── service/                # Service 层实现
│   │   ├── configs/                    # 服务配置文件（运行时 config.yaml）
│   │   ├── bin/                        # 编译输出目录
│   │   ├── openapi.yaml                # 生成的 OpenAPI 文档
│   │   └── Makefile                    # 服务级 Makefile（include app.mk）
│   │
│   └── sayhello/service/               # SayHello 独立微服务（示例）
│       ├── openapi.yaml                # 生成的 OpenAPI 文档
│       └── Makefile                    # 服务级 Makefile
│
├── manifest/                           # 部署相关文件
│   ├── SQL/                            # 数据库脚本
│   ├── docker/                         # Docker 配置
│   └── kubernetes/                     # K8s 配置
│
├── pkg/                                # 项目内部共享的通用库
├── examples/                           # 示例项目
│
├── .env.example                        # 环境变量示例（需复制为 .env）
├── .env                                # 本地环境变量（建议加入 .gitignore）
├── app.mk                              # 通用服务 Makefile（所有服务共享）
└── Makefile                            # 根目录 Makefile（管理所有服务
```


### Proto 文件组织规范

项目采用以下 Proto 文件组织规范：

1. **HTTP 接口文件** (`i_*.proto`)
   - 位置：`api/protos/krathub/service/v1/i_*.proto`
   - 包名：统一使用 `krathub.service.v1`
   - 用途：包含 HTTP 注解的接口定义，用于生成 OpenAPI 文档
   - 示例：`i_auth.proto`, `i_user.proto`, `i_test.proto`

2. **gRPC 服务文件**
   - 位置：`api/protos/{service}/service/v1/{service}.proto`
   - 包名：独立包名 `{service}.service.v1`
   - 用途：纯 gRPC 接口定义，不包含 HTTP 注解
   - 示例：`auth/service/v1/auth.proto`, `user/service/v1/user.proto`

3. **独立微服务**
   - 位置：`api/protos/{service}/service/v1/{service}.proto`
   - 包名：独立包名 `{service}.service.v1`
   - 用途：完全独立的微服务，可包含 HTTP 注解
   - 示例：`sayhello/service/v1/sayhello.proto`

## 📝 开发流程

推荐的开发顺序如下，以确保依赖关系正确：

### 1. 定义 API (`api/protos/`)

在 `.proto` 文件中定义服务接口：

```protobuf
// HTTP 接口：api/protos/krathub/service/v1/i_example.proto
syntax = "proto3";
package krathub.service.v1;

import "google/api/annotations.proto";

service Example {
  rpc GetExample(GetExampleRequest) returns (GetExampleResponse) {
    option (google.api.http) = {
      get: "/api/v1/example/{id}"
    };
  }
}
```

### 2. 生成代码

```shell
# 在根目录生成所有 protobuf 代码
make api

# 为所有服务生成 OpenAPI 文档
make openapi

# 或者一次性生成所有代码
make gen
```

### 3. 实现业务逻辑

按照 Kratos 的分层架构实现：

1. **业务逻辑层** (`internal/biz/`): 定义业务接口和实现
2. **数据访问层** (`internal/data/`): 实现数据持久化
3. **服务层** (`internal/service/`): 实现 API 接口

### 4. 依赖注入

在 `cmd/server/wire.go` 中注册新组件：

```go
//go:build wireinject
// +build wireinject

func wireApp(...) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        newApp,
    ))
}
```

生成依赖注入代码：

```shell
# 在服务目录下
cd app/krathub/service
make wire
```

### 5. 运行和测试

```shell
# 在服务目录下运行
cd app/krathub/service
make run

# 或者在根目录构建所有服务
make build
```

## 🛠️ 常用命令

### 根目录命令（管理所有服务）

#### 初始化和设置
- `make init` - 安装所有开发工具（buf, wire, protoc 插件等）
- `make plugin` - 安装 protoc 插件
- `make cli` - 安装 CLI 工具（kratos, buf, wire 等）

#### 代码生成
- `make api` - 生成所有 protobuf Go 代码
- `make openapi` - 为所有服务生成 OpenAPI 文档
- `make wire` - 为所有服务生成 wire 依赖注入代码
- `make gen` - 生成所有代码（api + openapi + wire）

#### 构建和运行
- `make build` - 构建所有服务（包含代码生成）
- `make build_only` - 仅构建所有服务（不生成代码）
- `make all` - 生成并构建所有服务

#### 代码质量
- `make test` - 运行所有测试
- `make lint` - 运行代码检查
- `make vet` - 运行静态分析

#### 其他
- `make clean` - 清理所有构建产物
- `make env` - 显示环境变量
- `make help` - 显示帮助信息
- `make compose.dev.build` - 构建 Air 开发镜像（基于 `docker-compose.yaml` + `docker-compose.dev.yaml`）
- `make compose.dev.up` - 启动 Air 热重载开发容器
- `make compose.dev.ps` - 查看 Air 开发容器状态
- `make compose.dev.logs` - 查看 Air 开发容器日志
- `make compose.dev.down` - 停止 Air 开发容器
- `make compose.build` - 构建生产镜像（krathub + sayhello）
- `make compose.up` - 启动生产 compose 全栈（中间件 + 微服务）
- `make compose.rebuild` - 重建生产镜像并启动生产 compose 全栈
- `make compose.ps` - 查看生产 compose 服务状态
- `make compose.logs` - 查看生产 compose 服务日志
- `make compose.down` - 停止生产 compose 全栈

### 服务级命令（在服务目录下执行）

进入服务目录：`cd app/krathub/service`

#### 开发命令
- `make run` - 运行服务（包含代码生成）
- `make build` - 构建服务（包含代码生成）
- `make build_only` - 仅构建服务
- `make app` - 生成并构建服务

#### 代码生成
- `make api` - 生成 protobuf 代码
- `make openapi` - 生成 OpenAPI 文档
- `make wire` - 生成 wire 代码
- `make ent` - 生成 ent 代码（如果使用 ent）
- `make gen.dao` - 生成 GORM GEN 的 PO 和 DAO 代码（如果使用 gorm-gen）
- `make gen` - 生成所有代码

#### 其他
- `make clean` - 清理构建产物
- `make docker-build` - 构建 Docker 镜像
- `make env` - 显示环境变量

### Buf 相关命令

- `make lint-proto` - 检查 proto 文件规范
- `make buf-update` - 更新 buf 依赖

### 添加新服务

1. 创建服务目录结构：
```shell
mkdir -p app/newservice/service
```

2. 创建服务 Makefile：
```shell
echo "include ../../../app.mk" > app/newservice/service/Makefile
```

3. 创建 OpenAPI 生成配置：
```shell
# 复制并修改现有配置
cp api/buf.krathub.openapi.gen.yaml api/buf.newservice.openapi.gen.yaml
# 修改配置中的服务名称和路径
```

4. 定义 proto 文件并生成代码：
```shell
# 创建 proto 文件
mkdir -p api/protos/newservice/service/v1
# 编写 proto 文件...

# 生成代码
make api
make openapi
```

5. 实现服务代码（参考 krathub 服务结构）

根目录的 Makefile 会自动发现并管理新服务！

## 🗄️ 数据库支持

项目支持多种数据库后端，通过配置文件即可切换：

- **MySQL** (默认) - 企业级关系型数据库
- **PostgreSQL** - 高性能开源关系型数据库
- **SQLite** - 轻量级嵌入式数据库

数据库配置示例请参考 `app/krathub/service/configs/config.yaml`。

## 🐳 Docker 支持

每个微服务都有独立的 Dockerfile，位于 `app/{service}/service/Dockerfile`。

项目根目录采用两份 Compose 文件：

- `docker-compose.yaml`：生产部署编排（中间件 + 微服务，按镜像运行）
- `docker-compose.dev.yaml`：开发覆盖层（仅替换 Go 服务为 Air 热重载）
- `Dockerfile.air`：开发热重载镜像

```shell
# 构建所有服务的 Docker 镜像
make docker-build

# 构建单个服务的 Docker 镜像
cd app/krathub/service
make docker-build
```

## 🔧 配置说明

每个服务的配置文件位于 `app/{service}/service/configs/config.yaml`，支持通过环境变量覆盖默认值。

**核心配置项：**

- **服务配置** - HTTP/gRPC 服务地址、TLS 证书等
- **数据库配置** - 支持 MySQL/PostgreSQL/SQLite
- **Redis 配置** - 缓存和会话存储
- **服务治理** - Consul/Nacos/etcd 注册发现
- **配置中心** - 支持从 Consul/Nacos/etcd 动态加载配置
- **JWT 认证** - 用户认证和授权
- **日志配置** - 日志级别、文件轮转等
- **链路追踪** - Jaeger/Zipkin 配置
- **监控指标** - Prometheus 配置

完整的配置示例请参考 `app/krathub/service/configs/config-example.yaml` 文件。

## 📚 技术栈

- **框架**: [Kratos v2](https://go-kratos.dev/) - Go 微服务框架
- **Protobuf 管理**: [Buf](https://buf.build/) - 现代化的 Protobuf 工具链
- **依赖注入**: [Wire](https://github.com/google/wire) - 编译时依赖注入
- **ORM**: [GORM](https://gorm.io/) - Go ORM 库
- **缓存**: [Redis](https://redis.io/) - 内存数据库
- **服务发现**: Consul / Nacos / etcd
- **链路追踪**: Jaeger / Zipkin
- **监控**: Prometheus + Grafana

## 🎯 最佳实践

### Proto 文件管理

1. **HTTP 接口统一管理**: 所有带 HTTP 注解的接口放在 `krathub/service/v1/i_*.proto`
2. **gRPC 服务独立包**: 每个 gRPC 服务使用独立的包名，避免类型冲突
3. **使用 Buf 管理依赖**: 通过 `buf.build` 远程依赖，无需本地 `third_party`

### 服务开发

1. **分层架构**: 严格遵循 biz -> data -> service 的分层结构
2. **依赖注入**: 使用 Wire 管理依赖，保持代码清晰
3. **配置管理**: 使用配置中心进行动态配置管理
4. **错误处理**: 使用 Kratos 的错误处理机制，统一错误码

### 代码组织

1. **服务独立**: 每个服务在 `app/{service}/service` 下独立开发
2. **共享代码**: 通用代码放在 `pkg/` 目录下
3. **配置独立**: 每个服务有自己的配置文件
4. **文档自动生成**: OpenAPI 文档自动生成到服务目录

## 📞 客户端层 (Client)

`internal/client` 是一个自定义的层，用于管理对外部 gRPC 等服务的调用。它提供了一个客户端工厂，可以方便地基于服务发现创建和复用连接，专为服务间的通信而设计。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

在提交 PR 之前，请确保：

1. 代码通过 `make lint` 检查
2. 所有测试通过 `make test`
3. 更新相关文档

## 📄 许可协议

本项目遵循 [LICENSE](LICENSE) 文件中的许可协议。
