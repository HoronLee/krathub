package client

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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

// GrpcConn gRPC连接实现
type GrpcConn struct {
	conn gogrpc.ClientConnInterface
	ref  int32 // 引用计数
	mu   sync.RWMutex
}

// NewGrpcConn 创建gRPC连接封装
func NewGrpcConn(conn gogrpc.ClientConnInterface) *GrpcConn {
	return &GrpcConn{
		conn: conn,
	}
}

// Value 返回原始gRPC连接
func (g *GrpcConn) Value() interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.conn
}

// Close 减少引用计数（参考pool示例）
func (g *GrpcConn) Close() error {
	newRef := atomic.AddInt32(&g.ref, -1)
	if newRef < 0 {
		panic(fmt.Sprintf("negative ref: %d", newRef))
	}
	return nil
}

// IsHealthy 检查连接健康状态
func (g *GrpcConn) IsHealthy() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.conn != nil
}

// GetType 返回连接类型
func (g *GrpcConn) GetType() ConnType {
	return GRPC
}

// reset 实际关闭连接（由工厂管理）
func (g *GrpcConn) reset() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.conn != nil {
		if closer, ok := g.conn.(interface{ Close() error }); ok {
			err := closer.Close()
			g.conn = nil
			return err
		}
	}
	return nil
}

// createGrpcConnection 创建gRPC连接的内部函数
func createGrpcConnection(ctx context.Context, serviceName string, dataCfg *conf.Data,
	traceCfg *conf.Trace, discovery registry.Discovery, logger log.Logger) (gogrpc.ClientConnInterface, error) {

	// 默认超时时间
	timeout := 5 * time.Second
	endpoint := fmt.Sprintf("discovery:///%s", serviceName)

	// 尝试获取服务特定配置
	for _, c := range dataCfg.Client.GetGrpc() {
		if c.ServiceName == serviceName {
			if c.Timeout != nil {
				timeout = c.Timeout.AsDuration()
			}
			if c.Endpoint != "" {
				endpoint = c.Endpoint
				logger.Log(log.LevelInfo, "msg", "using configured endpoint",
					"service_name", serviceName, "endpoint", endpoint)
			}
			break
		}
	}

	// 准备中间件
	middleware := []middleware.Middleware{
		recovery.Recovery(),
		logging.Client(logger),
		circuitbreaker.Client(),
	}

	if traceCfg != nil && traceCfg.Endpoint != "" {
		middleware = append(middleware, tracing.Client())
	}

	// 创建gRPC连接
	var conn *gogrpc.ClientConn
	var err error

	if endpoint == fmt.Sprintf("discovery:///%s", serviceName) && discovery != nil {
		conn, err = grpc.DialInsecure(
			ctx,
			grpc.WithEndpoint(endpoint),
			grpc.WithDiscovery(discovery),
			grpc.WithTimeout(timeout),
			grpc.WithMiddleware(middleware...),
		)
	} else {
		conn, err = grpc.DialInsecure(
			ctx,
			grpc.WithEndpoint(endpoint),
			grpc.WithTimeout(timeout),
			grpc.WithMiddleware(middleware...),
		)
	}

	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to create grpc client",
			"service_name", serviceName, "error", err)
		return nil, fmt.Errorf("failed to create grpc client for service %s: %w", serviceName, err)
	}

	logger.Log(log.LevelDebug, "msg", "successfully created grpc client",
		"service_name", serviceName, "endpoint", endpoint)

	return conn, nil
}
