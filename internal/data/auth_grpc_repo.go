package data

import (
	"context"
	hellov1 "krathub/api/hello/v1"
	"krathub/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type authGrpcRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthGrpcRepo(data *Data, logger log.Logger) biz.AuthGrpcRepo {
	return &authGrpcRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *authGrpcRepo) Hello(ctx context.Context, in string) (string, error) {
	r.log.Debugf("Saying hello with greeting: %s", in)

	// 直接使用 CreateGrpcConn 方法获取连接
	conn, err := r.data.clientFactory.CreateGrpcConn(ctx, "hello.service")
	if err != nil {
		r.log.Errorf("Failed to create grpc connection: %v", err)
		return "", err
	}

	// 使用连接创建客户端
	helloClient := hellov1.NewHelloServiceClient(conn)

	ret, err := helloClient.SayHello(ctx, &hellov1.HelloRequest{Greeting: &in})
	if err != nil {
		r.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return ret.Reply, nil
}
