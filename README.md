# Krathub

[English](README.en-US.md) | 简体中文

Krathub 是一个基于 Go Kratos 框架的微服务项目模板。它集成了一系列最佳实践和常用组件，旨在帮助开发者快速构建一个功能完善、结构清晰、易于扩展的现代化 Go 应用。

## ✨ 核心特性

- **微服务架构**: 基于 Kratos v2 构建，天然支持微服务。
- **双协议支持**: 同时提供 gRPC 和 HTTP 接口，满足不同场景需求。
- **服务治理**: 集成 Consul 和 Nacos，提供开箱即用的服务注册与发现能力。
- **配置中心**: 支持通过 Consul 或 Nacos 进行动态配置管理。
- **数据库集成**: 采用 GORM 作为 ORM，并提供 `make gen.db` 快速生成模型代码。
- **依赖注入**: 使用 Wire 进行依赖注入，清晰化对象之间的依赖关系。
- **代码生成**: 大量使用 `make` 命令简化 `proto`、`wire` 等代码的生成。
- **认证鉴权**: 内置 JWT 中间件，方便实现用户认证。
- **容器化**: 提供 `Dockerfile` 和 `docker-compose.yml`，轻松实现容器化部署。
- **可观测性**: 已集成 `Metrics` (Prometheus) 和 `Trace` (Jaeger) 的基础配置。

## 📖 API 文档

您可以在以下地址查看并测试项目的 API：

