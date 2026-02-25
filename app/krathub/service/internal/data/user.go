package data

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/horonlee/krathub/app/krathub/service/internal/biz"
	"github.com/horonlee/krathub/app/krathub/service/internal/biz/entity"
	dataent "github.com/horonlee/krathub/app/krathub/service/internal/data/ent"
	entuser "github.com/horonlee/krathub/app/krathub/service/internal/data/ent/user"
	"github.com/horonlee/krathub/pkg/helpers/hash"
	pkglogger "github.com/horonlee/krathub/pkg/logger"
	"github.com/horonlee/krathub/pkg/mapper"
)

type userRepo struct {
	data   *Data
	log    *log.Helper
	mapper *mapper.CopierMapper[entity.User, dataent.User]
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:   data,
		log:    log.NewHelper(pkglogger.With(logger, pkglogger.WithModule("user/data/krathub-service"))),
		mapper: mapper.New[entity.User, dataent.User]().RegisterConverters(mapper.AllBuiltinConverters()),
	}
}

func (r *userRepo) SaveUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	entUser := r.mapper.ToEntity(user)
	b := r.data.entClient.User.Create().
		SetName(entUser.Name).
		SetEmail(entUser.Email).
		SetPassword(entUser.Password).
		SetNillablePhone(entUser.Phone).
		SetNillableAvatar(entUser.Avatar).
		SetNillableBio(entUser.Bio).
		SetNillableLocation(entUser.Location).
		SetNillableWebsite(entUser.Website).
		SetRole(entUser.Role)

	if entUser.ID > 0 {
		b.SetID(entUser.ID)
	}

	created, err := b.Save(ctx)
	if err != nil {
		r.log.Errorf("SaveUser failed: %v", err)
		return nil, err
	}
	return r.mapper.ToDomain(created), nil
}

func (r *userRepo) GetUserById(ctx context.Context, id int64) (*entity.User, error) {
	entUser, err := r.data.entClient.User.Query().Where(entuser.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(entUser), nil
}

func (r *userRepo) DeleteUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := r.data.entClient.User.DeleteOneID(user.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	entUser := r.mapper.ToEntity(user)
	updated, err := r.data.entClient.User.UpdateOneID(user.ID).
		SetName(entUser.Name).
		SetEmail(entUser.Email).
		SetPassword(entUser.Password).
		SetNillablePhone(entUser.Phone).
		SetNillableAvatar(entUser.Avatar).
		SetNillableBio(entUser.Bio).
		SetNillableLocation(entUser.Location).
		SetNillableWebsite(entUser.Website).
		SetRole(entUser.Role).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapper.ToDomain(updated), nil
}

func (r *userRepo) ListUsers(ctx context.Context, page int32, pageSize int32) ([]*entity.User, int64, error) {
	offset := int((page - 1) * pageSize)
	limit := int(pageSize)

	query := r.data.entClient.User.Query().Order(entuser.ByID(sql.OrderDesc()))
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	entUsers, err := query.Offset(offset).Limit(limit).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, 0, len(entUsers))
	for _, entUser := range entUsers {
		users = append(users, r.mapper.ToDomain(entUser))
	}

	return users, int64(total), nil
}
