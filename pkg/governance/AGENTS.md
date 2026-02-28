<!-- Parent: ../AGENTS.md -->
# 服务治理 (pkg/governance)

**最后更新时间**: 2026-02-09

## 模块目的
封装服务发现、注册以及配置中心的功能。支持多种后端：Etcd, Consul, Nacos, Kubernetes。

## 核心组件

### 1. 服务注册与发现 (registry)
实现 Kratos `registry.Registrar` 和 `registry.Discovery` 接口。
- `etcd.go`: Etcd 实现，包含心跳检测和自动重试机制。
- `kubernetes.go`: Kubernetes 原生服务发现。

### 2. 配置中心 (configCenter)
提供统一的配置加载适配器。

## 使用示例

### Etcd 注册中心初始化
```go
import "github.com/horonlee/servora/pkg/governance/registry"

reg, err := registry.NewEtcdRegistry(etcdConfig)
if err != nil {
    panic(err)
}
```

## 测试指南
由于涉及外部依赖，部分测试可能需要特定环境：
```bash
go test -v ./pkg/governance/registry/...
```
🤖 Generated with [Claude Code](https://claude.com/claude-code)
