package biz

import (
	"context"
	authv1 "krathub/api/auth/v1"
	userv1 "krathub/api/user/v1"
	"krathub/internal/conf"
	"krathub/internal/data/model"
	"krathub/pkg/jwt"

	"github.com/go-kratos/kratos/v2/log"
)

type UserDBRepo interface {
	SaveUser(context.Context, *model.User) (*model.User, error)
	GetUserById(context.Context, int64) (*model.User, error)
	DeleteUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
}

type UserUsecase struct {
	repo    UserDBRepo
	log     *log.Helper
	cfg     *conf.App
	aDBRepo AuthDBRepo // 只依赖 AuthDBRepo
}

func NewUserUsecase(repo UserDBRepo, logger log.Logger, cfg *conf.App, aDBRepo AuthDBRepo) *UserUsecase {
	uc := &UserUsecase{
		repo:    repo,
		log:     log.NewHelper(logger),
		cfg:     cfg,
		aDBRepo: aDBRepo,
	}
	return uc
}

func (uc *UserUsecase) CurrentUserInfo(ctx context.Context) (*model.User, error) {
	// 从context中获取当前登录用户信息
	claims, ok := jwt.FromContext(ctx)
	if !ok {
		return nil, authv1.ErrorUnauthorized("user not authenticated")
	}

	// 直接使用JWT中的信息构建用户模型，避免数据库查询
	user := &model.User{
		ID:   claims.ID,
		Name: claims.Name,
		Role: claims.Role,
	}

	return user, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// 先获取当前用户信息
	currentUser, err := uc.repo.GetUserById(ctx, user.ID)
	if err != nil {
		return nil, userv1.ErrorUserNotFound("user not found: %v", err)
	}

	// 只有当用户名发生变化时才检查重复
	if user.Name != currentUser.Name {
		existingUsers, err := uc.aDBRepo.ListUserByUserName(ctx, user.Name)
		if err != nil {
			return nil, authv1.ErrorUserNotFound("failed to check username: %v", err)
		}
		if len(existingUsers) > 0 {
			return nil, authv1.ErrorUserAlreadyExists("username already exists")
		}
	}

	// 只有当邮箱发生变化时才检查重复
	if user.Email != currentUser.Email {
		existingEmails, err := uc.aDBRepo.ListUserByEmail(ctx, user.Email)
		if err != nil {
			return nil, authv1.ErrorUserNotFound("failed to check email: %v", err)
		}
		if len(existingEmails) > 0 {
			return nil, authv1.ErrorUserAlreadyExists("email already exists")
		}
	}

	updatedUser, err := uc.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, userv1.ErrorUpdateUserFailed("failed to update user: %v", err)
	}
	return updatedUser, nil
}

func (uc *UserUsecase) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := uc.checkUserExists(ctx, user); err != nil {
		return nil, err
	}

	savedUser, err := uc.repo.SaveUser(ctx, user)
	if err != nil {
		return nil, userv1.ErrorSaveUserFailed("failed to save user: %v", err)
	}
	return savedUser, nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, user *model.User) (success bool, err error) {
	_, err = uc.repo.DeleteUser(ctx, user)
	if err != nil {
		return false, userv1.ErrorDeleteUserFailed("failed to delete user: %v", err)
	}
	return true, nil
}

func (uc *UserUsecase) checkUserExists(ctx context.Context, user *model.User) error {
	if users, err := uc.aDBRepo.ListUserByUserName(ctx, user.Name); err != nil {
		return authv1.ErrorUserNotFound("failed to check username: %v", err)
	} else if len(users) > 0 {
		return authv1.ErrorUserAlreadyExists("username already exists")
	}

	if emails, err := uc.aDBRepo.ListUserByEmail(ctx, user.Email); err != nil {
		return authv1.ErrorUserNotFound("failed to check email: %v", err)
	} else if len(emails) > 0 {
		return authv1.ErrorUserAlreadyExists("email already exists")
	}
	return nil
}
