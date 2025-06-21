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
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	err := r.data.query.User.
		WithContext(ctx).
		Create(user)
	if err != nil {
		r.log.Errorf("SaveUser failed: %v", err)
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
