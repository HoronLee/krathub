package client

import (
	"context"
	hellov1 "krathub/api/hello/v1"
)

// NewHelloClient 创建Hello服务的客户端
func (f *grpcClientFactory) NewHelloClient(ctx context.Context) (hellov1.HelloServiceClient, error) {
	conn, err := f.createGrpcConn(ctx, "hello")
	if err != nil {
		return nil, err
	}
	return hellov1.NewHelloServiceClient(conn), nil
}
