package biz

import (
	"context"
	authv1 "krathub/api/auth/v1"
	"krathub/internal/conf"
	"krathub/internal/data/model"
	"krathub/pkg/hash"
	jwtpkg "krathub/pkg/jwt"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthRepo 统一的认证仓库接口，包含数据库和 grpc 操作
type AuthRepo interface {
	// 数据库操作
	SaveUser(context.Context, *model.User) (*model.User, error)
	ListUserByEmail(context.Context, string) ([]*model.User, error)
	ListUserByUserName(context.Context, string) ([]*model.User, error)
	// grpc 操作
	Hello(ctx context.Context, in string) (string, error)
}

// AuthUsecase is a Auth usecase.
type AuthUsecase struct {
	repo            AuthRepo
	log             *log.Helper
	cfg             *conf.App
	adminRegistered bool        // 是否已经注册了 admin 用户
	jwt             *jwtpkg.JWT // 新增 jwt 字段
}

// NewAuthUsecase new an auth usecase.
func NewAuthUsecase(repo AuthRepo, logger log.Logger, cfg *conf.App) *AuthUsecase {
	uc := &AuthUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
		cfg:  cfg,
		jwt: jwtpkg.NewJWT(&jwtpkg.Config{
			SecretKey: cfg.Jwt.SecretKey,
			Expire:    cfg.Jwt.Expire,
			Audience:  cfg.Jwt.Audience,
			Issuer:    cfg.Jwt.Issuer,
		}),
	}
	// 初始化 adminRegistered
	admins, err := repo.ListUserByUserName(context.Background(), "admin")
	if err == nil && len(admins) > 0 {
		uc.adminRegistered = true
	}
	return uc
}

// SignupByEmail 使用邮件注册
func (uc *AuthUsecase) SignupByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	// 检查 admin 用户是否已存在
	if !uc.adminRegistered {
		// 第一次注册，用户名必须为 admin
		if user.Name != "admin" {
			return nil, authv1.ErrorInvalidCredentials("the first user must be named admin")
		}
		user.Role = "admin"
	} else {
		// 后续注册，用户名可以任意，但角色为 user
		// 检查用户名是否已存在
		existingUsers, err := uc.repo.ListUserByUserName(ctx, user.Name)
		if err != nil {
			return nil, authv1.ErrorUserNotFound("failed to check username: %v", err)
		}
		if len(existingUsers) > 0 {
			return nil, authv1.ErrorUserAlreadyExists("username already exists")
		}
		user.Role = "user"
	}

	// 检查邮箱是否已存在
	existingEmails, err := uc.repo.ListUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, authv1.ErrorUserNotFound("failed to check email: %v", err)
	}
	if len(existingEmails) > 0 {
		return nil, authv1.ErrorUserAlreadyExists("email already exists")
	}

	createdUser, err := uc.repo.SaveUser(ctx, user)
	if err == nil && !uc.adminRegistered && user.Name == "admin" {
		uc.adminRegistered = true // 注册成功后更新状态
	}
	return createdUser, err
}

// generateToken 签发 JWT token
func (uc *AuthUsecase) generateToken(id int64, name, role string) (string, error) {
	return uc.jwt.GenerateToken(id, name, role)
}

// LoginByEmailPassword 邮箱密码登录
func (uc *AuthUsecase) LoginByEmailPassword(ctx context.Context, user *model.User) (token string, err error) {
	users, err := uc.repo.ListUserByEmail(ctx, user.Email)
	if err != nil {
		return "", authv1.ErrorUserNotFound("failed to get user: %v", err)
	}
	if len(users) == 0 {
		uc.log.Warnf("user %s does not exist", user.Email)
		return "", authv1.ErrorUserNotFound("user %s does not exist", user.Email)
	}
	if !hash.BcryptCheck(user.Password, users[0].Password) {
		return "", authv1.ErrorIncorrectPassword("incorrect password for user: %s", user.Email)
	}
	// 登录成功，签发 token
	token, err = uc.generateToken(users[0].ID, users[0].Name, users[0].Role)
	if err != nil {
		return "", authv1.ErrorTokenGenerationFailed("failed to generate token: %v", err)
	}
	return token, nil
}

// Hello 通过 repo 实现
func (uc *AuthUsecase) Hello(ctx context.Context, in *string) (string, error) {
	uc.log.Debugf("Saying hello with greeting: %s", *in)
	response, err := uc.repo.Hello(ctx, *in)
	if err != nil {
		uc.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return response, nil
}
