package jwt

import (
	"crypto/md5"
	"fmt"
	"time"

	"krathub/internal/conf"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	secretKey      []byte
	expirationTime time.Duration
	audience       string
	issuer         string
}

func NewJWT(cfg *conf.Jwt) *JWT {
	return &JWT{
		secretKey:      []byte(cfg.SecretKey),
		expirationTime: time.Duration(cfg.Expire) * time.Hour,
		audience:       cfg.Audience,
		issuer:         cfg.Issuer,
	}
}

type UserClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func (j *JWT) GenerateToken(name, identity string) (string, error) {
	claims := &UserClaims{
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
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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
