package data

import (
	"context"
	"krathub/internal/biz"
	"krathub/internal/data/model"

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

func (r *authRepo) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.data.query.User.
		WithContext(ctx).
		Save(user)
	return user, err
}
