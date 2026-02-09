# AGENTS.md - API 定义层

<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-02-09 | Updated: 2026-02-09 -->

## 目的

`api/` 目录是 Krathub 项目的 API 定义中心，负责管理所有 Protobuf 定义文件和生成的代码。该目录使用 Buf 工具链进行现代化的 Protobuf 管理，支持双协议（gRPC + HTTP）和自动 OpenAPI 文档生成。

**核心职责**：
- 定义所有微服务的 API 契约（gRPC 和 HTTP）
- 管理 Protobuf 依赖（通过 Buf Schema Registry）
- 生成 Go protobuf 代码、gRPC 服务代码、HTTP 路由代码
- 生成 Kratos 错误代码和参数验证代码
- 生成 OpenAPI 文档（每个服务独立配置）

**设计原则**：
- **契约优先**：所有 API 变更从 proto 定义开始
- **代码生成**：通过 `make gen` 统一生成所有代码，不手动编辑生成代码
- **版本化管理**：使用 `v1`, `v2` 等版本号管理 API 演进
- **依赖隔离**：通过 Buf 远程依赖管理外部 proto（如 googleapis）

## 关键文件

### Buf 配置文件

- **buf.gen.yaml** - Buf 代码生成主配置
  - 配置 Go protobuf 代码生成插件（protoc-gen-go, protoc-gen-go-grpc）
  - 配置 Kratos 插件（protoc-gen-go-http, protoc-gen-go-errors）
  - 配置参数验证插件（protoc-gen-validate）
  - 管理 `go_package` 路径映射（使用 managed mode）
  - 支持 Python 代码生成（已注释）

- **buf.work.yaml** - Buf workspace 配置
  - 定义 workspace 目录列表
  - 当前只包含 `protos/` 目录

- **buf.krathub.openapi.gen.yaml** - Krathub 服务 OpenAPI 生成配置
  - 生成 `app/krathub/service/openapi.yaml`
  - 仅包含 `krathub/service/v1/i_*.proto`（HTTP 接口）

- **buf.sayhello.openapi.gen.yaml** - SayHello 服务 OpenAPI 生成配置
  - 生成 `app/sayhello/service/openapi.yaml`
  - 包含 `sayhello/service/v1/sayhello.proto`

### 重要说明

- **生成的代码不提交**：`gen/` 目录已加入 `.gitignore`
- **所有修改从 proto 开始**：修改 proto 后必须运行 `make gen`
- **包名映射规则**：在 `buf.gen.yaml` 中统一管理，避免包名冲突

## 子目录结构

### `protos/` - Protobuf 源文件

Proto 定义文件的根目录，按服务和功能组织。

**关键文件**：
- **buf.yaml** - Proto 模块配置
  - 定义 Buf 远程依赖（BSR）：
    - `buf.build/googleapis/googleapis` - Google APIs（HTTP 注解等）
    - `buf.build/kratos/apis` - Kratos 框架 APIs（错误定义等）
    - `buf.build/bufbuild/protovalidate` - 参数验证
    - `buf.build/gnostic/gnostic` - OpenAPI 生成支持
  - 配置 breaking change 检测和 lint 规则

**目录组织**：

```
protos/
├── buf.yaml                            # Buf 模块配置
├── buf.lock                            # 依赖锁定文件
│
├── conf/v1/                            # 配置定义
│   ├── conf.proto                      # 配置结构定义（Bootstrap, Server, Data 等）
│   ├── config-example.yaml             # 配置示例文件
│   └── config-*.yaml                   # 各环境配置文件
│
├── krathub/service/v1/                 # Krathub HTTP 接口
│   ├── i_auth.proto                    # 认证 HTTP API（登录、注册、登出）
│   ├── i_user.proto                    # 用户 HTTP API（用户 CRUD）
│   ├── i_test.proto                    # 测试 HTTP API
│   └── krathub_doc.proto               # OpenAPI 文档元数据
│
├── auth/service/v1/                    # Auth gRPC 服务
│   └── auth.proto                      # 认证 gRPC 接口（VerifyToken 等）
│
├── user/service/v1/                    # User gRPC 服务
│   └── user.proto                      # 用户 gRPC 接口（GetUser, UpdateUser 等）
│
├── test/service/v1/                    # Test gRPC 服务
│   └── test.proto                      # 测试 gRPC 接口
│
└── sayhello/service/v1/                # SayHello 独立微服务
    ├── sayhello.proto                  # SayHello gRPC/HTTP 接口
    └── sayhello_doc.proto              # OpenAPI 文档元数据
```

