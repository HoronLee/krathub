# Krathub

> 基于Kratos框架编写的快开框架，目前处于开发初期阶段

开发顺序: api -> config -> service -> biz -> data -> client

功能编写完成后需要使用`make wire`来进行依赖注入，并且需要在`internal/server`中的`NewServer`方法中添加用法，注意需要手动在方法签名中添加依赖

## Client层

client层是本人自己新增的客户端层，目前包含了grpc客户端的工厂方法。这个层的是用于调用外部grpc服务而设计的，后续可能会添加http客户端的功能，但是考虑到微服务环境下大多还是以grpc为沟通协议，所以暂不实现。

Client层编码顺序：先在`GrpcClientFactory`添加所需实现的结构，后续在`grpc_client.go`中进行实现，这样可以保持代码的整洁。

随后可以在data层进行调用，首先创建一个客户端`helloClient, err := r.data.clientFactory.NewHelloClient(ctx)`，然后进行该服务的一些方法的调用`helloClient.SayHello(ctx, &hellov1.HelloRequest{Greeting: &in})`

## 项目依赖

直接执行`make init`即可下载所需软件

## Data层编码须知

编写 data 层代码之前需要先修改configs目录下的config.yaml文件来配置数据库等相关信息。然后再通过`make gendb`来生成 orm 代码

## Docker Compose 部署

可用的docker compose文件在项目的deployment/docker-compose目录下，首次运行请把model.sql放入initdb文件夹中，这样数据库首次运行就会导入数据。配置文件放于data/conf目录下。