package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWT 是一个用于处理 JWT 操作的通用结构体。
// 类型参数 'T' 应该是你的自定义 claims 结构体（例如，MyUserClaims）。
// 你的结构体指针 (*T) 必须实现 jwt.Claims 接口。
// 最简单的方法是在你的结构体中嵌入 jwt.RegisteredClaims。
type JWT[T any] struct {
	secretKey []byte
}

// Config 保存 JWT 服务的配置。
type Config struct {
	SecretKey string
}

// NewJWT 创建一个新的通用 JWT 服务。
func NewJWT[T any](cfg *Config) *JWT[T] {
	return &JWT[T]{
		secretKey: []byte(cfg.SecretKey),
	}
}

// GenerateToken 使用提供的 claims 创建一个新的 JWT 令牌。
// claims 参数必须是你的自定义 claims 结构体的指针。
func (j *JWT[T]) GenerateToken(claims *T) (string, error) {
	// claims 结构体的指针必须实现 jwt.Claims。
	// 我们在这里执行运行时检查。
	jwtClaims, ok := any(claims).(jwt.Claims)
	if !ok {
		return "", fmt.Errorf("claims type *%T does not implement jwt.Claims. Did you forget to embed jwt.RegisteredClaims?", *claims)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(j.secretKey)
}

// ParseToken 解析令牌字符串并返回填充的自定义 claims。
func (j *JWT[T]) ParseToken(tokenString string) (*T, error) {
	// 创建一个指向 T 类型零值 claims 对象的新指针。
	claims := new(T)

	// claims 结构体的指针必须实现 jwt.Claims。
	// 在将其传递给解析器之前，我们通过类型断言来检查这一点。
	// 这是由于 Go 泛型的一个限制，即编译器
	// 无法证明 '*T' 实现了一个接口，即使它在运行时会实现。
	claimsInterface, ok := any(claims).(jwt.Claims)
	if !ok {
		return nil, fmt.Errorf("claims type *%T does not implement jwt.Claims. Did you forget to embed jwt.RegisteredClaims?", *claims)
	}

	token, err := jwt.ParseWithClaims(tokenString, claimsInterface, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 由于 claimsInterface 是 `claims` 指针的包装器，
	// `claims` 现在包含了已解析和验证的数据。
	return claims, nil
}

// authKey 是一个未导出的类型，用作在上下文中存储 claims 的键
// 以防止与其他包发生冲突。
type authKey struct{}

// NewContext 将用户 claims 存储到上下文中。
func NewContext[T any](ctx context.Context, claims *T) context.Context {
	return context.WithValue(ctx, authKey{}, claims)
}

// FromContext 从上下文中检索用户 claims。
func FromContext[T any](ctx context.Context) (*T, bool) {
	claims, ok := ctx.Value(authKey{}).(*T)
	return claims, ok
}