**Proto 文件组织规范**：

1. **HTTP 接口**（统一包名）
   - 路径：`krathub/service/v1/i_*.proto`
   - 包名：`krathub.service.v1`
   - 命名：以 `i_` 前缀标识 HTTP 接口
   - 用途：生成 Kratos HTTP 路由（使用 `google.api.http` 注解）

2. **gRPC 服务**（独立包名）
   - 路径：`{service}/service/v1/{service}.proto`
   - 包名：`{service}.service.v1`
   - 用途：生成 gRPC 服务代码

3. **独立微服务**（完整定义）
   - 路径：`sayhello/service/v1/sayhello.proto`
   - 包名：`sayhello.service.v1`
   - 用途：独立微服务，可同时包含 HTTP 和 gRPC 定义

4. **配置定义**（特殊处理）
   - 路径：`conf/v1/conf.proto`
   - 包名：`conf.v1`
   - 用途：定义服务配置结构（YAML 配置文件映射）
   - 注意：`buf.gen.yaml` 中保留原始 `go_package` 设置

### `gen/go/` - 生成的 Go 代码

自动生成的 Go protobuf 代码，由 `make gen` 生成。

**目录结构**：
```
gen/go/
├── auth/service/v1/           # Auth 服务生成代码
│   ├── auth.pb.go             # Protobuf 消息定义
│   ├── auth_grpc.pb.go        # gRPC 服务接口
│   └── auth_errors.pb.go      # Kratos 错误定义
│
├── user/service/v1/           # User 服务生成代码
│   ├── user.pb.go
│   ├── user_grpc.pb.go
│   └── user_errors.pb.go
│
├── krathub/service/v1/        # Krathub HTTP 接口生成代码
│   ├── i_auth.pb.go
│   ├── i_auth_http.pb.go      # HTTP 路由和客户端
│   ├── i_user.pb.go
│   ├── i_user_http.pb.go
│   └── ...
│
└── conf/v1/                   # 配置定义生成代码
    └── conf.pb.go
```

**生成的文件类型**：
- `*.pb.go` - Protobuf 消息定义（protoc-gen-go）
- `*_grpc.pb.go` - gRPC 服务接口（protoc-gen-go-grpc）
- `*_http.pb.go` - Kratos HTTP 路由（protoc-gen-go-http）
- `*_errors.pb.go` - Kratos 错误定义（protoc-gen-go-errors）
- `*.pb.validate.go` - 参数验证代码（protoc-gen-validate）

**重要警告**：
- **不要手动编辑**：所有代码通过 `make gen` 自动生成
- **不要提交到 Git**：已在 `.gitignore` 中排除
- **每次修改 proto 后重新生成**：否则会导致代码不同步

## AI Agent 工作指南

### 添加新 API

#### 1. 添加 HTTP 接口（推荐）

**场景**：为 Krathub 主服务添加新的 HTTP 接口

**步骤**：

```bash
# 1. 在 protos/krathub/service/v1/ 创建新文件
# 文件名格式：i_{feature}.proto
```

```protobuf
// 示例：protos/krathub/service/v1/i_article.proto
syntax = "proto3";

package krathub.service.v1;

import "google/api/annotations.proto";
import "validate/validate.proto";

option go_package = "github.com/horonlee/krathub/api/gen/go/krathub/service/v1;krathubpb";

// 文章服务 HTTP 接口
service Article {
  // 创建文章
  rpc CreateArticle (CreateArticleRequest) returns (CreateArticleReply) {
    option (google.api.http) = {
      post: "/api/v1/articles"
      body: "*"
    };
  }

  // 获取文章列表
  rpc ListArticles (ListArticlesRequest) returns (ListArticlesReply) {
    option (google.api.http) = {
      get: "/api/v1/articles"
    };
  }
}

message CreateArticleRequest {
  string title = 1 [(validate.rules).string = {min_len: 1, max_len: 200}];
  string content = 2 [(validate.rules).string.min_len = 1];
}

message CreateArticleReply {
  int64 id = 1;
}

message ListArticlesRequest {
  int32 page = 1 [(validate.rules).int32 = {gte: 1}];
  int32 page_size = 2 [(validate.rules).int32 = {gte: 1, lte: 100}];
}

message ListArticlesReply {
  repeated Article articles = 1;
  int64 total = 2;
}

message Article {
  int64 id = 1;
  string title = 2;
  string content = 3;
}
```

