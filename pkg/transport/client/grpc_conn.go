package client

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	pkglogger "github.com/horonlee/krathub/pkg/logger"

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
func (g *GrpcConn) Value() any {
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

// GetGRPCConn 创建并提取指定服务的 gRPC 连接接口。
func GetGRPCConn(ctx context.Context, c Client, serviceName string) (gogrpc.ClientConnInterface, error) {
	connWrapper, err := c.CreateConn(ctx, GRPC, serviceName)
	if err != nil {
		return nil, err
	}

	conn, ok := connWrapper.Value().(gogrpc.ClientConnInterface)
	if !ok {
		return nil, fmt.Errorf("unexpected grpc connection type: %T", connWrapper.Value())
	}

	return conn, nil
}

// createGrpcConnection 创建gRPC连接的内部函数
func createGrpcConnection(ctx context.Context, serviceName string, dataCfg *conf.Data,
	traceCfg *conf.Trace, discovery registry.Discovery, logger log.Logger) (gogrpc.ClientConnInterface, error) {
	setupLogger := log.NewHelper(pkglogger.With(logger, pkglogger.WithField("operation", "createGrpcConnection")))

	timeout := 5 * time.Second
	defaultEndpoint := fmt.Sprintf("discovery:///%s", serviceName)
	endpoint := defaultEndpoint

	if dataCfg != nil && dataCfg.Client != nil {
		for _, c := range dataCfg.Client.GetGrpc() {
			if c == nil || c.ServiceName != serviceName {
				continue
			}
			if cfgTimeout := c.GetTimeout(); cfgTimeout != nil {
				if d := cfgTimeout.AsDuration(); d > 0 {
					timeout = d
				}
			}
			if configuredEndpoint := strings.TrimSpace(c.GetEndpoint()); configuredEndpoint != "" {
				endpoint = configuredEndpoint
				setupLogger.Infof("using configured endpoint: service_name=%s endpoint=%s", serviceName, endpoint)
			}
			break
		}
	}

	mds := []middleware.Middleware{
		recovery.Recovery(),
		logging.Client(logger),
		circuitbreaker.Client(),
	}

	if traceCfg != nil && traceCfg.Endpoint != "" {
		mds = append(mds, tracing.Client())
	}

	opts := []grpc.ClientOption{
		grpc.WithEndpoint(endpoint),
		grpc.WithTimeout(timeout),
		grpc.WithMiddleware(mds...),
	}
	if endpoint == defaultEndpoint && discovery != nil {
		opts = append(opts, grpc.WithDiscovery(discovery))
	}

	conn, err := grpc.DialInsecure(ctx, opts...)

	if err != nil {
		setupLogger.Errorf("failed to create grpc client: service_name=%s error=%v", serviceName, err)
		return nil, fmt.Errorf("failed to create grpc client for service %s: %w", serviceName, err)
	}

	setupLogger.Infof("successfully created grpc client: service_name=%s endpoint=%s timeout=%s", serviceName, endpoint, timeout.String())

	return conn, nil
}
