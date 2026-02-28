# kratos-bootstrap 项目分析（面向 micro-forge）

## 1. 项目定位与价值

`kratos-bootstrap` 是一个偏“启动层/基础设施装配层”的 Kratos 扩展库，目标是把以下能力标准化：

- 应用启动入口与生命周期封装
- 配置加载（本地 + 远程）
- 日志、注册中心、追踪初始化
- 组件工厂注册与按配置动态创建

对应源码主线可参考：

- 启动主链路：`/Users/horonlee/projects/go/kratos-bootstrap/bootstrap/bootstrap.go`
- 配置装载：`/Users/horonlee/projects/go/kratos-bootstrap/config/config.go`
- 配置工厂注册：`/Users/horonlee/projects/go/kratos-bootstrap/config/factory.go`
- 注册中心工厂：`/Users/horonlee/projects/go/kratos-bootstrap/registry/discovery.go`、`/Users/horonlee/projects/go/kratos-bootstrap/registry/registrar.go`

---

## 2. 核心设计模式

## 2.1 启动链路收敛（Bootstrap Pipeline）

`RunApp -> bootstrap -> initApp -> app.Run` 的固定链路把初始化时序固化下来：

1. 打印应用信息
2. 加载配置
3. 初始化 logger
4. 初始化 registrar
5. 初始化 tracer
6. 调用业务方 `initApp` 构建 kratos.App
7. 运行 app 并处理退出

优点：

- 新服务只关心 `initApp` 业务装配，基础设施初始化不再重复写。
- 初始化故障在早期阶段失败（fail fast），降低线上半初始化状态风险。

关键参考：`/Users/horonlee/projects/go/kratos-bootstrap/bootstrap/bootstrap.go`。

## 2.2 工厂注册 + 插件化扩展

配置中心与注册中心都采用“`Type -> Factory`”注册表模式：

- 配置源：`MustRegisterFactory/RegisterFactory/NewProvider`
- 服务发现：`RegisterDiscoveryFactory/NewDiscovery`
- 服务注册：`RegisterRegistrarFactory/NewRegistrar`

这使新增实现（如 Nacos、Polaris）时只需新增子包并在 `init` 注册，无需改主流程。

关键参考：

- `/Users/horonlee/projects/go/kratos-bootstrap/config/factory.go`
- `/Users/horonlee/projects/go/kratos-bootstrap/registry/discovery.go`
- `/Users/horonlee/projects/go/kratos-bootstrap/registry/registrar.go`

## 2.3 配置聚合与可扩展扫描

`config/bootstrap_config.go` 通过 `RegisterConfig` 维护需要扫描的配置列表，`LoadBootstrapConfig` 统一 `cfg.Scan`。模式上支持：

- 内置配置结构（server/data/logger/trace/...）
- 业务方追加自定义 proto 配置

关键参考：

- `/Users/horonlee/projects/go/kratos-bootstrap/config/bootstrap_config.go`
- `/Users/horonlee/projects/go/kratos-bootstrap/config/config.go`

## 2.4 注册名规范化

对不同注册中心（如 Consul/Kubernetes）做服务名归一化，避免非法字符和长度问题。

关键参考：`/Users/horonlee/projects/go/kratos-bootstrap/registry/normalize.go`。

## 2.5 CLI 与守护进程化入口

通过 cobra/flags 统一注入 `conf/env/daemon` 等参数，为后续子命令扩展预留空间。

关键参考：

- `/Users/horonlee/projects/go/kratos-bootstrap/bootstrap/flag.go`
- `/Users/horonlee/projects/go/kratos-bootstrap/bootstrap/bootstrap.go`

---

## 3. 对 micro-forge 的可迁移建议（优先级）

## P0：形成统一 Runtime 启动骨架

当前 micro-forge 已有 `pkg/bootstrap`（`/Users/horonlee/projects/go/micro-forge/pkg/bootstrap/bootstrap.go`），建议继续收敛为固定阶段：

1. 配置加载
2. 日志初始化
3. 注册中心初始化
4. Trace/Metrics 初始化
5. Wire 组装 app
6. app.Run

并把每阶段的失败日志标准化，便于故障定位。

## P0：抽象“可插拔配置源工厂”

现状中配置加载与治理能力分散在不同包。建议引入与 `kratos-bootstrap/config/factory.go` 等价的统一注册表，覆盖：

- 本地文件
- etcd/consul/nacos/kubernetes

收益：新增配置源无需改核心启动代码。

## P1：抽象“注册中心工厂”并统一命名规范

在 `pkg/governance/registry` 现有基础上，新增统一 `NewRegistrar/NewDiscovery` 工厂入口，并沿用服务名规范化策略（特别是 K8s DNS label）。

## P1：配置扫描扩展点

参考 `RegisterConfig` 模式，为业务模块提供注册自定义配置 proto 的机制，减少在 `main/wire` 里手工穿透配置结构。

## P2：CLI 子命令能力

可在 `cmd/server` 增加子命令约定（如 `serve`, `version`, `doctor`），并复用统一 flags，提升运维可用性。

---

## 4. 结合 Kratos 官方实践的校准

与官方建议一致的点：

- 配置与代码分离（运行期加载）
- Wire 做显式初始化，避免全局变量
- 中间件顺序敏感（recovery/tracing/logging）

建议 micro-forge 在后续文档中明确“生产配置不打包镜像”的约束，并把 Runtime 初始化阶段作为模板固化。

参考：

- https://go-kratos.dev/docs/component/config/
- https://go-kratos.dev/docs/guide/wire
- https://pkg.go.dev/github.com/go-kratos/kratos/v2/middleware

---

## 5. 风险与注意事项

- `kratos-bootstrap` 本身集成面较大，micro-forge 建议先迁移“模式”而非一次性搬全部组件。
- 工厂注册表需要避免重复注册 panic 对启动稳定性的影响（建议启动期统一检查并输出已注册列表）。
- 命名规范化要兼顾历史服务名兼容策略，避免注册名突变引发流量切换问题。
