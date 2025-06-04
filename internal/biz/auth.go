package biz

import (
	"context"
	"krathub/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthRepo is a Auth repo.
type AuthRepo interface {
	SaveUser(context.Context, *model.User) (*model.User, error)
	// DeleteUser(context.Context, *model.User) (*model.User, error)
	// UpdateUser(context.Context, *model.User) (*model.User, error)
	// ListUserByEmail(context.Context, string) ([]*model.User, error)
	// ListUserByID(context.Context, int64) (*model.User, error)
	// ListUserByUsername(context.Context, string) ([]*model.User, error)
	// ListUserByPhone(context.Context, string) ([]*model.User, error)
}

// AuthUsecase is a Auth usecase.
type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

// NewAuthUsecase new an auth usecase.
func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *AuthUsecase) SignupByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	uc.log.WithContext(ctx).Debugf("[biz] Signing up user: %v", user)
	// 1. 数据校验
	// 2. 发送验证码
	// 3. 验证验证码
	// 4. 生成用户信息
	// 5. 保存用户信息
	return uc.repo.SaveUser(ctx, user)
}
