package biz

import (
	"context"
	v1 "krathub/api/user/v1"
	"krathub/internal/conf"
	"krathub/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo interface {
	// SaveUser(context.Context, *model.User) (*model.User, error)
	DeleteUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
	cfg  *conf.App
}

func NewUserUsecase(repo UserRepo, logger log.Logger, cfg *conf.App) *UserUsecase {
	uc := &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
		cfg:  cfg,
	}
	return uc
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, user *model.User) (success bool, err error) {
	_, err = uc.repo.DeleteUser(ctx, user)
	if err != nil {
		return false, v1.ErrorDeleteUserFailed("failed to delete user: %v", err)
	}
	return true, nil
}
