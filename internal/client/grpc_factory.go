package client

import (
	"context"
	"fmt"
	"time"

	hellov1 "krathub/api/hello/v1"
	"krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
)

// GrpcClientFactory 定义客户端工厂接口
type GrpcClientFactory interface {
	// 创建Hello服务客户端
	NewHelloClient(ctx context.Context) (hellov1.HelloServiceClient, error)
}

// grpcClientFactory 是 GrpcClientFactory接口的实现
type grpcClientFactory struct {
	config    *conf.Data         // 配置信息
	discovery registry.Discovery // 服务发现客户端
	logger    *log.Helper        // 日志助手
}

// NewGrpcClientFactory 创建一个新的 GRPC 客户端工厂
func NewGrpcClientFactory(config *conf.Data, discovery registry.Discovery, logger log.Logger) (GrpcClientFactory, error) {
	return &grpcClientFactory{
		config:    config,
		discovery: discovery,
		logger:    log.NewHelper(logger),
	}, nil
}

// createGrpcConn 创建gRPC连接
func (f *grpcClientFactory) createGrpcConn(ctx context.Context, serviceName string) (ggrpc.ClientConnInterface, error) {
	// 获取服务配置
	var serviceConfig *conf.Data_Client_GRPC
	for _, c := range f.config.Client.GetGrpc() {
		if c.ServiceName == serviceName {
			serviceConfig = c
			break
		}
	}

	if serviceConfig == nil {
		return nil, fmt.Errorf("no grpc client config found for service: %s", serviceName)
	}

	// 设置默认超时时间
	timeout := 5 * time.Second
	if serviceConfig.Timeout != nil {
		timeout = serviceConfig.Timeout.AsDuration()
	}

	// 创建gRPC连接
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", serviceName)),
		grpc.WithDiscovery(f.discovery),
		grpc.WithTimeout(timeout),
		grpc.WithMiddleware(
			recovery.Recovery(),
			// 可以在这里添加更多中间件
		),
	)
	if err != nil {
		f.logger.Errorf("failed to create grpc client for service %s: %v", serviceName, err)
		return nil, fmt.Errorf("failed to create grpc client for service %s: %w", serviceName, err)
	}

	f.logger.Debugf("successfully created grpc client for service: %s", serviceName)
	return conn, nil
}
