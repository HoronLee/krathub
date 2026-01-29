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

_待补充_
