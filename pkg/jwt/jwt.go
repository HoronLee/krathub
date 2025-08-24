package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	secretKey      []byte
	expirationTime time.Duration
	audience       string
	issuer         string
}

type Config struct {
	SecretKey string
	Expire    int32
	Audience  string
	Issuer    string
}

func NewJWT(cfg *Config) *JWT {
	return &JWT{
		secretKey:      []byte(cfg.SecretKey),
		expirationTime: time.Duration(cfg.Expire) * time.Hour,
		audience:       cfg.Audience,
		issuer:         cfg.Issuer,
	}
}

type UserClaims struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func (j *JWT) GenerateToken(id int64, name, identity string) (string, error) {
	claims := &UserClaims{
		ID:   id,
		Name: name,
		Role: identity,
		StandardClaims: jwt.StandardClaims{
			Audience:  j.audience,
			ExpiresAt: time.Now().Add(j.expirationTime).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWT) AnalyseToken(tokenString string) (*UserClaims, error) {
	claims := new(UserClaims)
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

type authKey struct{}

// NewContext 将用户claims信息存入context
func NewContext(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, authKey{}, claims)
}

// FromContext 从context中提取用户claims信息
func FromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(authKey{}).(*UserClaims)
	return claims, ok
}
