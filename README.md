# Krathub

> 基于Kratos框架编写的快开框架，目前处于开发初期阶段

开发顺序: api -> config -> service -> biz -> data
功能编写完成后需要依次进行依赖注入然后wire生成代码，并且需要在internal/server中添加用法

## 项目依赖

直接执行`make init`即可下载所需软件，下方仅为新增以来软件的列表

项目中需要使用wire依赖注入,安装wire:

```bash
go install github.com/google/wire/cmd/wire@latest
```

protobuf 参数校验插件

```bash
go install github.com/envoyproxy/protoc-gen-validate@latest
```

protobuf 错误处理插件

```bash
go install github.com/envoyproxy/protoc-gen-validate@latest
```

## Data层编码须知

编写 data 层代码之前需要先修改configs目录下的config.yaml文件来配置数据库等相关信息。然后再通过`make gendb`来生成 orm 代码

## Docker Compose 部署

可用的docker compose文件在项目的deployment/docker-compose目录下，首次运行请把model.sql放入initdb文件夹中，这样数据库首次运行就会导入数据。配置文件放于data/conf目录下。