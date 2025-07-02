package client

import (
	"context"
	"fmt"
	"time"

	"krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	ggrpc "google.golang.org/grpc"
)

// GrpcClientFactory 定义客户端工厂接口
type GrpcClientFactory interface {

	// CreateGrpcConn 创建gRPC连接
	CreateGrpcConn(ctx context.Context, serviceName string) (ggrpc.ClientConnInterface, error)
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

// CreateGrpcConn 创建gRPC连接
func (f *grpcClientFactory) CreateGrpcConn(ctx context.Context, serviceName string) (ggrpc.ClientConnInterface, error) {
	// 默认超时时间
	timeout := 5 * time.Second

	// 尝试获取服务特定配置（如果存在）
	for _, c := range f.config.Client.GetGrpc() {
		if c.ServiceName == serviceName {
			// 使用服务特定的超时设置（如果有）
			if c.Timeout != nil {
				timeout = c.Timeout.AsDuration()
			}
			break
		}
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
