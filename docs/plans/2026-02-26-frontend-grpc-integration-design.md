# Frontend gRPC Client Integration Design

## 1. Goal
Implement a modern, type-safe gRPC client for the Vue 3 frontend that communicates with the Go Kratos backend, replacing the current custom HTTP generator.

## 2. Comparison Summary

| Tool | Protocol Support | Proxy Req | Maintenance | Recommendation |
| :--- | :--- | :--- | :--- | :--- |
| **grpc-web** | gRPC-Web | Envoy/Middleware | Legacy/Maintenance | Not recommended for new TS projects |
| **Connect-ES** | Connect, gRPC-Web, gRPC | Middleware (Optional) | Active (Buf) | **Highly Recommended** |
| **ts-proto** | Custom/gRPC-Web | Depends | Active (Community) | Good for plain objects |

## 3. Recommended Approach: Connect-ES + Kratos gRPC-Web
We will use the **Connect-ES** ecosystem from Buf. It provides the best developer experience in 2026 for TypeScript/Vite projects.

### Frontend (Vue 3 / Vite)
- **Generator**: `@connectrpc/protoc-gen-connect-es` and `@bufbuild/protoc-gen-protobuf-es`.
- **Runtime**: `@connectrpc/connect` and `@connectrpc/connect-web`.
- **Transport**: `createGrpcWebTransport` to ensure compatibility with Kratos gRPC backend.
- **Integration**: Update `buf.gen.yaml` to automate code generation.

### Backend (Go Kratos)
- **Middleware**: Use `github.com/improbable-eng/grpc-web/go/grpcweb` to wrap the Kratos gRPC server.
- **Configuration**: Expose necessary CORS headers (e.g., `grpc-status`, `grpc-message`) to allow browser to read gRPC errors.

## 4. Implementation Steps
1. **API Update**: Update `buf.gen.yaml` to include Connect-ES plugins.
2. **Frontend Setup**: Install `@connectrpc/connect`, `@connectrpc/connect-web`, `@bufbuild/protobuf`.
3. **Backend Middleware**: Add gRPC-Web wrapper to Kratos `internal/server/grpc.go`.
4. **Client Instance**: Create a shared transport instance in `web/src/api/client.ts`.

## 5. Browser Limitations & Mitigations
- **No HTTP/2 Trailers**: Mitigated by using the gRPC-Web protocol (standardized by Connect-ES).
- **CORS**: Mitigated by explicit header exposure in Kratos/gRPC-Web middleware.
