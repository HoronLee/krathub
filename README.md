# Krathub

> 基于Kratos框架编写的快开框架

开发顺序: api -> config -> service -> biz -> data -> service

项目中推荐使用wire依赖注入,安装wire:

```bash
go install github.com/google/wire/cmd/wire@latest
```

```

使用之前需要先修改configs目录下的config.yaml文件来配置数据库等相关信息

数据库可使用Gen来进行操作
