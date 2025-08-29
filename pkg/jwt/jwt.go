package jwt

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWT is a generic struct for handling JWT operations.
// The type parameter 'T' should be your custom claims struct (e.g., MyUserClaims).
// A pointer to your struct (*T) must implement the jwt.Claims interface.
// The easiest way to do this is to embed jwt.RegisteredClaims in your struct.
type JWT[T any] struct {
	secretKey []byte
}

// Config holds the configuration for the JWT service.
type Config struct {
	SecretKey string
}

// NewJWT creates a new generic JWT service.
func NewJWT[T any](cfg *Config) *JWT[T] {
	return &JWT[T]{
		secretKey: []byte(cfg.SecretKey),
	}
}

// GenerateToken creates a new JWT token with the provided claims.
// The claims argument must be a pointer to your custom claims struct.
func (j *JWT[T]) GenerateToken(claims *T) (string, error) {
	// A pointer to the claims struct must implement jwt.Claims.
	// We perform a runtime check here.
	jwtClaims, ok := any(claims).(jwt.Claims)
	if !ok {
		return "", fmt.Errorf("claims type *%T does not implement jwt.Claims. Did you forget to embed jwt.RegisteredClaims?", *claims)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(j.secretKey)
}

// AnalyseToken parses a token string and returns the populated custom claims.
func (j *JWT[T]) AnalyseToken(tokenString string) (*T, error) {
	// Create a new pointer to a zero-value claims object of type T.
	claims := new(T)

	// A pointer to the claims struct must implement jwt.Claims.
	// We check this with a type assertion before passing it to the parser.
	// This is required due to a limitation in Go's generics where the compiler
	// cannot prove that '*T' implements an interface even if it will at runtime.
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

	// Since claimsInterface was a wrapper around the `claims` pointer,
	// `claims` now contains the parsed and validated data.
	return claims, nil
}

// authKey is an unexported type used as a key for storing claims in context
// to prevent collisions with other packages.
type authKey struct{}

// NewContext stores user claims into the context.
func NewContext[T any](ctx context.Context, claims *T) context.Context {
	return context.WithValue(ctx, authKey{}, claims)
}

// FromContext retrieves user claims from the context.
func FromContext[T any](ctx context.Context) (*T, bool) {
	claims, ok := ctx.Value(authKey{}).(*T)
	return claims, ok
}
