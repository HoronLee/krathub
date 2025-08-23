# Krathub

[English](README.en-US.md) | 简体中文

> 基于Kratos框架编写的快开框架，目前处于开发初期阶段

## 如何使用

使用kratos layout功能快速通过krathub模板创建本地项目

### Multirepo单仓模式

```bash
kratos new PeojectName -r https://github.com/HoronLee/krathub.git
```

### Monorepo大仓模式

使用 --nomod 添加服务，共用 go.mod ，大仓模式，其中app/user可以自定义为想要的微服务结构

```shell
kratos new helloworld
cd helloworld
kratos new app/user --nomod -r https://github.com/HoronLee/krathub.git
```

使用大仓模式会启用远程proto仓库，如单仓下的import为`authv1 "krathub/api/auth/v1"`，但是使用大仓模式初始化的项目会变为具体的url`authv1  "github.com/ExampleUser/projectName/api/auth/v1"`
也就是说大仓模式下需要将proto文件全部放在仓库的根目录的api文件下！

## 开发须知

开发顺序: api -> config -> service -> biz -> data -> client

功能编写完成后需要使用`make wire`来进行依赖注入，并且需要在`internal/server`中的`NewServer`方法中添加用法，注意需要手动在方法签名中添加依赖

## 项目依赖

直接执行`make init`即可下载所需软件

## Data层编码须知

### 数据库

编写 data 层代码之前需要先修改configs目录下的config.yaml文件来配置数据库等相关信息。然后再通过`make gendb`来生成 orm 代码

## Client层

client层是本人自己新增的客户端层，级别上来说比 data 层低一层，目前包含了grpc客户端的工厂方法。这个层的是用于调用外部grpc服务而设计的，后续可能会添加http客户端的功能，但是考虑到微服务环境下大多还是以grpc为沟通协议，所以暂不实现。


## Docker Compose 部署

可用的docker compose文件在项目的deployment/docker-compose目录下，首次运行请把model.sql放入initdb文件夹中，这样数据库首次运行就会导入数据。配置文件放于data/conf目录下。

## 配置文件示例

`./configs/config.yaml`
```yaml
server:
  http:
    addr: "${HADDR:0.0.0.0:8000}"
    timeout: "${HTIMEOUT:1s}"
    tls:
      enable: false
      cert_path: "${HTTPS_CERT_PATH:../../manifest/certs/server.cert}"
      key_path: "${HTTPS_KEY_PATH:../../manifest/certs/server.key}"
  grpc:
    addr: "${GADDR:0.0.0.0:8001}"
    timeout: "${GTIMEHOUT:1s}"
    tls:
      enable: false
      cert_path: "${GRPCS_CERT_PATH:../../manifest/certs/server.cert}"
      key_path: "${GRPCS_KEY_PATH:../../manifest/certs/server.key}"

data:
  database:
    driver: "${DB_DRIVER:mysql}"
    source: "${DB_SOURCE:projectName:123456@tcp(127.0.0.1:3306)/projectName?parseTime=True&loc=Local}"
  redis:
    addr: "${RADDR:127.0.0.1:6379}"
    user_name: "${RUSER_NAME:}"  # Redis用户名
    password: "${RPASSWORD:redisHoron}"  # Redis密码
    db: "${RDB:5}"  # Redis数据库编号
    read_timeout: "${RREAD_TIMEOUT:0.2s}"
    write_timeout: "${RWRITE_TIMEOUT:0.2s}"
  client:
    # 这里可以配置第三方服务的客户端
    # 默认不用配置，而是在代码中直接服务发现
    grpc:
      # - service_name: hello.grpc  # nacos需要添加协议后缀
      #   endpoint: "${GRPC_ENDPOINT:127.0.0.1:8003}"
      #   enable_tls: false
      #   timeout: 5s

app:
  env: "${ENV:dev}" # dev, test, prod
  jwt:
    secret_key: "${JWT_SECRETK_KEY:projectName_secret}"
    expire: "${JWT_EXPIRE:24}"
    issuer: "${JWT_ISSUER:projectName}"
    # audience: "${JWT_AUDIENCE:projectName}"
  log:
    level: "${LOG_LEVEL:-1}"  # -1debug,0info,1warn,2error,3dpanic,4panic,5fatal
    filename: "${LOG_FILENAME:projectName.log}"  # 日志文件夹为根目录logs
    max_size: "${LOG_MAX_SIZE:20}"  # 日志文件最大大小，单位MB
    max_age: "${LOG_MAX_AGE:30}"  # 日志文件最大保存天数
    max_backups: "${LOG_MAX_BACKUPS:10}"  # 日志文件最大备份数

# 注册中心配置
registry:
  # 使用 Consul 作为注册中心
  consul:
    addr: consul.r430.com:30080
    scheme: http
    datacenter: dc1
    timeout: 5s

  # 或者使用 Nacos 作为注册中心
  # nacos:
    # addr: "${NACOSR_ADDR:127.0.0.1}"
    # port: "${NACOSR_PORT:8848}"
    # namespace: "${NACOSR_NAMESPACE:public}"
    # group: "${NACOSR_GROUP:DEFAULT_GROUP}"
    # username: "${NACOSR_USERNAME:nacos}"
    # password: "${NACOSR_PASSWORD:nacos}"
    # timeout: "${NACOSR_TIMEOUT:5s}"

# 服务发现配置，一般和注册中心配置相同
discovery:
  consul:
    addr: consul.r430.com:30080
    scheme: http
    datacenter: dc1
    timeout: 5s

  # nacos:
  #   addr: "${NACOSD_ADDR:127.0.0.1}"
  #   port: "${NACOSD_PORT:8848}"
  #   namespace: "${NACOSD_NAMESPACE:public}"
  #   group: "${NACOSD_GROUP:DEFAULT_GROUP}"
  #   username: "${NACOSD_USERNAME:nacos}"
  #   password: "${NACOSD_PASSWORD:nacos}"
  #   timeout: "${NACOSD_TIMEOUT:5s}"

# 配置中心配置
config:

  # 使用 Nacos 作为配置中心
  # nacos:
  #   addr: "${NACOSC_ADDR:127.0.0.1}"
  #   port: "${NACOSC_PORT:8848}"
  #   namespace: "${NACOSC_NAMESPACE:public}"
  #   group: "${NACOSC_GROUP:DEFAULT_GROUP}"
  #   data_id: "${NACOSC_DATA_ID:projectName.yaml}"
  #   username: "${NACOSC_USERNAME:nacos}"
  #   password: "${NACOSC_PASSWORD:nacos}"
  #   timeout: "${NACOSC_TIMEOUT:5s}"

  # 使用 consul 作为配置中心
  # consul:
  #   addr: consul.r430.com:30080
  #   scheme: http
  #   datacenter: dc1
  #   timeout: 5s
  #   key: projectName/config.yaml  # 配置键名

# 链路追踪配置
trace:
  # 使用 Jaeger 作为链路追踪
  endpoint: otlp.jaeger.r430.com:30080
```
