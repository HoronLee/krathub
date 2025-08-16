# 配置中心使用示例

## Consul 配置中心

### 1. 配置文件示例 (config.yaml)

```yaml
config:
  consul:
    addr: "127.0.0.1:8500"
    scheme: "http"
    token: ""
    datacenter: "dc1"
    key: "krathub/config"  # 配置在 Consul KV 中的键名
    timeout: "5s"
```

### 2. 在 Consul 中存储配置

```bash
# 使用 Consul CLI 存储配置
consul kv put krathub/config @config.yaml

# 或者使用 HTTP API
curl -X PUT \
  http://127.0.0.1:8500/v1/kv/krathub/config \
  --data-binary @config.yaml
```

### 3. 编程方式使用

```go
import "krathub/pkg/configCenter"

consulConfigSource := configcenter.NewConsulConfigSource(&configcenter.ConsulConfig{
    Addr:       "127.0.0.1:8500",
    Scheme:     "http",
    Token:      "",
    Datacenter: "dc1",
    Key:        "krathub/config",
    Timeout:    time.Second * 5,
})

config := config.New(
    config.WithSource(
        file.NewSource("config.yaml"),
        consulConfigSource,
    ),
)
```

## Etcd 配置中心

### 1. 配置文件示例 (config.yaml)

```yaml
config:
  etcd:
    endpoints:
      - "127.0.0.1:2379"
      - "127.0.0.1:2380"
    username: ""
    password: ""
    key: "/krathub/config"  # 配置在 Etcd 中的键名
    timeout: "5s"
```

### 2. 在 Etcd 中存储配置

```bash
# 使用 etcdctl 存储配置
etcdctl put /krathub/config "$(cat config.yaml)"

# 或者设置多个配置项
etcdctl put /krathub/config/database.yaml "$(cat database.yaml)"
etcdctl put /krathub/config/redis.yaml "$(cat redis.yaml)"
```

### 3. 编程方式使用

```go
import "krathub/pkg/configCenter"

etcdConfigSource := configcenter.NewEtcdConfigSource(&configcenter.EtcdConfig{
    Endpoints: []string{"127.0.0.1:2379", "127.0.0.1:2380"},
    Username:  "",
    Password:  "",
    Key:       "/krathub/config",
    Timeout:   time.Second * 5,
})

config := config.New(
    config.WithSource(
        file.NewSource("config.yaml"),
        etcdConfigSource,
    ),
)
```

## Nacos 配置中心

### 1. 配置文件示例 (config.yaml)

```yaml
config:
  nacos:
    addr: "127.0.0.1"
    port: 8848
    namespace: "public"
    data_id: "krathub.yaml"
    group: "DEFAULT_GROUP"
    username: ""
    password: ""
    timeout: "5s"
```

### 2. 编程方式使用

```go
import "krathub/pkg/configCenter"

nacosConfigSource := configcenter.NewNacosConfigSource(&configcenter.NacosConfig{
    Addr:      "127.0.0.1",
    Port:      8848,
    Namespace: "public",
    Username:  "",
    Password:  "",
    Group:     "DEFAULT_GROUP",
    DataId:    "krathub.yaml",
    Timeout:   time.Second * 5,
})

config := config.New(
    config.WithSource(
        file.NewSource("config.yaml"),
        nacosConfigSource,
    ),
)
```

## 完整的配置加载示例

```go
package main

import (
    "krathub/internal/conf"
    "krathub/pkg/configCenter"
    
    "github.com/go-kratos/kratos/v2/config"
    "github.com/go-kratos/kratos/v2/config/env"
    "github.com/go-kratos/kratos/v2/config/file"
)

func loadConfig() (*conf.Bootstrap, config.Config, error) {
    // 创建基本配置源
    c := config.New(
        config.WithSource(
            env.NewSource("KRATHUB_"),
            file.NewSource("configs/config.yaml"),
        ),
    )

    // 加载基本配置
    if err := c.Load(); err != nil {
        return nil, nil, err
    }

    // 扫描基本配置到结构体
    var bc conf.Bootstrap
    if err := c.Scan(&bc); err != nil {
        return nil, nil, err
    }

    // 检查是否配置了远程配置中心
    if configSrc := bc.Config; configSrc != nil {
        var configSource config.Source

        switch configType := configSrc.Config.(type) {
        case *conf.Config_Nacos_:
            configSource = configcenter.NewNacosConfigSource(&configcenter.NacosConfig{
                Addr:      configType.Nacos.Addr,
                Port:      configType.Nacos.Port,
                Namespace: configType.Nacos.Namespace,
                Username:  configType.Nacos.Username,
                Password:  configType.Nacos.Password,
                Group:     configType.Nacos.Group,
                DataId:    configType.Nacos.DataId,
                Timeout:   configType.Nacos.Timeout,
            })
        case *conf.Config_Consul_:
            configSource = configcenter.NewConsulConfigSource(&configcenter.ConsulConfig{
                Addr:       configType.Consul.Addr,
                Scheme:     configType.Consul.Scheme,
                Token:      configType.Consul.Token,
                Datacenter: configType.Consul.Datacenter,
                Key:        configType.Consul.Key,
                Timeout:    configType.Consul.Timeout,
            })
        case *conf.Config_Etcd_:
            configSource = configcenter.NewEtcdConfigSource(&configcenter.EtcdConfig{
                Endpoints: configType.Etcd.Endpoints,
                Username:  configType.Etcd.Username,
                Password:  configType.Etcd.Password,
                Key:       configType.Etcd.Key,
                Timeout:   configType.Etcd.Timeout,
            })
        }

        if configSource != nil {
            // 创建新的配置对象，包含远程配置源
            newConfig := config.New(
                config.WithSource(
                    env.NewSource("KRATHUB_"),
                    file.NewSource("configs/config.yaml"),
                    configSource,
                ),
            )

            // 替换配置对象
            c.Close()
            c = newConfig

            // 重新加载配置
            if err := c.Load(); err != nil {
                return nil, nil, err
            }

            // 重新扫描配置到结构体
            if err := c.Scan(&bc); err != nil {
                return nil, nil, err
            }
        }
    }

    return &bc, c, nil
}
```

## 特性说明

### Consul 配置中心特性
- ✅ 支持 HTTP/HTTPS 连接
- ✅ 支持 Token 认证
- ✅ 支持多数据中心
- ✅ 支持自定义键名
- ✅ 支持配置变更监听

### Etcd 配置中心特性  
- ✅ 支持多节点集群
- ✅ 支持用户名/密码认证
- ✅ 支持前缀匹配查询
- ✅ 支持配置变更监听
- ✅ 自动识别配置格式 (yaml/json/toml)

### Nacos 配置中心特性
- ✅ 支持命名空间隔离
- ✅ 支持分组管理
- ✅ 支持用户名/密码认证  
- ✅ 支持配置变更监听
- ✅ 支持配置版本管理

## 最佳实践

1. **配置优先级**: 环境变量 > 远程配置中心 > 本地文件
2. **容错处理**: 远程配置中心不可用时，使用本地文件兜底
3. **配置监听**: 利用配置变更监听实现配置热更新
4. **安全考虑**: 敏感配置使用认证和加密传输
5. **键名规范**: 使用统一的键名规范，便于管理和维护