```bash
# 2. 生成代码
cd /Users/horonlee/projects/micro-service/krathub
make gen

# 3. 验证生成结果
ls api/gen/go/krathub/service/v1/i_article*
# 应该看到：
# - i_article.pb.go
# - i_article_http.pb.go
# - i_article.pb.validate.go

# 4. 在服务中实现接口
# 编辑 app/krathub/service/internal/service/article.go
```

#### 2. 添加 gRPC 服务

**场景**：添加新的独立 gRPC 服务

**步骤**：

```bash
# 1. 创建服务目录
mkdir -p protos/article/service/v1
```

```protobuf
// 示例：protos/article/service/v1/article.proto
syntax = "proto3";

package article.service.v1;

import "validate/validate.proto";

option go_package = "github.com/horonlee/krathub/api/gen/go/article/service/v1;articlepb";

// 文章服务 gRPC 接口
service ArticleService {
  rpc GetArticle (GetArticleRequest) returns (GetArticleReply);
  rpc UpdateArticle (UpdateArticleRequest) returns (UpdateArticleReply);
}

message GetArticleRequest {
  int64 id = 1 [(validate.rules).int64.gt = 0];
}

message GetArticleReply {
  Article article = 1;
}

message UpdateArticleRequest {
  int64 id = 1 [(validate.rules).int64.gt = 0];
  string title = 2;
  string content = 3;
}

message UpdateArticleReply {
  Article article = 1;
}

message Article {
  int64 id = 1;
  string title = 2;
  string content = 3;
  int64 created_at = 4;
  int64 updated_at = 5;
}
```

```bash
# 2. 在 buf.gen.yaml 中添加 go_package 映射
# 编辑 api/buf.gen.yaml，在 override 部分添加：
```

```yaml
- file_option: go_package
  path: article/service/v1
  value: github.com/horonlee/krathub/api/gen/go/article/service/v1;articlepb
```

```bash
# 3. 生成代码
make gen

# 4. 验证生成结果
ls api/gen/go/article/service/v1/
# 应该看到：
# - article.pb.go
# - article_grpc.pb.go
# - article.pb.validate.go
```

#### 3. 添加错误定义

**在现有 proto 文件中添加错误定义**：

```protobuf
// 在 protos/krathub/service/v1/i_article.proto 中添加

import "errors/errors.proto";

enum ErrorReason {
  option (errors.default_code) = 500;

  ARTICLE_NOT_FOUND = 0 [(errors.code) = 404];
  ARTICLE_ALREADY_EXISTS = 1 [(errors.code) = 409];
  ARTICLE_INVALID_TITLE = 2 [(errors.code) = 400];
}
```

生成后可在代码中使用：
```go
import krathubpb "github.com/horonlee/krathub/api/gen/go/krathub/service/v1"

return krathubpb.ErrorArticleNotFound("article id=%d not found", id)
```

### 修改现有 API

```bash
# 1. 编辑对应的 .proto 文件
# 例如：api/protos/krathub/service/v1/i_user.proto

# 2. 重新生成代码
make gen

# 3. 更新服务实现代码
# 例如：app/krathub/service/internal/service/user.go

# 4. 运行测试验证
make test
```

**注意事项**：
- **向后兼容**：避免删除字段或修改字段编号
- **添加新字段**：使用新的字段编号，添加默认值
- **废弃字段**：使用 `reserved` 关键字标记
- **Breaking Change 检测**：Buf 会自动检测破坏性变更

### 生成 OpenAPI 文档

#### 为新服务添加 OpenAPI 配置

```bash
# 1. 复制现有配置
cp api/buf.krathub.openapi.gen.yaml api/buf.newservice.openapi.gen.yaml
```

```yaml
# 2. 编辑配置文件
version: v2

plugins:
  - remote: buf.build/community/pseudomuto-protoc-gen-doc
    out: app/newservice/service
    opt:
      - openapi.yaml,json
    strategy: all

inputs:
  - directory: protos
    include:
      - newservice/service/v1/newservice.proto  # 指定要生成文档的 proto 文件
```

```bash
# 3. 生成文档
buf generate --template api/buf.newservice.openapi.gen.yaml api/protos

# 4. 验证生成结果
cat app/newservice/service/openapi.yaml
```

### 管理 Protobuf 依赖

#### 添加新的 Buf 依赖

```bash
# 编辑 api/protos/buf.yaml
```

