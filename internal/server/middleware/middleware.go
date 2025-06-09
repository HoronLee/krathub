package middleware

import (
	"context"
	authV1 "krathub/api/v1/auth"
	userV1 "krathub/api/v1/user"
	"krathub/internal/conf"
	"krathub/pkg/jwt"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
)

var bc *conf.Bootstrap

func SetBootstrap(bootstrap *conf.Bootstrap) {
	bc = bootstrap
}

// Auth is a middleware for authentication service.
func Auth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				token := tr.RequestHeader().Get("Authorization")
				if token == "" {
					return nil, authV1.ErrorMissingToken("missing Authorization header")
				}
				token = strings.TrimPrefix(token, "Bearer ")
				userClaims, err := jwt.NewJWT(bc.App.Jwt).AnalyseToken(token)
				if err != nil {
					return nil, err
				} else if userClaims.Role == "" {
					return nil, authV1.ErrorUnauthorized("don't have permission to access this resource")
				}
				// 调用独立的特殊接口权限检查
				if err := checkSpecialPermission(tr.Operation(), userClaims.Role); err != nil {
					return nil, err
				}
				ctx = metadata.NewServerContext(ctx, metadata.New(map[string][]string{
					"username": {userClaims.Name},
					"role":     {userClaims.Role},
				}))
			}
			return handler(ctx, req)
		}
	}
}

// AuthWhiteListMatcher returns a selector.MatchFunc for auth service whitelist.
func AuthWhiteListMatcher() selector.MatchFunc {
	whiteList := map[string]struct{}{
		"/auth.v1.Auth/SignupByEmail":        {},
		"/auth.v1.Auth/LoginByEmailPassword": {},
	}
	return func(_ context.Context, operation string) bool {
		_, ok := whiteList[operation]
		return !ok
	}
}

// 检查特殊接口权限
func checkSpecialPermission(operation, role string) error {
	if operation == "/user.v1.User/DeleteUser" && role != "admin" {
		return userV1.ErrorDeleteUserFailed("only admin can delete user")
	}
	return nil
}
