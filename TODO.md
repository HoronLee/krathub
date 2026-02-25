# Krathub TODO

## Kubernetes 增强功能

### 高优先级

- [ ] **HTTP 健康检查端点 `/healthz`** - 让 K8s 可用 HTTP 探针而非 gRPC，提高兼容性
- [ ] **优雅关闭时服务注销** - 确保服务下线时从注册中心移除，避免流量打到已停止的 Pod

### 中优先级

- [ ] **HorizontalPodAutoscaler (HPA)** - 根据 CPU/内存自动扩缩容
- [ ] **Ingress 配置** - 通过域名暴露 HTTP 服务，支持 TLS

### 低优先级

- [ ] **Prometheus 监控** - 添加 ServiceMonitor + 指标采集，接入 Grafana 看板

## 业务功能

_待补充_

## 技术债务

- [x] **CORS 中间件简化** - 删除 `app/krathub/service/internal/server/middleware/cors.go` 封装层，让 `pkg/middleware/cors` 的 `Middleware()` 方法直接接受 `*conf.CORS` (来自 `api/gen/go/conf/v1/conf.pb.go`)，所有空字段检测、默认值填充等操作在 pkg 内部完成，`http.go` 可直接调用 `cors.Middleware(c.Http.Cors)`
- [x] **PostgreSQL init.sql 外置** - 优化 `app/krathub/service/deployment/kubernetes/postgres.yaml`，将 `init.sql: |` 内联字段外置为独立文件（使用 ConfigMap 挂载或 kustomize configMapGenerator）
- [ ] **gRPC 客户端配置查找优化** - `pkg/transport/client/grpc_conn.go` 的 `createGrpcConnection` 每次建连都遍历 `dataCfg.Client.GetGrpc()`；评估在 `NewClient` 启动阶段构建 `service_name -> config` 索引缓存，避免热路径线性扫描

## Oracle 审查记录（2026-02-26）

- 结论：**无 Critical**
- 审查摘要：当前 observability 迁移可运行，但存在安全与生产稳态风险，建议优先收敛以下事项。

### 待处理（高优先级）

- [ ] **Grafana 默认弱口令 + 对宿主机暴露** - `docker-compose.yaml` 中 `GF_SECURITY_ADMIN_PASSWORD` 默认回退为 `admin` 且 `3001:3000` 暴露
- [ ] **Loki 未鉴权 + 对外暴露** - `manifests/loki/loki.yaml` 为 `auth_enabled: false`，compose 暴露 `3100`
- [ ] **Docker 配置存在明文凭据** - `app/krathub/service/configs/config-docker.yaml` 中 DB/Redis/JWT 明文
- [ ] **OTel Collector 以 root 身份运行** - `docker-compose.yaml` 中 `user: "0:0"`

### 待处理（中优先级）

- [ ] **Tracing 明文传输 + 全采样** - `pkg/governance/telemetry/tracing.go` 使用 `WithInsecure()` 且采样率 `1.0`
- [ ] **Collector debug exporter 在生产链路** - `manifests/otel/otel-collector.yaml` 的 traces/logs pipeline 包含 `debug`
- [ ] **Collector 未纳入健康依赖链** - 已配置 health extension，但服务依赖仍是 `service_started`
- [ ] **sayhello metrics 可观测性闭环不完整** - 已接入 middleware，Prometheus 当前仅抓取 `krathub`

### 待处理（低优先级）

- [ ] **Prometheus 抓取范围偏窄** - 建议补充 otel/loki/jaeger/grafana 的组件级采集
