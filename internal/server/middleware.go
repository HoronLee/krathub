package server

import (
	"context"
	authV1 "krathub/api/auth/v1"
	"krathub/internal/conf"
	"krathub/internal/consts"
	"krathub/pkg/jwt"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// MiddlewareManager 中间件管理器
type MiddlewareManager struct {
	appConf *conf.App
}

// NewMiddlewareManager 创建中间件管理器
func NewMiddlewareManager(appConf *conf.App) *MiddlewareManager {
	return &MiddlewareManager{
		appConf: appConf,
	}
}

// Auth 认证中间件
// 当 minRole 为 0 时，允许无 Authorization 头访问（如注册接口）
func (m *MiddlewareManager) Auth(minRole consts.UserRole) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, authV1.ErrorMissingToken("missing transport context")
			}
			authHeader := tr.RequestHeader().Get("Authorization")
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// 如果未设置 minRole（即为 0），允许无 token 访问
			if minRole == 0 && tokenString == "" {
				return handler(ctx, req)
			}

			if tokenString == "" {
				return nil, authV1.ErrorMissingToken("missing Authorization header")
			}

			// 创建JWT实例并解析Token
			jwtInstance := jwt.NewJWT(m.appConf.Jwt)
			claims, err := jwtInstance.AnalyseToken(tokenString)
			if err != nil {
				return nil, authV1.ErrorUnauthorized("invalid token: " + err.Error())
			}

			// 验证用户角色
			var userRole consts.UserRole
			switch claims.Role {
			case "guest":
				userRole = consts.Guest
			case "user":
				userRole = consts.User
			case "admin":
				userRole = consts.Admin
			case "operator":
				userRole = consts.Operator
			default:
				return nil, authV1.ErrorUnauthorized("unknown role")
			}

			if userRole < minRole {
				return nil, authV1.ErrorUnauthorized("permission denied, you at least need " + minRole.String() + " role")
			}

			// 将用户claims存入context
			ctx = jwt.NewContext(ctx, claims)

			return handler(ctx, req)
		}
	}
}
