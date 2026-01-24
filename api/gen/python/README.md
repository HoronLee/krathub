# Python gRPC Client

这个目录包含从 Protobuf 定义自动生成的 Python gRPC 客户端代码。

## 安装依赖

```bash
pip install grpcio grpcio-tools
```

或使用 uv（推荐）：

```bash
uv pip install grpcio grpcio-tools
```

## 使用示例

### 1. 导入生成的代码

```python
import sys
sys.path.append('/path/to/krathub/api/gen/python')

from auth.service.v1 import auth_pb2
from auth.service.v1 import auth_pb2_grpc
import grpc
```

### 2. 创建 gRPC 客户端

```python
def create_auth_client():
    channel = grpc.insecure_channel('localhost:9001')
    stub = auth_pb2_grpc.AuthServiceStub(channel)
    return stub

def login(username, password):
    stub = create_auth_client()
    request = auth_pb2.LoginRequest(
        username=username,
        password=password
    )
    
    try:
        response = stub.Login(request)
        print(f"Login successful! Token: {response.token}")
        return response
    except grpc.RpcError as e:
        print(f"Login failed: {e.details()}")
        return None
```

### 3. 完整示例

```python
#!/usr/bin/env python3
import sys
import grpc

sys.path.append('/path/to/krathub/api/gen/python')

from user.service.v1 import user_pb2
from user.service.v1 import user_pb2_grpc

def main():
    with grpc.insecure_channel('localhost:9001') as channel:
        stub = user_pb2_grpc.UserServiceStub(channel)
        
        request = user_pb2.GetUserRequest(id=1)
        
        try:
            response = stub.GetUser(request)
            print(f"User: {response.user.username}")
            print(f"Email: {response.user.email}")
        except grpc.RpcError as e:
            print(f"Error: {e.code()} - {e.details()}")

if __name__ == '__main__':
    main()
```

## 目录结构

```
python/
├── auth/service/v1/          # 认证服务
│   ├── auth_pb2.py          # 消息定义
│   └── auth_pb2_grpc.py     # gRPC 服务定义
├── user/service/v1/          # 用户服务
│   ├── user_pb2.py
│   └── user_pb2_grpc.py
├── test/service/v1/          # 测试服务
│   ├── test_pb2.py
│   └── test_pb2_grpc.py
└── krathub/service/v1/       # HTTP 接口定义
    ├── i_auth_pb2.py
    ├── i_user_pb2.py
    └── ...
```

## 注意事项

1. **自动生成** - 这些文件是自动生成的，不要手动修改
2. **重新生成** - 运行 `make api` 会重新生成所有代码
3. **导入路径** - 确保将 `api/gen/python` 添加到 Python 的 `sys.path`
4. **gRPC 端口** - 默认 gRPC 服务运行在 `localhost:9001`

## 与 Go 服务通信

Python 客户端可以直接与 Krathub 的 Go gRPC 服务通信：

```python
# 连接到 Krathub gRPC 服务
channel = grpc.insecure_channel('localhost:9001')

# 使用任何服务的 stub
auth_stub = auth_pb2_grpc.AuthServiceStub(channel)
user_stub = user_pb2_grpc.UserServiceStub(channel)
test_stub = test_pb2_grpc.TestServiceStub(channel)
```

## 开发建议

### 使用虚拟环境

```bash
# 使用 uv 创建虚拟环境
cd /path/to/your/python/project
uv venv
source .venv/bin/activate

# 安装依赖
uv pip install grpcio grpcio-tools
```

### 创建 requirements.txt

```txt
grpcio>=1.60.0
grpcio-tools>=1.60.0
```

### 项目结构建议

```
your-python-project/
├── .venv/                    # 虚拟环境
├── requirements.txt          # Python 依赖
├── client.py                 # 你的客户端代码
└── krathub_api/             # 符号链接到生成的代码
    -> /path/to/krathub/api/gen/python
```

## 更多信息

- [gRPC Python 官方文档](https://grpc.io/docs/languages/python/)
- [Protobuf Python 教程](https://protobuf.dev/getting-started/pythontutorial/)
