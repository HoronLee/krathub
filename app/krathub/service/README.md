# Krathub Service

Krathub 是本项目的主示例微服务，展示了基于 Kratos 框架的完整微服务实现。

## 技术选型

- **ORM**: GORM + GORM GEN（自动生成 PO 和 DAO）
- **缓存**: Redis
- **认证**: JWT
- **API**: gRPC + HTTP 双协议

## 目录结构

```
.
├── cmd/
│   ├── server/          # 服务启动入口
│   └── genDao/          # GORM GEN 代码生成工具
├── configs/             # 配置文件
├── internal/
│   ├── biz/            # 业务逻辑层
│   ├── data/           # 数据访问层
│   │   ├── dao/        # GORM GEN 生成的 DAO
│   │   └── po/         # GORM GEN 生成的 PO
│   ├── server/         # gRPC/HTTP 服务器配置
│   └── service/        # Service 层实现
└── Makefile
```

## 开发命令

```shell
# 生成 GORM GEN 的 PO 和 DAO 代码
make genDao

# 生成 wire 依赖注入代码
make wire

# 运行服务
make run

# 构建服务
make build
```

## GORM GEN 使用说明

本服务使用 GORM GEN 自动生成数据访问层代码：

1. 确保数据库配置正确（`configs/config.yaml`）
2. 运行 `make genDao` 生成代码
3. 生成的 PO 位于 `internal/data/po/`
4. 生成的 DAO 位于 `internal/data/dao/`

生成工具会连接数据库，为所有表自动生成对应的模型和查询方法。

## 配置

复制示例配置并修改：

```shell
cp configs/config-example.yaml configs/config.yaml
```

主要配置项：
- 数据库连接（MySQL/PostgreSQL/SQLite）
- Redis 连接
- JWT 密钥
- 服务端口
