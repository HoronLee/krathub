# micro-forge

简体中文

micro-forge 是一个基于 **Go Kratos v2** 的微服务示例项目，采用 **DDD 分层** 与 **契约优先（Proto First）** 的开发方式，覆盖从 API 定义、代码生成、服务开发到可观测性与容器化部署的完整链路。

## ✨ 核心能力

- **微服务模板化**：统一的服务目录约定与 `app.mk` 共享构建流程
- **Proto First**：使用 Buf 管理 Protobuf 依赖与代码生成
- **双协议接口**：同时支持 gRPC 与 HTTP（含 OpenAPI 生成）
- **DDD 分层**：`service -> biz -> data`，职责边界清晰
- **依赖注入**：使用 Wire 进行编译期依赖注入
- **数据访问**：Ent + GORM GEN 双工具链并行
- **服务治理**：支持 Consul / Nacos / etcd 注册发现与配置中心
- **可观测性**：OTel Collector + Jaeger + Loki + Prometheus + Grafana
- **开发体验**：支持 Docker Compose + Air 热重载开发

## 🧱 技术栈

- 框架：Kratos v2
- API：Protobuf + Buf
- DI：Google Wire
- ORM：Ent（主）+ GORM GEN（并行）
- 存储：MySQL / PostgreSQL / SQLite + Redis
- 前端：Vue 3 + Vite + Bun（位于根目录 `web/`）
- 观测：OTel / Jaeger / Loki / Prometheus / Grafana

## 🗂️ 项目结构

```text
.
├── api/                         # Proto 定义、Buf 配置、生成代码
│   ├── protos/
│   ├── gen/go/
│   ├── buf.gen.yaml
│   └── buf.*.openapi.gen.yaml
├── app/
│   ├── micro-forge/service/         # 主服务（DDD 分层）
│   └── sayhello/service/        # 独立示例服务
├── pkg/                         # 项目共享库
├── web/                         # Vue 3 前端项目（根目录）
├── manifests/                   # 可观测性与证书等基础设施清单
├── deployment/                  # 部署相关配置
├── docker-compose.yaml          # 生产编排
├── docker-compose.dev.yaml      # 开发覆盖层（Air）
├── app.mk                       # 服务级通用 Makefile
└── Makefile                     # 根目录统一入口
```

## 🚀 快速开始

### 1) 前置要求

- Go 1.21+
- Make
- Docker / Docker Compose

> 首次初始化工具链（buf、wire、lint 等）：

```bash
make init
```

### 2) 克隆与基础配置

```bash
git clone https://github.com/horonlee/micro-forge.git
cd micro-forge

# 复制示例配置（主服务）
cp api/protos/conf/v1/config-example.yaml app/micro-forge/service/configs/config.yaml
```

按需修改 `app/micro-forge/service/configs/config.yaml` 中的数据库、Redis、注册中心等配置。

### 3) 生成代码

```bash
make gen
```

该命令会统一执行代码生成流程：`wire + api + openapi + ent`。

### 4) 本地运行（两种方式）

**方式 A：直接运行主服务**

```bash
cd app/micro-forge/service
make run
```

默认示例端口：

- HTTP: `0.0.0.0:8000`
- gRPC: `0.0.0.0:8001`

**方式 B：Compose + Air（推荐开发）**

```bash
# 在仓库根目录
make gen
make compose.dev.build
make compose.dev.up
```

查看日志与停止：

```bash
make compose.dev.logs
make compose.dev.down
```

## 🧭 开发工作流

推荐顺序：

1. 修改/新增 `.proto`（`api/protos/`）
2. 运行 `make gen` 同步生成代码
3. 按 DDD 分层实现：`internal/service -> internal/biz -> internal/data`
4. 若修改了 Wire 依赖图，运行 `make wire`（或直接 `make gen`）
5. 运行 `make test`、`make lint` 验证质量

## 🛠️ 常用命令

### 根目录命令

```bash
# 初始化工具
make init

# 代码生成
make gen
make api
make openapi
make wire
make ent

# 构建与质量
make build
make build_only
make test
make lint
make vet

# Compose（生产）
make compose.build
make compose.up
make compose.rebuild
make compose.ps
make compose.logs
make compose.down

# Compose（开发 Air）
make compose.dev.build
make compose.dev.up
make compose.dev.ps
make compose.dev.logs
make compose.dev.down
```

### 服务级命令（示例：`app/micro-forge/service/`）

```bash
make run
make build
make build_only
make app
make gen
make wire
make gen.ent
make gen.dao
make openapi
make docker-build
```

### 前端命令（`web/`）

```bash
cd web
bun install
bun dev
bun test:unit
bun test:e2e
bun lint
```

## 📦 配置说明

- 主服务配置：`app/micro-forge/service/configs/config.yaml`
- 示例配置：`api/protos/conf/v1/config-example.yaml`
- 支持环境变量覆盖默认值（详见示例配置中的 `${VAR:default}` 写法）

核心配置块包括：

- `server`（HTTP/gRPC、TLS、CORS）
- `data`（数据库、Redis、客户端）
- `registry` / `discovery` / `config`（治理与配置中心）
- `trace` / `metrics`（观测）

## 🔭 可观测性

项目默认集成观测组件（Compose 生产栈）：

- Grafana: `http://localhost:3001`
- Prometheus: `http://localhost:9090`
- Jaeger: `http://localhost:16686`
- Loki: `http://localhost:3100`
- OTel Collector: `4317/4318`

## 🧪 质量与约束

- 不要手动编辑生成代码（如 `api/gen/go/`、`wire_gen.go`、`openapi.yaml`）
- 修改 Proto 后务必执行 `make gen`
- 修改 Wire 配置后务必重新生成（`make wire` 或 `make gen`）

## 🤝 贡献

欢迎提交 Issue / PR。提交前请至少确保：

```bash
make lint
make test
```

## 📄 License

MIT，详见 `LICENSE`。
