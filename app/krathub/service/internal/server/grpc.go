package server

import (
	"crypto/tls"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/horonlee/krathub/pkg/governance/telemetry"
	logpkg "github.com/horonlee/krathub/pkg/logger"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GRPCMiddleware 用于 Wire 注入的中间件切片包装类型
type GRPCMiddleware []middleware.Middleware

// NewGRPCMiddleware 创建 gRPC 中间件
func NewGRPCMiddleware(
	trace *conf.Trace,
	m *telemetry.Metrics,
	logger log.Logger,
) GRPCMiddleware {
	grpcLogger := logpkg.With(logger, logpkg.WithModule("grpc/server/krathub-service"))

	var ms []middleware.Middleware
	ms = append(ms,
		recovery.Recovery(),
		logging.Server(grpcLogger),
		ratelimit.Server(),
		validate.ProtoValidate(),
	)

	// 开启链路追踪
	if trace != nil && trace.Endpoint != "" {
		ms = append(ms, tracing.Server())
	}

	// 开启 metrics
	if m != nil {
		ms = append(ms, metrics.Server(
			metrics.WithSeconds(m.Seconds),
			metrics.WithRequests(m.Requests),
		))
	}

	return ms
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	middlewares GRPCMiddleware,
	logger log.Logger,
) *grpc.Server {
	grpcLogger := logpkg.With(logger, logpkg.WithModule("grpc/server/krathub-service"))
	helper := log.NewHelper(grpcLogger)

	opts := []grpc.ServerOption{
		grpc.Middleware(middlewares...),
		grpc.Logger(grpcLogger),
	}
	if c != nil && c.Grpc != nil {
		if c.Grpc.Network != "" {
			opts = append(opts, grpc.Network(c.Grpc.Network))
		}
		if c.Grpc.Addr != "" {
			opts = append(opts, grpc.Address(c.Grpc.Addr))
		}
		if c.Grpc.Timeout != nil {
			opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
		}
	}
	if c != nil && c.Grpc != nil && c.Grpc.Tls != nil && c.Grpc.Tls.Enable {
		cert, err := tls.LoadX509KeyPair(c.Grpc.Tls.CertPath, c.Grpc.Tls.KeyPath)
		if err != nil {
			helper.Fatalf("gRPC Server TLS: Failed to load key pair: %v", err)
		}
		creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		opts = append(opts, grpc.Options(gogrpc.Creds(creds)))
	}

	srv := grpc.NewServer(opts...)

	// 注册服务

	return srv
}
