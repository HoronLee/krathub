package server

import (
	"crypto/tls"
	"github.com/horonlee/krathub/internal/conf"
	mw "github.com/horonlee/krathub/internal/server/middleware"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	gogrpc "google.golang.org/grpc" // 引入官方 gRPC 包并重命名
	"google.golang.org/grpc/credentials"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, trace *conf.Trace, mM *mw.MiddlewareManager, logger log.Logger) *grpc.Server {
	var mds []middleware.Middleware
	mds = []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(logger),
	}
	// 开启链路追踪
	if trace != nil && trace.Endpoint != "" {
		mds = append(mds, tracing.Server())
	}

	var opts = []grpc.ServerOption{
		grpc.Middleware(mds...),
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

	// Add TLS configuration
	if c.Grpc.Tls != nil && c.Grpc.Tls.Enable {
		cert, err := tls.LoadX509KeyPair(c.Grpc.Tls.CertPath, c.Grpc.Tls.KeyPath)
		if err != nil {
			logger.Log(log.LevelFatal, "msg", "gRPC Server TLS: Failed to load key pair", "error", err)
		}
		creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		opts = append(opts, grpc.Options(gogrpc.Creds(creds)))
	}

	srv := grpc.NewServer(opts...)
	return srv
}
