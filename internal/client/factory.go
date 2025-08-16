package client

import (
	"context"

	ggrpc "google.golang.org/grpc"
)

// ClientFactory 定义客户端工厂接口
type ClientFactory interface {

	// CreateGrpcConn 创建gRPC连接
	CreateGrpcConn(ctx context.Context, serviceName string) (ggrpc.ClientConnInterface, error)
}
