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
