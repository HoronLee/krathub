package client

import (
	"context"
	"fmt"
	"time"

	"github.com/horonlee/krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gogrpc "google.golang.org/grpc"
)

// grpcClientFactory 是 GrpcClientFactory接口的实现
type grpcClientFactory struct {
	dataCfg   *conf.Data         // 数据配置信息
	traceCfg  *conf.Trace        // 链路追踪配置
	discovery registry.Discovery // 服务发现客户端
	logger    log.Logger
}

// NewGrpcClientFactory 创建一个新的 GRPC 客户端工厂
func NewGrpcClientFactory(dataCfg *conf.Data, traceCfg *conf.Trace, discovery registry.Discovery, logger log.Logger) (ClientFactory, error) {
	return &grpcClientFactory{
		dataCfg:   dataCfg,
		traceCfg:  traceCfg,
		discovery: discovery,
		logger:    logger,
	}, nil
}

// CreateGrpcConn 创建gRPC连接
func (f *grpcClientFactory) CreateGrpcConn(ctx context.Context, serviceName string) (gogrpc.ClientConnInterface, error) {
	// 默认超时时间
	timeout := 5 * time.Second

	// 默认使用服务发现
	endpoint := fmt.Sprintf("discovery:///%s", serviceName)

	// 尝试获取服务特定配置（如果存在）
	for _, c := range f.dataCfg.Client.GetGrpc() {
		if c.ServiceName == serviceName {
			// 使用服务特定的超时设置（如果有）
			if c.Timeout != nil {
				timeout = c.Timeout.AsDuration()
			}

			// 检查是否配置了特定的endpoint
			if c.Endpoint != "" {
				// 使用配置的endpoint替代服务发现
				endpoint = c.Endpoint
				f.logger.Log(log.LevelInfo, "msg", "using configured endpoint", "service_name", serviceName, "endpoint", endpoint)

			}
			break
		}
	}

	// 创建gRPC连接
	var conn *gogrpc.ClientConn
	var err error

	// 准备中间件
	middleware := []middleware.Middleware{
		recovery.Recovery(),
		logging.Client(f.logger),
		circuitbreaker.Client(),
	}

	// 如果开启了链路追踪，则添加客户端追踪中间件
	if f.traceCfg != nil && f.traceCfg.Endpoint != "" {
		middleware = append(middleware, tracing.Client())
	}

	if endpoint == fmt.Sprintf("discovery:///%s", serviceName) && f.discovery != nil {
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
		f.logger.Log(log.LevelError, "msg", "failed to create grpc client", "service_name", serviceName, "error", err)
		return nil, fmt.Errorf("failed to create grpc client for service %s: %w", serviceName, err)
	}

	f.logger.Log(log.LevelDebug, "msg", "successfully created grpc client", "service_name", serviceName, "endpoint", endpoint)
	return conn, nil
}