```yaml
deps:
  - "buf.build/googleapis/googleapis"
  - "buf.build/kratos/apis"
  - "buf.build/bufbuild/protovalidate"
  - "buf.build/gnostic/gnostic"
  - "buf.build/your-org/your-module"  # 添加新依赖
```

```bash
# 更新依赖锁定文件
cd api/protos
buf mod update

# 验证依赖
buf dep
```

### 代码生成工作流

**完整的代码生成流程**：

```bash
# 在项目根目录执行
make gen
```

该命令会依次执行：
1. `buf generate --template api/buf.gen.yaml api/protos` - 生成 Go 代码
2. `buf generate --template api/buf.krathub.openapi.gen.yaml api/protos` - 生成 Krathub OpenAPI
3. `buf generate --template api/buf.sayhello.openapi.gen.yaml api/protos` - 生成 SayHello OpenAPI

**仅生成特定服务的代码**：

```bash
# 仅生成 Go 代码
buf generate --template api/buf.gen.yaml api/protos

# 仅生成 Krathub OpenAPI
buf generate --template api/buf.krathub.openapi.gen.yaml api/protos
```

### 常见问题排查

#### 问题 1：生成代码失败

**错误信息**：`protoc-gen-go: program not found or is not executable`

**解决方案**：
```bash
# 安装 protoc 插件
make plugin

# 验证安装
which protoc-gen-go
which protoc-gen-go-grpc
which protoc-gen-go-http
```

#### 问题 2：导入路径错误

**错误信息**：`import "google/api/annotations.proto" not found`

**解决方案**：
```bash
# 更新 Buf 依赖
cd api/protos
buf mod update

# 清理并重新生成
cd ../..
make gen
```

#### 问题 3：go_package 路径冲突

**错误信息**：`duplicate go_package path`

**解决方案**：
检查 `api/buf.gen.yaml` 中的 `override` 配置，确保每个服务有唯一的包名：

```yaml
override:
  - file_option: go_package
    path: newservice/service/v1
    value: github.com/horonlee/krathub/api/gen/go/newservice/service/v1;newservicepb
```

#### 问题 4：Breaking Change 警告

**错误信息**：`Field "1" on message "User" changed type`

**解决方案**：
```bash
# 查看详细的 breaking change 报告
buf breaking --against '.git#branch=main'

# 如果是有意的破坏性变更，更新文档并通知使用方
# 如果是错误，回退修改
```

### 测试 API 定义

#### Lint 检查

```bash
cd api/protos
buf lint

# 查看所有 lint 规则
buf config ls-lint-rules
```

#### Breaking Change 检测

```bash
# 检测相对于 main 分支的破坏性变更
buf breaking --against '.git#branch=main'

# 检测相对于特定 tag 的破坏性变更
buf breaking --against '.git#tag=v1.0.0'
```

#### 格式化 Proto 文件

```bash
cd api/protos
buf format -w  # -w 表示写入文件
```

### 集成开发工作流

**完整的 API 开发流程**：

```bash
# 1. 创建功能分支
git checkout -b feature/add-article-api

# 2. 定义 API（编辑 proto 文件）
vim api/protos/krathub/service/v1/i_article.proto

# 3. Lint 检查
cd api/protos && buf lint && cd ../..

# 4. 生成代码
make gen

# 5. 实现服务
# 编辑 app/krathub/service/internal/service/article.go
# 编辑 app/krathub/service/internal/biz/article.go
# 编辑 app/krathub/service/internal/data/article.go

# 6. 更新依赖注入
cd app/krathub/service
make wire

# 7. 运行测试
make test

# 8. 本地运行验证
make run

# 9. 提交代码
git add .
git commit -m "feat(api): add article API"
git push origin feature/add-article-api
```

## 依赖关系

### 构建依赖

- **Buf CLI** - Protobuf 管理工具
  - 安装：`go install github.com/bufbuild/buf/cmd/buf@latest`
  - 版本要求：v1.0.0+

- **protoc 插件**（通过 `make plugin` 安装）
  - `protoc-gen-go` - Go protobuf 代码生成
  - `protoc-gen-go-grpc` - gRPC 服务代码生成
  - `protoc-gen-go-http` - Kratos HTTP 路由生成
  - `protoc-gen-go-errors` - Kratos 错误代码生成
  - `protoc-gen-validate` - 参数验证代码生成
  - `protoc-gen-openapi` - OpenAPI 文档生成

### 远程依赖（Buf Schema Registry）

在 `api/protos/buf.yaml` 中定义：

