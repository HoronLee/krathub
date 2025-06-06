package biz

import (
	"context"
	"fmt"
	"krathub/internal/conf"
	"krathub/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthRepo is a Auth repo.
type AuthRepo interface {
	SaveUser(context.Context, *model.User) (*model.User, error)
	DeleteUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
	ListUserByEmail(context.Context, string) ([]*model.User, error)
	// ListUserByID(context.Context, int64) (*model.User, error)
	ListUserByUserName(context.Context, string) ([]*model.User, error)
	ListUserByPhone(context.Context, string) ([]*model.User, error)
}

// AuthUsecase is a Auth usecase.
type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
	cfg  *conf.App
}

// NewAuthUsecase new an auth usecase.
func NewAuthUsecase(repo AuthRepo, logger log.Logger, cfg *conf.App) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
		cfg:  cfg,
	}
}

// SignupByEmail 使用邮件注册
func (uc *AuthUsecase) SignupByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	uc.log.WithContext(ctx).Debugf("\n[biz] Signing up user: %v", user)
	// 开发模式下，只有 admin 用户可以注册
	if user.Name == "admin" && uc.cfg.Env == "dev" {
		users, err := uc.repo.ListUserByUserName(ctx, user.Name)
		if err != nil {
			uc.log.Errorf("failed to list users by username: %v", err)
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		if len(users) > 0 {
			uc.log.Warnf("user %s already exists", user.Name)
			return nil, fmt.Errorf("user %s already exists", user.Name)
		}
		return uc.repo.SaveUser(ctx, user)
	} else {
		return nil, fmt.Errorf("only admin can sign up in dev mode")
	}

	// 发送验证码
	// 验证验证码
	// 生成用户信息
	// 保存用户信息
	return uc.repo.SaveUser(ctx, user)
}

// LoginByPassword 用户密码登录
func (uc *AuthUsecase) LoginByPassword(ctx context.Context, user *model.User) (token string, err error) {
	uc.log.WithContext(ctx).Debugf("\n[biz] Logging in user: %v", user)
	// TODO: Implement user password login logic
	if uc.cfg.Env == "dev" && (user.Email == "admin@example.com" || *user.Phone == "18888888888") {
		return "TestToken", nil
	}
	return "", fmt.Errorf("invalid credentials")
}
