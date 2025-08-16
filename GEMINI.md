# Krathub

> 基于Kratos框架编写的快开框架，目前处于开发初期阶段

## 开发须知

开发顺序: api -> config -> service -> biz -> data -> client

功能编写完成后需要使用`make wire`来进行依赖注入，并且需要在`internal/server`中的`NewServer`方法中添加用法，注意需要手动在方法签名中添加依赖

## 项目依赖

直接执行`make init`即可下载所需软件

## Data层编码须知

### 数据库

编写 data 层代码之前需要先修改configs目录下的config.yaml文件来配置数据库等相关信息。然后再通过`make gendb`来生成 orm 代码

### GRPC调用

与原版 kratos 不同的是，现在每个 biz 层的服务都有两种 repo，一种是传统的数据库 dbRepo，用于和各种数据库进行交互；还有一种是 grpcRepo，用于调用 client 层实现的 grpc 客户端来进行远程调用操作

## Client层

client层是本人自己新增的客户端层，级别上来说比 data 层低一层，目前包含了grpc客户端的工厂方法。这个层的是用于调用外部grpc服务而设计的，后续可能会添加http客户端的功能，但是考虑到微服务环境下大多还是以grpc为沟通协议，所以暂不实现。
