package data

import (
	"context"
	"krathub/internal/biz"
	"krathub/internal/data/model"
	"krathub/pkg/hash"

	"github.com/go-kratos/kratos/v2/log"
)

type authRepo struct {
	data *Data
	log  *log.Helper
}

// NewAuthRepo .
func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// SaveUser 保存用户信息
func (r *authRepo) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	r.log.Debugf("[data] SaveUser called: %+v", user)
	bcryptPassword, err := hash.BcryptHash(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = bcryptPassword
	r.log.Debugf("[data] bcryptPassword: %v", user.Password)
	err = r.data.query.User.
		WithContext(ctx).
		Save(user)
	if err != nil {
		r.log.Errorf("SaveUser failed: %v", err)
	}
	return user, err
}

// DeleteUser 删除用户信息
func (r *authRepo) DeleteUser(ctx context.Context, user *model.User) (*model.User, error) {
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

// ListUserByUserName 根据用户名获取用户信息
func (r *authRepo) ListUserByUserName(ctx context.Context, name string) ([]*model.User, error) {
	user, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.Name.Eq(name)).
		Find()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ListUserByEmail 根据邮箱获取用户信息
func (r *authRepo) ListUserByEmail(ctx context.Context, email string) ([]*model.User, error) {
	user, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.Email.Eq(email)).
		Find()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ListUserByPhone 根据手机号获取用户信息
func (r *authRepo) ListUserByPhone(ctx context.Context, phone string) ([]*model.User, error) {
	user, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.Phone.Eq(phone)).
		Find()
	if err != nil {
		return nil, err
	}
	return user, nil
}
