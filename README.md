# Krathub

> 基于Kratos框架编写的快开框架，目前处于开发初期阶段

开发顺序: api -> config -> service -> biz -> data

## 项目依赖
项目中需要使用wire依赖注入,安装wire:

```bash
go install github.com/google/wire/cmd/wire@latest
```

protobuf 参数校验插件

```bash
go install github.com/envoyproxy/protoc-gen-validate@latest
```

## Data层编码须知

编写 data 层代码之前需要先修改configs目录下的config.yaml文件来配置数据库等相关信息。然后再通过`make gendb`来生成 orm 代码