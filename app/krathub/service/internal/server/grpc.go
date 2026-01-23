package server

import (
	"crypto/tls"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	pkglogger "github.com/horonlee/krathub/pkg/logger"

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

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, trace *conf.Trace, m *Metrics, logger log.Logger) *grpc.Server {
	grpcLogger := pkglogger.WithModule(logger, "grpc/server/krathub-service")
	var mws []middleware.Middleware
	mws = []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(grpcLogger),
		ratelimit.Server(),
		validate.ProtoValidate(),
	}
	// 开启链路追踪
	if trace != nil && trace.Endpoint != "" {
		mws = append(mws, tracing.Server())
	}

	// 开启 metrics
	if m != nil {
		mws = append(mws, metrics.Server(
			metrics.WithSeconds(m.Seconds),
			metrics.WithRequests(m.Requests),
		))
	}
	opts := []grpc.ServerOption{
		grpc.Middleware(mws...),
		grpc.Logger(grpcLogger),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	if c.Grpc.Tls != nil && c.Grpc.Tls.Enable {
		cert, err := tls.LoadX509KeyPair(c.Grpc.Tls.CertPath, c.Grpc.Tls.KeyPath)
		if err != nil {
			grpcLogger.Log(log.LevelFatal, "msg", "gRPC Server TLS: Failed to load key pair", "error", err)
		}
		creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		opts = append(opts, grpc.Options(gogrpc.Creds(creds)))
	}
	srv := grpc.NewServer(opts...)

	// 注册服务

	return srv
}
