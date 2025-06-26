package data

import (
	"context"
	"krathub/internal/biz"
	"krathub/internal/data/model"
	"krathub/pkg/hash"

	"github.com/go-kratos/kratos/v2/log"
)

type authDBRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthDBRepo(data *Data, logger log.Logger) biz.AuthDBRepo {
	return &authDBRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *authDBRepo) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	err := r.data.query.User.WithContext(ctx).Create(user)
	if err != nil {
		r.log.Errorf("SaveUser failed: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *authDBRepo) ListUserByUserName(ctx context.Context, name string) ([]*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.Name.Eq(name)).Find()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authDBRepo) ListUserByEmail(ctx context.Context, email string) ([]*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.Email.Eq(email)).Find()
	if err != nil {
		return nil, err
	}
	return user, nil
}
