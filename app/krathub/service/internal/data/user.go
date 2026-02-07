package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/horonlee/krathub/app/krathub/service/internal/biz"
	po "github.com/horonlee/krathub/app/krathub/service/internal/data/po"
	"github.com/horonlee/krathub/pkg/helpers/hash"
	pkglogger "github.com/horonlee/krathub/pkg/logger"
	"github.com/horonlee/krathub/pkg/mapper"
)

type userRepo struct {
	data   *Data
	log    *log.Helper
	mapper *mapper.CopierMapper[biz.User, po.User]
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:   data,
		log:    log.NewHelper(pkglogger.WithModule(logger, "user/data/krathub-service")),
		mapper: mapper.New[biz.User, po.User]().RegisterConverters(mapper.AllBuiltinConverters()),
	}
}

func (r *userRepo) SaveUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	poUser := r.mapper.ToEntity(user)
	err := r.data.query.User.
		WithContext(ctx).
		Save(poUser)
	if err != nil {
		r.log.Errorf("SaveUser failed: %v", err)
		return nil, err
	}
	return r.mapper.ToDomain(poUser), nil
}

func (r *userRepo) GetUserById(ctx context.Context, id int64) (*biz.User, error) {
	poUser, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.ID.Eq(id)).
		First()
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(poUser), nil
}

func (r *userRepo) DeleteUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	poUser := r.mapper.ToEntity(user)
	_, err := r.data.query.User.
		WithContext(ctx).
		Delete(poUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (*biz.User, error) {
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	poUser := r.mapper.ToEntity(user)
	_, err := r.data.query.User.
		WithContext(ctx).
		Where(r.data.query.User.ID.Eq(user.ID)).
		Updates(poUser)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(poUser), nil
}
