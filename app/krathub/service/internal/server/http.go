package server

import (
	"crypto/tls"

	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	krathubv1 "github.com/horonlee/krathub/api/gen/go/krathub/service/v1"
	"github.com/horonlee/krathub/app/krathub/service/internal/consts"
	mwinter "github.com/horonlee/krathub/app/krathub/service/internal/server/middleware"
	"github.com/horonlee/krathub/app/krathub/service/internal/service"
	logpkg "github.com/horonlee/krathub/pkg/logger"
	mwpkg "github.com/horonlee/krathub/pkg/middleware"
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

// HTTPMiddleware 用于 Wire 注入的中间件切片包装类型
type HTTPMiddleware []middleware.Middleware

// NewHTTPMiddleware 创建 HTTP 中间件（使用白名单机制）
func NewHTTPMiddleware(
	trace *conf.Trace,
	m *Metrics,
	logger log.Logger,
	authJWT mwinter.AuthJWT,
) HTTPMiddleware {
	httpLogger := logpkg.WithModule(logger, "http/server/krathub-service")

	var ms []middleware.Middleware
	ms = append(ms,
		recovery.Recovery(),
		logging.Server(httpLogger),
		ratelimit.Server(),
		validate.ProtoValidate(),
	)

	if trace != nil && trace.Endpoint != "" {
		ms = append(ms, tracing.Server())
	}

	if m != nil {
		ms = append(ms, metrics.Server(
			metrics.WithSeconds(m.Seconds),
			metrics.WithRequests(m.Requests),
		))
	}

	// 公开接口白名单（无需认证）
	publicWhitelist := mwpkg.NewWhiteList(mwpkg.Exact,
		krathubv1.OperationAuthServiceLoginByEmailPassword,
		krathubv1.OperationAuthServiceRefreshToken,
		krathubv1.OperationAuthServiceSignupByEmail,
		krathubv1.OperationTestServiceTest,
		krathubv1.OperationTestServiceHello,
	)

	// User 级接口白名单（需要 User 权限但跳过 Admin 检查）
	userWhitelist := mwpkg.NewWhiteList(mwpkg.Exact,
		krathubv1.OperationUserServiceCurrentUserInfo,
		krathubv1.OperationUserServiceUpdateUser,
		krathubv1.OperationTestServicePrivateTest,
	)

	// Admin 权限排除白名单 = 公开接口 ∪ User 级接口
	adminExcludeWhitelist := publicWhitelist.Merge(userWhitelist)

	ms = append(ms,
		selector.Server(authJWT(consts.User)).
			Match(publicWhitelist.MatchFunc()).
			Build(),
		selector.Server(authJWT(consts.Admin)).
			Match(adminExcludeWhitelist.MatchFunc()).
			Build(),
	)

	return ms
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	middlewares HTTPMiddleware,
	m *Metrics,
	logger log.Logger,
	auth *service.AuthService,
	user *service.UserService,
	test *service.TestService,
) *http.Server {
	httpLogger := logpkg.WithModule(logger, "http/server/krathub-service")

	var opts = []http.ServerOption{
		http.Middleware(middlewares...),
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
		corsOptions := mwinter.CORS(c.Http.Cors)
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

	if m != nil {
		srv.Handle("/metrics", m.Handler)
	}

	krathubv1.RegisterAuthServiceHTTPServer(srv, auth)
	krathubv1.RegisterUserServiceHTTPServer(srv, user)
	krathubv1.RegisterTestServiceHTTPServer(srv, test)

	return srv
}
