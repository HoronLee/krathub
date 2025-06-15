package server

import (
	authV1 "krathub/api/auth/v1"
	userV1 "krathub/api/user/v1"
	"krathub/internal/conf"
	"krathub/internal/consts"
	"krathub/internal/server/middleware"
	"krathub/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, auth *service.AuthService, user *service.UserService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate.Validator(),
			selector.Server(middleware.Auth(consts.UserRole(0))).
				Prefix("/krathub.auth.v1.Auth/").
				Build(),
			selector.Server(middleware.Auth(consts.UserRole(1))).
				Prefix("/krathub.user.v1.User/").
				Build(),
			selector.Server(middleware.Auth(consts.UserRole(3))).
				Path("/krathub.user.v1.User/DeleteUser").
				Build(),
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
	srv := grpc.NewServer(opts...)
	authV1.RegisterAuthServer(srv, auth)
	userV1.RegisterUserServer(srv, user)
	return srv
}
