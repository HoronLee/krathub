package data

import (
	"context"
	"krathub/internal/biz"
	"krathub/internal/data/model"
	"krathub/pkg/hash"

	"github.com/go-kratos/kratos/v2/log"
)

type userDBRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserDBRepo(data *Data, logger log.Logger) biz.UserDBRepo {
	return &userDBRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// SaveUser 保存用户信息
func (r *userDBRepo) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	err := r.data.query.User.
		WithContext(ctx).
		Save(user)
	if err != nil {
		r.log.Errorf("SaveUser failed: %v", err)
		return nil, err
	}
	return user, nil
}

// GetUserById 根据用户ID获取用户信息
func (r *userDBRepo) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	user, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.ID.Eq(id)).
		First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser 删除用户
func (r *userDBRepo) DeleteUser(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := r.data.query.User.
		WithContext(ctx).
		Delete(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (r *userDBRepo) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// 判断密码是否已加密，未加密则加密
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	_, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.ID.Eq(user.ID)).
		Updates(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
