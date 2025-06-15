package middleware

import (
	"context"
	authV1 "krathub/api/auth/v1"
	"krathub/internal/conf"
	"krathub/internal/consts"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
)

var bc *conf.Bootstrap

func SetBootstrap(bootstrap *conf.Bootstrap) {
	bc = bootstrap
}

// Auth is a middleware for authentication service.
func Auth(minRole consts.UserRole) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, authV1.ErrorMissingToken("missing transport context")
			}
			authHeader := tr.RequestHeader().Get("Authorization")
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return nil, authV1.ErrorMissingToken("missing Authorization header")
			}

			// 解析 JWT Token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				return []byte(bc.App.Jwt.SecretKey), nil
			})
			if err != nil || !token.Valid {
				return nil, authV1.ErrorUnauthorized("invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, authV1.ErrorUnauthorized("invalid claims")
			}

			roleStr, ok := claims["role"].(string)
			if !ok {
				return nil, authV1.ErrorUnauthorized("role not found in token")
			}

			var userRole consts.UserRole
			switch roleStr {
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

			// 可选：将用户信息写入 metadata
			ctx = metadata.NewServerContext(ctx, metadata.New(map[string][]string{
				"role": {roleStr},
			}))

			return handler(ctx, req)
		}
	}
}

// AuthWhiteListMatcher returns a selector.MatchFunc for auth service whitelist.
func AuthWhiteListMatcher() selector.MatchFunc {
	whiteList := map[string]struct{}{
		"/krathub.auth.v1.Auth/SignupByEmail":        {},
		"/krathub.auth.v1.Auth/LoginByEmailPassword": {},
	}
	return func(_ context.Context, operation string) bool {
		_, ok := whiteList[operation]
		return !ok
	}
}