- **buf.build/googleapis/googleapis** - Google APIs
  - 提供 `google/api/annotations.proto`（HTTP 注解）
  - 提供 `google/api/http.proto`（HTTP 规则）

- **buf.build/kratos/apis** - Kratos 框架 APIs
  - 提供 `errors/errors.proto`（错误定义）
  - 提供 Kratos 特有的 proto 扩展

- **buf.build/bufbuild/protovalidate** - 参数验证
  - 提供 `validate/validate.proto`（验证规则）

- **buf.build/gnostic/gnostic** - OpenAPI 支持
  - 用于 OpenAPI 文档生成

### 被依赖关系

生成的代码被以下模块使用：

- **app/krathub/service/** - Krathub 主服务
  - 导入 `krathub/service/v1`（HTTP 接口）
  - 导入 `auth/service/v1`（Auth gRPC）
  - 导入 `user/service/v1`（User gRPC）
  - 导入 `test/service/v1`（Test gRPC）

- **app/sayhello/service/** - SayHello 服务
  - 导入 `sayhello/service/v1`

- **pkg/** - 共享库
  - 可能导入 `conf/v1`（配置定义）

**导入示例**：
```go
import (
    krathubpb "github.com/horonlee/krathub/api/gen/go/krathub/service/v1"
    authpb "github.com/horonlee/krathub/api/gen/go/auth/service/v1"
    userpb "github.com/horonlee/krathub/api/gen/go/user/service/v1"
)
```

## 注意事项

### 文件生成规则

- **生成的代码不提交**：`gen/` 目录已在 `.gitignore` 中排除
- **每次修改 proto 后重新生成**：运行 `make gen`
- **不要手动编辑生成代码**：会在下次生成时被覆盖
- **Wire 集成**：生成代码后需运行 `make wire` 更新依赖注入

### Proto 编写规范

- **使用 `validate` 约束**：为所有输入参数添加验证规则
- **添加注释**：proto 注释会生成到 OpenAPI 文档
- **字段编号管理**：已使用的编号不要重复使用，废弃字段用 `reserved`
- **包名约定**：使用 `{service}.service.v1` 格式
- **HTTP 路径规范**：RESTful 风格，使用 `/api/v1/` 前缀

### Breaking Change 管理

- **向后兼容优先**：尽量保持向后兼容
- **版本化策略**：重大变更通过版本号管理（v1 → v2）
- **CI 集成**：在 CI 中运行 `buf breaking` 检测
- **文档更新**：破坏性变更需更新 CHANGELOG 和迁移指南

### OpenAPI 文档

- **自动生成**：通过 `buf.{service}.openapi.gen.yaml` 配置
- **文档位置**：生成到 `app/{service}/service/openapi.yaml`
- **访问方式**：服务启动后通过 `/q/swagger-ui` 访问
- **元数据管理**：通过 `{service}_doc.proto` 定义文档元数据

## 快速参考

### 常用命令

```bash
# 生成所有代码
make gen

# Lint 检查
cd api/protos && buf lint

# Breaking change 检测
cd api/protos && buf breaking --against '.git#branch=main'

# 格式化 proto 文件
cd api/protos && buf format -w

# 更新依赖
cd api/protos && buf mod update

# 查看依赖
cd api/protos && buf dep
```

### Proto 文件模板

```protobuf
syntax = "proto3";

package {service}.service.v1;

import "google/api/annotations.proto";
import "validate/validate.proto";
import "errors/errors.proto";

option go_package = "github.com/horonlee/krathub/api/gen/go/{service}/service/v1;{service}pb";

// 服务定义
service {Service} {
  rpc {Method} ({Request}) returns ({Reply}) {
    option (google.api.http) = {
      post: "/api/v1/{path}"
      body: "*"
    };
  }
}

// 请求消息
message {Request} {
  string field = 1 [(validate.rules).string.min_len = 1];
}

// 响应消息
message {Reply} {
  int64 id = 1;
}

// 错误定义
enum ErrorReason {
  option (errors.default_code) = 500;

  {ERROR_NAME} = 0 [(errors.code) = 404];
}
```

### 目录创建检查清单

添加新服务时需要：
- [ ] 创建 `protos/{service}/service/v1/{service}.proto`
- [ ] 在 `buf.gen.yaml` 添加 `go_package` 映射
- [ ] 创建 `buf.{service}.openapi.gen.yaml`（如需 OpenAPI）
- [ ] 运行 `make gen` 生成代码
- [ ] 在服务中实现接口
- [ ] 更新根目录 Makefile（如果是新服务）