- **[https://jqovxjvrtw.apifox.cn](https://jqovxjvrtw.apifox.cn)**

## 🚀 快速开始

请确保您已安装 Go、Docker 以及 `make` 工具。

1.  **克隆项目**
    ```bash
    git clone https://github.com/HoronLee/krathub.git
    cd krathub
    ```

2.  **安装依赖工具**
    此命令将安装 `protoc` 插件、`kratos` 工具、`wire` 等开发依赖。
    ```bash
    make init
    ```

3.  **生成所有代码**
    此命令会清理旧文件、生成 `proto`、数据库模型、`wire` 依赖注入代码等。
    ```bash
    make all
    ```

4.  **配置项目**
    - 根据您的环境修改 `configs/config.yaml` 中的数据库、Redis、Consul/Nacos 等配置。

5.  **运行项目**
    ```bash
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

## 🛠️ 常用命令

项目通过 `Makefile` 提供了丰富的命令来简化开发流程。

- `make help`: 显示所有可用的 make 命令。
- `make init`: 初始化开发环境，安装所需工具。
- `make proto`: 生成所有 Protobuf 相关代码 (api, errors, config)。
- `make gen.db`: 根据 `configs/config.yaml` 中的数据库配置生成 GORM 模型。
- `make wire`: 在 `cmd/server/` 目录下运行 `wire` 生成依赖注入代码。
- `make all`: 清理并执行所有代码生成和构建任务。
- `make run`: 启动服务。
- `make build`: 编译和构建二进制文件到 `bin/` 目录。
- `make clean`: 清理所有生成的文件和构建产物。

## 📝 开发流程

推荐的开发顺序如下，以确保依赖关系正确：

1.  **API 定义 (`api/`)**: 在 `.proto` 文件中定义 gRPC 服务和消息体。
2.  **生成代码 (`make proto`)**: 运行命令生成 gRPC、HTTP、Errors 的桩代码。
3.  **业务逻辑 (`internal/biz/`)**: 定义业务逻辑的接口和实现，这是不依赖任何框架的核心。
4.  **数据访问 (`internal/data/`)**: 实现 `biz` 层定义的接口，负责与数据库、缓存等交互。
5.  **服务实现 (`internal/service/`)**: 实现 `api` 层定义的 gRPC 服务接口，它会调用 `biz` 层的逻辑。
6.  **依赖注入 (`cmd/server/wire.go`)**: 将新的 `service`, `biz`, `data` 组件注入到 `wire.go` 中。
7.  **运行 `make wire`**: 生成最终的依赖注入代码。
8.  **启动与测试**: 运行 `make run` 并进行测试。

## 📞 客户端层 (Client)

`internal/client` 是一个自定义的层，用于管理对外部 gRPC 等服务的调用。它提供了一个客户端工厂，可以方便地基于服务发现创建和复用连接，专为服务间的通信而设计。

## ⚙️ 配置文件示例

以下是 `configs/config.yaml` 的一个示例，展示了所有可用的配置项。您可以使用环境变量覆盖这些值。

```yaml
server:
  http:
    addr: "${HADDR:0.0.0.0:8000}"
    timeout: "${HTIMEOUT:1s}"
    tls:
      enable: false
      cert_path: "${HTTPS_CERT_PATH:../../manifest/certs/server.cert}"
      key_path: "${HTTPS_KEY_PATH:../../manifest/certs/server.key}"
  grpc:
    addr: "${GADDR:0.0.0.0:8001}"
    timeout: "${GTIMEHOUT:1s}"
    tls:
      enable: false
      cert_path: "${GRPCS_CERT_PATH:../../manifest/certs/server.cert}"
      key_path: "${GRPCS_KEY_PATH:../../manifest/certs/server.key}"

data:
  database:
    driver: "${DB_DRIVER:mysql}"
    source: "${DB_SOURCE:projectName:123456@tcp(127.0.0.1:3306)/projectName?parseTime=True&loc=Local}"
  client:
    # 这里可以配置第三方服务的客户端
    # 默认不用配置，而是在代码中直接服务发现
    grpc:
      # - service_name: hello.grpc  # nacos需要添加协议后缀
      #   endpoint: "${GRPC_ENDPOINT:127.0.0.1:8003}"
      #   enable_tls: false
      #   timeout: 5s

app:
  env: "${ENV:dev}" # dev, test, prod
  jwt:
    secret_key: "${JWT_SECRETK_KEY:projectName_secret}"
    expire: "${JWT_EXPIRE:24}"
    issuer: "${JWT_ISSUER:projectName}"
    # audience: "${JWT_AUDIENCE:projectName}"
  log:
    level: "${LOG_LEVEL:-1}"  # -1debug,0info,1warn,2error,3dpanic,4panic,5fatal
    filename: "${LOG_FILENAME:projectName.log}"  # 日志文件夹为根目录logs
    max_size: "${LOG_MAX_SIZE:20}"  # 日志文件最大大小，单位MB
    max_age: "${LOG_MAX_AGE:30}"  # 日志文件最大保存天数
    max_backups: "${LOG_MAX_BACKUPS:10}"  # 日志文件最大备份数

# 注册中心配置
registry:
  # 使用 Consul 作为注册中心
  consul:
    addr: consul.r430.com:30080
    scheme: http
    datacenter: dc1
    timeout: 5s
    # 为服务自定义 tag
    tags:
     - "traefik.enable=true"
     - "traefik.http.routers.krathub-router.rule=Host(`krathub.r430.com`)"
     - "traefik.http.services.krathub-svc.loadbalancer.server.port=8000"

  # 或者使用 Nacos 作为注册中心
  # nacos:
    # addr: "${NACOSR_ADDR:127.0.0.1}"
    # port: "${NACOSR_PORT:8848}"
    # namespace: "${NACOSR_NAMESPACE:public}"
    # group: "${NACOSR_GROUP:DEFAULT_GROUP}"
    # username: "${NACOSR_USERNAME:nacos}"
    # password: "${NACOSR_PASSWORD:nacos}"
    # timeout: "${NACOSR_TIMEOUT:5s}"

# 服务发现配置，一般和注册中心配置相同
discovery:
  consul:
    addr: consul.r430.com:30080
    scheme: http
    datacenter: dc1
    timeout: 5s

  # nacos:
  #   addr: "${NACOSD_ADDR:127.0.0.1}"
  #   port: "${NACOSD_PORT:8848}"
  #   namespace: "${NACOSD_NAMESPACE:public}"
  #   group: "${NACOSD_GROUP:DEFAULT_GROUP}"
  #   username: "${NACOSD_USERNAME:nacos}"
  #   password: "${NACOSD_PASSWORD:nacos}"
  #   timeout: "${NACOSD_TIMEOUT:5s}"

# 配置中心配置
config:

  # 使用 Nacos 作为配置中心
  # nacos:
  #   addr: "${NACOSC_ADDR:127.0.0.1}"
  #   port: "${NACOSC_PORT:8848}"
  #   namespace: "${NACOSC_NAMESPACE:public}"
  #   group: "${NACOSC_GROUP:DEFAULT_GROUP}"
  #   data_id: "${NACOSC_DATA_ID:projectName.yaml}"
  #   username: "${NACOSC_USERNAME:nacos}"
  #   password: "${NACOSC_PASSWORD:nacos}"
  #   timeout: "${NACOSC_TIMEOUT:5s}"

  # 使用 consul 作为配置中心
  # consul:
  #   addr: consul.r430.com:30080
  #   scheme: http
  #   datacenter: dc1
  #   timeout: 5s
  #   key: projectName/config.yaml  # 配置键名

# 链路追踪配置
trace:
  # 使用 Jaeger 作为链路追踪
  endpoint: otlp.jaeger.r430.com:30080
# 监控配置
metrics:
    # 当前只支持Prometheus
    prometheus:
        addr: ":8000"
    meterName: "krathub"
```

## 📄 许可协议

本项目遵循 [LICENSE](LICENSE) 文件中的许可协议。