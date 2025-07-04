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
    "github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, auth *service.AuthService, user *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate.ProtoValidate(),
			// 登录等无需鉴权接口
			selector.Server(middleware.Auth(consts.UserRole(0))).
				Prefix("/krathub.auth.v1.Auth/").
				Build(),
			// 需要User权限的接口
			selector.Server(middleware.Auth(consts.UserRole(2))).
				Prefix("/krathub.user.v1.User/").
				Build(),
			// 需要Admin权限的接口
			selector.Server(middleware.Auth(consts.UserRole(3))).
				Path("/krathub.user.v1.User/DeleteUser", "/krathub.user.v1.User/SaveUser").
				Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	authV1.RegisterAuthHTTPServer(srv, auth)
	userV1.RegisterUserHTTPServer(srv, user)
	return srv
}
