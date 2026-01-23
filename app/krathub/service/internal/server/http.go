package server

import (
	"crypto/tls"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	krathubv1 "github.com/horonlee/krathub/api/gen/go/krathub/service/v1"
	"github.com/horonlee/krathub/app/krathub/service/internal/consts"
	mw "github.com/horonlee/krathub/app/krathub/service/internal/server/middleware"
	"github.com/horonlee/krathub/app/krathub/service/internal/service"
	pkglogger "github.com/horonlee/krathub/pkg/logger"
	"github.com/horonlee/krathub/pkg/middleware/cors"

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
func NewHTTPServer(c *conf.Server, trace *conf.Trace, authJWT mw.AuthJWT, m *Metrics, logger log.Logger, auth *service.AuthService, user *service.UserService, test *service.TestService) *http.Server {
	httpLogger := pkglogger.WithModule(logger, "http/server/krathub-service")
	var mws []middleware.Middleware
	mws = []middleware.Middleware{
		recovery.Recovery(),
		logging.Server(httpLogger),
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
	// 配置特殊路由
	mws = append(mws, configureRoutes(authJWT)...)
	var opts = []http.ServerOption{
		http.Middleware(mws...),
		http.Logger(httpLogger),
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
	if c.Http.Cors != nil {
		corsOptions := mw.CORS(c.Http.Cors)
		if len(corsOptions.AllowedOrigins) > 0 {
			opts = append(opts, http.Filter(cors.Middleware(corsOptions)))
			httpLogger.Log(log.LevelInfo, "msg", "CORS middleware enabled", "allowed_origins", corsOptions.AllowedOrigins)
		}
	}
	if c.Http.Tls != nil && c.Http.Tls.Enable {
		if c.Http.Tls.CertPath == "" || c.Http.Tls.KeyPath == "" {
			httpLogger.Log(log.LevelFatal, "msg", "Server TLS: can't find TLS key pairs")
		}
		cert, err := tls.LoadX509KeyPair(c.Http.Tls.CertPath, c.Http.Tls.KeyPath)
		if err != nil {
			httpLogger.Log(log.LevelFatal, "msg", "Server TLS: Failed to load key pair", "error", err)
		}
		opts = append(opts, http.TLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}}))
	}
	srv := http.NewServer(opts...)
	// 开启 metrics 路由
	if m != nil {
		srv.Handle("/metrics", m.Handler)
	}

	// 注册服务 - 使用 krathub HTTP 服务（i_*.proto 生成的）
	krathubv1.RegisterAuthServiceHTTPServer(srv, auth)
	krathubv1.RegisterUserServiceHTTPServer(srv, user)
	krathubv1.RegisterTestServiceHTTPServer(srv, test)
	return srv
}

// configureRoutes 配置权限路由中间件
func configureRoutes(authMiddleware mw.AuthJWT) []middleware.Middleware {
	return []middleware.Middleware{
		// 需要User权限的接口
		selector.Server(authMiddleware(consts.UserRole(2))).
			Prefix("/krathub.service.v1.UserService/", "/krathub.service.v1.TestService/").
			Build(),
		// 需要Admin权限的接口
		selector.Server(authMiddleware(consts.UserRole(3))).
			Path("/krathub.service.v1.UserService/DeleteUser", "/krathub.service.v1.UserService/SaveUser").
			Build(),
	}
}
