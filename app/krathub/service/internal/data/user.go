package data

import (
	"context"

	"github.com/horonlee/krathub/app/krathub/service/internal/biz"
	po "github.com/horonlee/krathub/app/krathub/service/internal/data/po"
	"github.com/horonlee/krathub/pkg/hash"
	pkglogger "github.com/horonlee/krathub/pkg/logger"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(pkglogger.WithModule(logger, "user/data/krathub-service")),
	}
}

// SaveUser 保存用户信息
func (r *userRepo) SaveUser(ctx context.Context, user *po.User) (*po.User, error) {
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
func (r *userRepo) GetUserById(ctx context.Context, id int64) (*po.User, error) {
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
func (r *userRepo) DeleteUser(ctx context.Context, user *po.User) (*po.User, error) {
	_, err := r.data.query.User.
		WithContext(ctx).
		Delete(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (r *userRepo) UpdateUser(ctx context.Context, user *po.User) (*po.User, error) {
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
