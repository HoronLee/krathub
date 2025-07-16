package client

import (
	"context"
	"fmt"
	"time"

	"krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
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
	rawLogger log.Logger         // 原始日志记录器（用于中间件）
}

// NewGrpcClientFactory 创建一个新的 GRPC 客户端工厂
func NewGrpcClientFactory(config *conf.Data, discovery registry.Discovery, logger log.Logger) (GrpcClientFactory, error) {
	return &grpcClientFactory{
		config:    config,
		discovery: discovery,
		logger:    log.NewHelper(logger),
		rawLogger: logger,
	}, nil
}

// CreateGrpcConn 创建gRPC连接
func (f *grpcClientFactory) CreateGrpcConn(ctx context.Context, serviceName string) (ggrpc.ClientConnInterface, error) {
	// 默认超时时间
	timeout := 5 * time.Second

	// 默认使用服务发现
	endpoint := fmt.Sprintf("discovery:///%s", serviceName)
	enableTLS := false

	// 尝试获取服务特定配置（如果存在）
	for _, c := range f.config.Client.GetGrpc() {
		if c.ServiceName == serviceName {
			// 使用服务特定的超时设置（如果有）
			if c.Timeout != nil {
				timeout = c.Timeout.AsDuration()
			}

			// 检查是否配置了特定的endpoint
			if c.Endpoint != "" {
				// 使用配置的endpoint替代服务发现
				endpoint = c.Endpoint
				f.logger.Infof("using configured endpoint for service %s: %s", serviceName, endpoint)

				// 检查是否需要启用TLS
				enableTLS = c.EnableTls
			}
			break
		}
	}

	// 创建gRPC连接
	var conn *ggrpc.ClientConn
	var err error

	// 准备中间件
	middleware := []middleware.Middleware{
		recovery.Recovery(),
		tracing.Client(),
		// 使用Logger()方法获取底层的Logger接口实现
		logging.Client(f.rawLogger),
	}

	if enableTLS {
		// 使用TLS连接
		conn, err = grpc.Dial(
			ctx,
			grpc.WithEndpoint(endpoint),
			grpc.WithTLSConfig(nil), // 这里可以根据需要配置TLS证书
			grpc.WithTimeout(timeout),
			grpc.WithMiddleware(middleware...),
		)
	} else if endpoint == fmt.Sprintf("discovery:///%s", serviceName) && f.discovery != nil {
		// 使用服务发现
		conn, err = grpc.DialInsecure(
			ctx,
			grpc.WithEndpoint(endpoint),
			grpc.WithDiscovery(f.discovery),
			grpc.WithTimeout(timeout),
			grpc.WithMiddleware(middleware...),
		)
	} else {
		// 使用直连但不启用TLS
		conn, err = grpc.DialInsecure(
			ctx,
			grpc.WithEndpoint(endpoint),
			grpc.WithTimeout(timeout),
			grpc.WithMiddleware(middleware...),
		)
	}

	if err != nil {
		f.logger.Errorf("failed to create grpc client for service %s: %v", serviceName, err)
		return nil, fmt.Errorf("failed to create grpc client for service %s: %w", serviceName, err)
	}

	f.logger.Debugf("successfully created grpc client for service: %s with endpoint: %s", serviceName, endpoint)
	return conn, nil
}
