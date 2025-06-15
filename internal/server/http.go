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
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, auth *service.AuthService, user *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate.Validator(),
			selector.Server(middleware.Auth(consts.UserRole(3))).
				Path("/krathub.user.v1.User/DeleteUser").
				Build(),
			selector.Server(middleware.Auth(consts.UserRole(2))).
				Match(middleware.AuthWhiteListMatcher()).
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
