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
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, trace *conf.Trace, auth *service.AuthService, user *service.UserService, mM *middleware.MiddlewareManager, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			// 登录等无需鉴权接口
			selector.Server(mM.Auth(consts.UserRole(0))).
				Prefix("/krathub.auth.v1.Auth/").
				Build(),
			// 需要User权限的接口
			selector.Server(mM.Auth(consts.UserRole(2))).
				Prefix("/krathub.user.v1.User/").
				Build(),
			// 需要Admin权限的接口
			selector.Server(mM.Auth(consts.UserRole(3))).
				Path("/krathub.user.v1.User/DeleteUser", "/krathub.user.v1.User/SaveUser").
				Build(),
			mM.ProtoValidate(),
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
	// 开启链路追踪
	if trace != nil && trace.Endpoint != "" {
		opts = append(opts, http.Middleware(tracing.Server()))
	}
	// 参数校验中间件
	// 将"github.com/go-kratos/kratos/contrib/middleware/validate/v2"下载到本地作为中间件
	// TODO: 此中间件有未知问题【如果不放在最后注册，则会失效】
	opts = append(opts, http.Middleware(mM.ProtoValidate()))

	srv := http.NewServer(opts...)
	authV1.RegisterAuthHTTPServer(srv, auth)
	userV1.RegisterUserHTTPServer(srv, user)
	return srv
}
