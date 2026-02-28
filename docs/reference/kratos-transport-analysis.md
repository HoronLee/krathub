# kratos-transport 项目分析（面向 servora）

## 1. 项目定位与价值

`kratos-transport` 的核心目标是：把“非 Kratos 原生传输协议/消息系统”统一适配为 `kratos transport.Server`，从而纳入同一生命周期管理。

典型覆盖：

- HTTP/RPC 扩展（gin、fasthttp、graphql、thrift、trpc、hertz、iris）
- 消息队列传输（kafka、nats、nsq、rabbitmq、rocketmq、mqtt、redis、pulsar 等）
- 实时协议（websocket、sse、tcp、socketio、signalr、webtransport、mcp）

关键目录：

- 传输实现：`/Users/horonlee/projects/go/kratos-transport/transport`
- Broker 抽象：`/Users/horonlee/projects/go/kratos-transport/broker`

---

## 2. 核心抽象与实现模式

## 2.1 统一生命周期契约

几乎所有传输实现都对齐 `Start(context.Context) error` / `Stop(context.Context) error`，并实现 `Endpoint() (*url.URL, error)`。

关键模式：

- `listenAndEndpoint` 负责监听器与 endpoint 初始化
- `Start` 做幂等检查、启动主服务
- `Stop` 做优雅关闭与资源清理

参考示例：

- gin：`/Users/horonlee/projects/go/kratos-transport/transport/gin/server.go`
- kafka：`/Users/horonlee/projects/go/kratos-transport/transport/kafka/server.go`
- keepalive：`/Users/horonlee/projects/go/kratos-transport/transport/keepalive/server.go`

## 2.2 Option 模式驱动配置注入

各 transport 子模块普遍采用 `NewServer(opts ...ServerOption)`：

- 传入网络地址、TLS、codec、middleware、tracing 等
- 通过 `init(opts...)` 完成内部组装

参考：

- gin option：`/Users/horonlee/projects/go/kratos-transport/transport/gin/options.go`
- kafka option：`/Users/horonlee/projects/go/kratos-transport/transport/kafka/options.go`

## 2.3 Broker 层统一消息抽象

`broker.Broker` 抽象统一了不同消息系统的操作：

- `Connect/Disconnect`
- `Publish/Subscribe/Request`

并搭配 `SubscribeOption/PublishOption` 与 tracing 选项做跨实现一致性。

参考：

- `/Users/horonlee/projects/go/kratos-transport/broker/broker.go`
- `/Users/horonlee/projects/go/kratos-transport/broker/options.go`

## 2.4 Keepalive 旁路健康服务

对 MQ 类 transport（外部无法直接探活）增加 keepalive gRPC health server，单独暴露可注册 endpoint，解决“消息 broker 服务本体不可探活”的治理问题。

参考：

- `/Users/horonlee/projects/go/kratos-transport/transport/keepalive/server.go`
- `/Users/horonlee/projects/go/kratos-transport/transport/keepalive/README.md`

## 2.5 可观测性注入点

主要体现在两类注入：

- HTTP/RPC 层：如 gin 使用 `otelgin` 中间件
- MQ/Broker 层：通过 broker option 注入 tracer provider / propagator

参考：

- gin tracing：`/Users/horonlee/projects/go/kratos-transport/transport/gin/options.go`
- broker tracing option：`/Users/horonlee/projects/go/kratos-transport/broker/options.go`

---

## 3. 对 servora 的可迁移建议（优先级）

## P0：定义统一的扩展传输接口层

在 servora 内建议抽象一个 `pkg/transport/serverext`（命名可调整），明确：

- 必须实现 Start/Stop/Endpoint
- 必须支持 Option 注入
- 必须定义统一日志字段和错误语义

先做最小闭环：`sse` 与 `websocket` 两个高价值扩展。

## P0：引入 keepalive 模式用于 MQ Worker

如果未来 servora 引入异步消费服务（Kafka/NSQ/RabbitMQ Worker），建议同步引入 keepalive 旁路探活，确保注册中心健康状态可信。

## P1：沉淀“transport 选项规范”

约定所有扩展 transport 的 Option 至少覆盖：

- `WithAddress`
- `WithTLSConfig`
- `WithMiddleware`（如适用）
- `WithTracerProvider/WithPropagator`

避免每个实现自定义杂乱参数。

## P1：统一 Endpoint 生成策略

借鉴 `AdjustAddress + NewRegistryEndpoint` 思路，统一“监听地址到注册地址”的转换，避免多网卡/0.0.0.0 场景注册错误。

参考：`/Users/horonlee/projects/go/kratos-transport/transport/utils.go`。

## P2：MQ 抽象可分阶段引入

`kratos-transport` 的 MQ 生态很全，但 servora 建议分阶段：

1. 先统一 broker 接口与 tracing/coding 选项
2. 首批接入 1~2 个队列（如 Kafka + Redis stream）
3. 再逐步扩展其他 broker

---

## 4. 与 servora 现状的对照

servora 已具备：

- 统一 bootstrap runtime：`/Users/horonlee/projects/go/servora/pkg/bootstrap/bootstrap.go`
- HTTP/gRPC 中间件链：`/Users/horonlee/projects/go/servora/app/servora/service/internal/server/http.go`、`/Users/horonlee/projects/go/servora/app/servora/service/internal/server/grpc.go`

下一步差距主要在：

- 扩展 transport 的统一规范与模板
- MQ worker 的可观测/可治理能力（keepalive）
- 非 gRPC/HTTP 协议（如 websocket/sse/mcp）的标准接入

---

## 5. 外部最佳实践校准（Kratos 官方/社区）

可用于校准 transport 设计的外部原则：

- 中间件顺序敏感：`recovery -> tracing -> logging -> ...`
- OTel 作为 tracing 标准，日志字段应包含 `trace_id/span_id`
- 服务生命周期统一通过 Server Start/Stop 管理

参考：

- https://pkg.go.dev/github.com/go-kratos/kratos/v2/middleware
- https://go-kratos.dev/docs/component/middleware/tracing
- https://go-kratos.dev/docs/component/transport/overview

---

## 6. 落地路线建议（可执行）

1. 在 servora 建立 `transport extension RFC`（接口、Option、日志、endpoint 规范）
2. 落地 `sse` 与 `websocket` 两个示范实现
3. 为异步 worker 增加 keepalive server 并接入 registry
4. 把 tracing/metrics/error mapping 模板化，沉淀到 `pkg/transport`
5. 补齐单元测试 + 集成测试（至少覆盖 Start/Stop 幂等、endpoint 正确性）

如果按此路线推进，可以在不破坏现有 gRPC/HTTP 主链路的情况下，逐步把 servora 升级为“多传输协议统一治理”的框架。
