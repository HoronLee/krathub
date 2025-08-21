package server

import (
	authV1 "krathub/api/auth/v1"
	"krathub/internal/conf"
	"krathub/internal/service"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, trace *conf.Trace, logger log.Logger, auth *service.AuthService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			validate.ProtoValidate(),
		),
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

	// 开启链路追踪
	if trace != nil && trace.Endpoint != "" {
		opts = append(opts, grpc.Middleware(tracing.Server()))
	}

	srv := grpc.NewServer(opts...)
	authV1.RegisterAuthServer(srv, auth)
	return srv
}
