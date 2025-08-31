package biz

import (
	"context"

	"github.com/horonlee/krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/golang-jwt/jwt/v5"
)

// UserClaims JWT 的自定义声明示例
// type UserClaims struct {
// ID   int64  `json:"id"`
// Name string `json:"name"`
// Role string `json:"role"`
// jwt.RegisteredClaims
// }

// CallHelloRepo 统一的认证仓库接口，包含数据库和 grpc 操作
type CallHelloRepo interface {
	// grpc 操作
	CallHello(ctx context.Context, in string) (string, error)

	// 数据库操作示例
	// SaveUser(context.Context, *model.User) (*model.User, error)
	// ListUserByEmail(context.Context, string) ([]*model.User, error)
	// ListUserByUserName(context.Context, string) ([]*model.User, error)
}

// CallHelloUsecase is a Auth usecase.
type CallHelloUsecase struct {
	repo CallHelloRepo
	log  *log.Helper
	cfg  *conf.App
	// jwt             *jwtpkg.JWT[UserClaims] // Use the generic JWT with UserClaims
}

// NewAuthUsecase new an auth usecase.
func NewAuthUsecase(repo CallHelloRepo, logger log.Logger, cfg *conf.App) *CallHelloUsecase {
	// JWT中间件使用示例
	// jwtService := jwtpkg.NewJWT[UserClaims](&jwtpkg.Config{
	// 	SecretKey: cfg.Jwt.SecretKey,
	// })
	uc := &CallHelloUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
		cfg:  cfg,
		// jwt:  jwtService,
	}
	// 初始化 adminRegistered
	// admins, err := repo.ListUserByUserName(context.Background(), "admin")
	// if err == nil && len(admins) > 0 {
	// 	uc.adminRegistered = true
	// }
	return uc
}

// Hello 通过 repo 实现
func (uc *CallHelloUsecase) Hello(ctx context.Context, in *string) (string, error) {
	uc.log.Debugf("Saying hello with greeting: %s", *in)
	response, err := uc.repo.CallHello(ctx, *in)
	if err != nil {
		uc.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return response, nil
}

// // SignupByEmail 使用邮件注册
// func (uc *AuthUsecase) SignupByEmail(ctx context.Context, user *model.User) (*model.User, error) {
// 	// 检查 admin 用户是否已存在
// 	if !uc.adminRegistered {
// 		// 第一次注册，用户名必须为 admin
// 		if user.Name != "admin" {
// 			return nil, authv1.ErrorInvalidCredentials("the first user must be named admin")
// 		}
// 		user.Role = "admin"
// 	} else {
// 		// 后续注册，用户名可以任意，但角色为 user
// 		// 检查用户名是否已存在
// 		existingUsers, err := uc.repo.ListUserByUserName(ctx, user.Name)
// 		if err != nil {
// 			return nil, authv1.ErrorUserNotFound("failed to check username: %v", err)
// 		}
// 		if len(existingUsers) > 0 {
// 			return nil, authv1.ErrorUserAlreadyExists("username already exists")
// 		}
// 		user.Role = "user"
// 	}

// 	// 检查邮箱是否已存在
// 	existingEmails, err := uc.repo.ListUserByEmail(ctx, user.Email)
// 	if err != nil {
// 		return nil, authv1.ErrorUserNotFound("failed to check email: %v", err)
// 	}
// 	if len(existingEmails) > 0 {
// 		return nil, authv1.ErrorUserAlreadyExists("email already exists")
// 	}

// 	createdUser, err := uc.repo.SaveUser(ctx, user)
// 	if err == nil && !uc.adminRegistered && user.Name == "admin" {
// 		uc.adminRegistered = true // 注册成功后更新状态
// 	}
// 	return createdUser, err
// }

// // generateToken 签发 JWT token
// func (uc *AuthUsecase) generateToken(claims *UserClaims) (string, error) {
// 	return uc.jwt.GenerateToken(claims)
// }

// // LoginByEmailPassword 邮箱密码登录
// func (uc *AuthUsecase) LoginByEmailPassword(ctx context.Context, user *model.User) (token string, err error) {
// 	users, err := uc.repo.ListUserByEmail(ctx, user.Email)
// 	if err != nil {
// 		return "", authv1.ErrorUserNotFound("failed to get user: %v", err)
// 	}
// 	if len(users) == 0 {
// 		uc.log.Warnf("user %s does not exist", user.Email)
// 		return "", authv1.ErrorUserNotFound("user %s does not exist", user.Email)
// 	}
// 	if !hash.BcryptCheck(user.Password, users[0].Password) {
// 		return "", authv1.ErrorIncorrectPassword("incorrect password for user: %s", user.Email)
// 	}
// 	// 登录成功，签发 token
// 	expirationTime := time.Duration(uc.cfg.Jwt.Expire) * time.Hour
// 	claims := &UserClaims{
// 		ID:   users[0].ID,
// 		Name: users[0].Name,
// 		Role: users[0].Role,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Audience:  jwt.ClaimStrings{uc.cfg.Jwt.Audience},
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			Issuer:    uc.cfg.Jwt.Issuer,
// 		},
// 	}

// 	token, err = uc.generateToken(claims)
// 	if err != nil {
// 		return "", authv1.ErrorTokenGenerationFailed("failed to generate token: %v", err)
// 	}
// 	return token, nil
// }
