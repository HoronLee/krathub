package data

import (
	"context"

	sayhellov1 "github.com/horonlee/krathub/api/sayhello/v1"
	"github.com/horonlee/krathub/internal/biz"
	"github.com/horonlee/krathub/internal/data/model"
	"github.com/horonlee/krathub/pkg/hash"

	"github.com/go-kratos/kratos/v2/log"
)

// authRepo 统一的认证仓库实现，同时包含数据库和 grpc 操作
type authRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 数据库操作方法

func (r *authRepo) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
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

func (r *authRepo) ListUserByUserName(ctx context.Context, name string) ([]*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.Name.Eq(name)).Find()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepo) ListUserByEmail(ctx context.Context, email string) ([]*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.Email.Eq(email)).Find()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Grpc 操作方法

// Hello 负责调用 hello 服务的 SayHello 方法
func (r *authRepo) Hello(ctx context.Context, in string) (string, error) {
	r.log.Debugf("Saying hello with greeting: %s", in)

	// 直接使用 CreateGrpcConn 方法获取连接
	conn, err := r.data.clientFactory.CreateGrpcConn(ctx, "hello")
	if err != nil {
		r.log.Errorf("Failed to create grpc connection: %v", err)
		return "", err
	}

	// 使用连接创建客户端
	helloClient := sayhellov1.NewSayHelloClient(conn)

	ret, err := helloClient.Hello(ctx, &sayhellov1.HelloRequest{Greeting: &in})
	if err != nil {
		r.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return ret.Reply, nil
}
