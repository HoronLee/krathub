# 服务治理 PKG 包使用文档

本项目已将服务注册、服务发现和配置中心相关功能重构到 `pkg` 目录下，提供了通用的、可复用的组件。

## 目录结构

```
pkg/
├── registry/        # 服务注册组件
├── discovery/       # 服务发现组件
└── configCenter/    # 配置中心组件
```

## 1. 服务注册 (pkg/registry)

### 支持的注册中心
- **Consul**: 基于 HashiCorp Consul 的服务注册
- **Nacos**: 基于阿里巴巴 Nacos 的服务注册

### 使用示例

```go
import "krathub/pkg/registry"

// Consul 服务注册
consulRegistrar := registry.NewConsulRegistrar(&registry.ConsulConfig{
    Addr:       "127.0.0.1:8500",
    Scheme:     "http",
    Token:      "",
    Datacenter: "dc1",
    Timeout:    time.Second * 5,
})

// Nacos 服务注册
nacosRegistrar := registry.NewNacosRegistrar(&registry.NacosConfig{
    Addr:      "127.0.0.1",
    Port:      8848,
    Namespace: "public",
    Username:  "",
    Password:  "",
    Group:     "DEFAULT_GROUP",
    Timeout:   time.Second * 5,
})
```

## 2. 服务发现 (pkg/discovery)

### 支持的服务发现
- **Consul**: 基于 HashiCorp Consul 的服务发现
- **Nacos**: 基于阿里巴巴 Nacos 的服务发现

### 使用示例

```go
import "krathub/pkg/discovery"

// Consul 服务发现
consulDiscovery := discovery.NewConsulDiscovery(&discovery.ConsulConfig{
    Addr:       "127.0.0.1:8500",
    Scheme:     "http",
    Token:      "",
    Datacenter: "dc1",
    Timeout:    time.Second * 5,
})

// Nacos 服务发现
nacosDiscovery := discovery.NewNacosDiscovery(&discovery.NacosConfig{
    Addr:      "127.0.0.1",
    Port:      8848,
    Namespace: "public",
    Username:  "",
    Password:  "",
    Group:     "DEFAULT_GROUP",
    Timeout:   time.Second * 5,
})
```

## 3. 配置中心 (pkg/configCenter)

### 支持的配置中心
- **Nacos**: 基于阿里巴巴 Nacos 的配置管理 ✅ **完整实现**
- **Consul**: 基于 HashiCorp Consul 的配置管理 ✅ **完整实现**
- **Etcd**: 基于 Etcd 的配置管理 ✅ **完整实现**

### 使用示例

```go
import "krathub/pkg/configCenter"

// Nacos 配置中心
nacosConfigSource := configCenter.NewNacosConfigSource(&configCenter.NacosConfig{
    Addr:      "127.0.0.1",
    Port:      8848,
    Namespace: "public",
    Username:  "",
    Password:  "",
    Group:     "DEFAULT_GROUP",
    DataId:    "config.yaml",
    Timeout:   time.Second * 5,
})

// Consul 配置中心
consulConfigSource := configCenter.NewConsulConfigSource(&configCenter.ConsulConfig{
    Addr:       "127.0.0.1:8500",
    Scheme:     "http",
    Token:      "",
    Datacenter: "dc1",
    Key:        "krathub/config",
    Timeout:    time.Second * 5,
})

// Etcd 配置中心
etcdConfigSource := configCenter.NewEtcdConfigSource(&configCenter.EtcdConfig{
    Endpoints: []string{"127.0.0.1:2379"},
    Username:  "",
    Password:  "",
    Key:       "/krathub/config",
    Timeout:   time.Second * 5,
})

// 将配置源添加到 Kratos 配置中
config := config.New(
    config.WithSource(
        file.NewSource("config.yaml"),
        nacosConfigSource, // 或 consulConfigSource, etcdConfigSource
    ),
)
```

## 4. 配置中心特性对比

| 特性 | Nacos | Consul | Etcd |
|------|-------|--------|------|
| 配置存储 | ✅ DataId + Group | ✅ KV 存储 | ✅ KV 存储 |
| 配置监听 | ✅ 实时推送 | ✅ 长轮询 | ✅ Watch 机制 |
| 命名空间 | ✅ 多租户隔离 | ✅ 多数据中心 | ❌ |
| 认证授权 | ✅ 用户名/密码 | ✅ Token | ✅ 用户名/密码 |
| 配置格式 | ✅ 自动识别 | ✅ 自动识别 | ✅ 自动识别 |
| 集群支持 | ✅ | ✅ | ✅ |
| 配置版本 | ✅ 版本管理 | ✅ 修改索引 | ✅ 修订版本 |

## 5. 在其他项目中使用

这些 pkg 包可以被其他项目直接引用：

```go
import (
    "your-project/pkg/registry"
    "your-project/pkg/discovery"
    "your-project/pkg/configCenter"
)
```

## 6. 扩展指南

### 添加新的注册中心

1. 在 `pkg/registry/registry.go` 中添加新的配置结构
2. 实现相应的 `NewXXXRegistrar` 函数
3. 返回实现了 `registry.Registrar` 接口的客户端

### 添加新的服务发现

1. 在 `pkg/discovery/discovery.go` 中添加新的配置结构
2. 实现相应的 `NewXXXDiscovery` 函数
3. 返回实现了 `registry.Discovery` 接口的客户端

### 添加新的配置中心

1. 在 `pkg/configCenter/configCenter.go` 中添加新的配置结构
2. 实现相应的 `NewXXXConfigSource` 函数
3. 返回实现了 `config.Source` 接口的配置源

## 7. 优势

- **通用性**: 可以在不同项目间复用
- **解耦**: 业务逻辑与基础设施组件分离
- **可扩展**: 易于添加新的注册中心、服务发现和配置中心支持
- **标准化**: 统一的接口和配置方式
- **可测试**: 独立的包便于单元测试
- **生产就绪**: 完整实现了三大主流配置中心支持

## 8. 详细文档

- [配置中心详细使用指南](./configCenter/README.md)