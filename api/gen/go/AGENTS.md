# AGENTS.md - 生成的 Go 代码

<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-02-09 | Updated: 2026-02-09 -->

## 重要说明

⚠️ **这是自动生成的代码，请勿手动编辑。**

本目录下的所有文件均由 `buf` 根据 `api/protos/` 中的 Protobuf 定义自动生成。任何手动修改都会在下次代码生成时被覆盖。

## 如何重新生成代码

若需更新此目录下的代码，请在项目根目录下执行：

```bash
make gen
```

该命令会调用 `buf generate` 并根据 `api/buf.gen.yaml` 的配置重新生成所有 Go 代码。

## 子目录结构

- `auth/`: 认证服务 (gRPC/Errors)
- `user/`: 用户服务 (gRPC/Errors)
- `krathub/`: Krathub 网关服务 (HTTP/Validate)
- `sayhello/`: 示例服务
- `test/`: 测试服务
- `conf/`: 静态配置定义
