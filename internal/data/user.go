package data

import (
	"context"
	"krathub/internal/biz"
	"krathub/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// DeleteUser 删除用户
func (r *userRepo) DeleteUser(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := r.data.query.User.
		WithContext(ctx).
		Delete(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (r *authRepo) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.ID.Eq(user.ID)).
		Updates(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
