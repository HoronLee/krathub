# Krathub

[English](README.en-US.md) | 简体中文

Krathub 是一个基于 Go Kratos 框架的微服务项目模板。它集成了一系列最佳实践和常用组件，旨在帮助开发者快速构建一个功能完善、结构清晰、易于扩展的现代化 Go 应用。

## ✨ 核心特性

- **微服务架构**: 基于 Kratos v2 构建，天然支持微服务。
- **双协议支持**: 同时提供 gRPC 和 HTTP 接口，满足不同场景需求。
- **服务治理**: 集成 Consul、Nacos 和 etcd，提供开箱即用的服务注册与发现能力。
- **配置中心**: 支持通过 Consul、Nacos 或 etcd 进行动态配置管理。
- **数据库集成**: 采用 GORM 作为 ORM，并提供 `make gen.db` 快速生成模型代码。
- **依赖注入**: 使用 Wire 进行依赖注入，清晰化对象之间的依赖关系。
- **代码生成**: 大量使用 `make` 命令简化 `proto`、`wire` 等代码的生成。
- **认证鉴权**: 内置 JWT 中间件，方便实现用户认证。
- **容器化**: 提供 `Dockerfile` 和 `docker-compose.yml`，轻松实现容器化部署。
- **可观测性**: 已集成 `Metrics` (Prometheus) 和 `Trace` (Jaeger) 的基础配置。

## 📖 项目文档

- **API 文档**: [在线文档](https://jqovxjvrtw.apifox.cn) - 查看和测试所有 API 接口
- **项目配置**: `configs/config.yaml` - 完整的配置选项说明

## 🚀 快速开始

### 使用 kratos layout 功能

1. 先安装 kratos cli 工具
```shell
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

2. 然后使用layout功能从模板生成项目
```shell
kratos new <ProjectName> -r https://github.com/horonlee/krathub.git
```

3. 进入项目目录并安装依赖
```shell
cd <ProjectName>
make setup
```

4. 从configs中的config-example.yaml开始配置项目

5. 启动项目
```shell
make run
```
 服务启动后，HTTP 服务将监听在 `0.0.0.0:8000`，gRPC 服务将监听在 `0.0.0.0:8001` (以默认配置为例)。

## 📁 项目结构

```
.
├── api/         # Protobuf API 定义 (gRPC & HTTP)
├── cmd/         # 主程序入口和启动逻辑
├── configs/     # 配置文件
├── internal/    # 核心业务逻辑 (不对外暴露)
│   ├── biz/     # 业务逻辑层 (struct 和 interface)
│   ├── client/  # 客户端层，用于服务间调用
│   ├── data/    # 数据访问层 (数据库、缓存)
│   ├── conf/    # Protobuf 定义的配置结构
│   ├── server/  # gRPC 和 HTTP 服务的创建和配置
│   ├── service/ # Service 层，实现 API 定义的接口
│   └── ...
├── manifest/    # 部署相关文件 (SQL, Docker, K8s, 证书)
├── pkg/         # 可在项目内部共享的通用库
└── third_party/ # 第三方 proto 文件和依赖
```

## 📝 开发流程

推荐的开发顺序如下，以确保依赖关系正确：

1. **API 定义 (`api/`)**: 在 `.proto` 文件中定义 gRPC 服务和消息体。
2. **生成代码 (`make proto`)**: 运行命令生成 gRPC、HTTP、Errors 的桩代码。
3. **业务逻辑 (`internal/biz/`)**: 定义业务逻辑的接口和实现，这是不依赖任何框架的核心。
4. **数据访问 (`internal/data/`)**: 实现 `biz` 层定义的接口，负责与数据库、缓存等交互。
5. **服务实现 (`internal/service/`)**: 实现 `api` 层定义的 gRPC 服务接口，它会调用 `biz` 层的逻辑。
6. **依赖注入 (`cmd/server/wire.go`)**: 将新的 `service`, `biz`, `data` 组件注入到 `wire.go` 中。
7. **运行 `make wire`**: 生成最终的依赖注入代码。
8. **启动与测试**: 运行 `make run` 并进行测试。

## 🛠️ 常用命令

### 🚀 开发命令

- `make setup` - 一键设置开发环境（推荐）
- `make all` - 生成所有代码并构建项目
- `make run` - 启动应用程序
- `make check` - 运行代码质量检查
- `make help` - 查看所有可用命令

### 🏗️ 核心功能

- `make gen.db` - 生成数据库模型代码
- `make proto` - 生成 Protobuf 代码
- `make wire` - 生成依赖注入代码
- `make test` - 运行所有测试
- `make build` - 构建应用程序

### 📦 其他工具

- `make docker-build` - 构建 Docker 镜像
- `make clean` - 清理生成文件
- `make mod-tidy` - 整理 Go 依赖

### 🗄️ 数据库支持

项目支持多种数据库后端，通过配置文件即可切换：

- **MySQL** (默认) - 企业级关系型数据库
- **PostgreSQL** - 高性能开源关系型数据库
- **SQLite** - 轻量级嵌入式数据库

数据库配置示例请参考 `configs/config.yaml`。

### 🐳 Docker 支持

- `make docker-build`: 构建 Docker 镜像。
- `make docker-run`: 运行 Docker 容器。
- `make docker-push`: 推送 Docker 镜像。
- `make docker-clean`: 清理 Docker 资源。

### 🧹 清理命令

- `make clean`: 清理生成的文件。
- `make clean-all`: 清理所有文件和缓存。
- `make clean-deps`: 清理并重新安装依赖。

### ℹ️ 其他命令

- `make version`: 显示项目版本信息。
- `make info`: 显示项目信息（同 version）。

### 🔄 兼容性说明

所有原有命令都保持向后兼容，包括：
- `make init`: 初始化开发环境（现在是 `install-dev` 的别名）。

## 📞 客户端层 (Client)

`internal/client` 是一个自定义的层，用于管理对外部 gRPC 等服务的调用。它提供了一个客户端工厂，可以方便地基于服务发现创建和复用连接，专为服务间的通信而设计。

## ⚙️ 配置说明

项目的配置文件位于 `configs/config.yaml`，支持通过环境变量覆盖默认值。数据库支持 MySQL、PostgreSQL 和 SQLite，可通过修改配置文件切换。

**核心配置项：**

- **服务配置** - HTTP/gRPC 服务地址、TLS 证书等
- **数据库配置** - 支持 MySQL/PostgreSQL/SQLite，可通过环境变量切换
- **Redis配置** - 缓存和会话存储
- **服务治理** - Consul/Nacos/etcd 注册发现
- **JWT认证** - 用户认证和授权
- **日志配置** - 日志级别、文件轮转等

完整的配置示例请参考 `configs/config-example.yaml` 文件。

## 📄 许可协议

本项目遵循 [LICENSE](LICENSE) 文件中的许可协议。
