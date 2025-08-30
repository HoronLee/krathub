package server

import (
	"crypto/tls"

	authV1 "github.com/horonlee/krathub/api/auth/v1"
	userV1 "github.com/horonlee/krathub/api/user/v1"
	"github.com/horonlee/krathub/internal/conf"
	"github.com/horonlee/krathub/internal/consts"
	mw "github.com/horonlee/krathub/internal/server/middleware"
	"github.com/horonlee/krathub/internal/service"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, trace *conf.Trace, auth *service.AuthService, user *service.UserService, mM *mw.MiddlewareManager, m *Metrics, logger log.Logger) *http.Server {
	var mds []middleware.Middleware
	mds = []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(logger),
		ratelimit.Server(),
		validate.ProtoValidate(),
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
	}
	// 开启链路追踪
	if trace != nil && trace.Endpoint != "" {
		mds = append(mds, tracing.Server())
	}

	// 开启 metrics
	if m != nil {
		mds = append(mds, metrics.Server(
			metrics.WithSeconds(m.Seconds),
			metrics.WithRequests(m.Requests),
		))
	}

	var opts = []http.ServerOption{
		http.Middleware(mds...),
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

	// Add TLS configuration
	if c.Http.Tls != nil && c.Http.Tls.Enable {
		if c.Http.Tls.CertPath == "" || c.Http.Tls.KeyPath == "" {
			logger.Log(log.LevelFatal, "msg", "Server TLS: can't find TLS key pairs")
		}
		cert, err := tls.LoadX509KeyPair(c.Http.Tls.CertPath, c.Http.Tls.KeyPath)
		if err != nil {
			logger.Log(log.LevelFatal, "msg", "Server TLS: Failed to load key pair", "error", err)
		}
		opts = append(opts, http.TLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}}))
	}

	srv := http.NewServer(opts...)
	// 开启 metrics 路由
	if m != nil {
		srv.Handle("/metrics", m.Handler)
	}
	// 注册服务
	authV1.RegisterAuthHTTPServer(srv, auth)
	userV1.RegisterUserHTTPServer(srv, user)
	return srv
}
