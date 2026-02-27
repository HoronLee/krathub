# micro-forge TODO

## Kubernetes 增强功能

- [ ] **HTTP 健康检查端点 `/healthz`** - 让 K8s 可用 HTTP 探针而非 gRPC，提高兼容性

## 技术债务

- [ ] **gRPC 客户端配置查找优化** - `pkg/transport/client/grpc_conn.go` 的 `createGrpcConnection` 每次建连都遍历 `dataCfg.Client.GetGrpc()`；评估在 `NewClient` 启动阶段构建 `service_name -> config` 索引缓存，避免热路径线性扫描

## Oracle 审查记录（2026-02-26）

- [ ] **Tracing 明文传输 + 全采样** - `pkg/governance/telemetry/tracing.go` 使用 `WithInsecure()` 且采样率 `1.0`
- [ ] **Collector debug exporter 在生产链路** - `manifests/otel/otel-collector.yaml` 的 traces/logs pipeline 包含 `debug`
- [ ] **Collector 未纳入健康依赖链** - 已配置 health extension，但服务依赖仍是 `service_started`
- [ ] **sayhello metrics 可观测性闭环不完整** - 已接入 middleware，Prometheus 当前仅抓取 `micro-forge`
- [ ] **Prometheus 抓取范围偏窄** - 建议补充 otel/loki/jaeger/grafana 的组件级采集
