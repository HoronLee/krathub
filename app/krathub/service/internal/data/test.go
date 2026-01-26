package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	sayhellopb "github.com/horonlee/krathub/api/gen/go/sayhello/service/v1"
	"github.com/horonlee/krathub/app/krathub/service/internal/biz"
	"github.com/horonlee/krathub/pkg/transport/client"
	pkglogger "github.com/horonlee/krathub/pkg/logger"
	gogrpc "google.golang.org/grpc"
)

type testRepo struct {
	data *Data
	log  *log.Helper
}

func NewTestRepo(data *Data, logger log.Logger) biz.TestRepo {
	return &testRepo{
		data: data,
		log:  log.NewHelper(pkglogger.WithModule(logger, "test/data/krathub-service")),
	}
}

func (r *testRepo) Hello(ctx context.Context, in string) (string, error) {
	r.log.Debugf("Saying hello with greeting: %s", in)

	connWrapper, err := r.data.client.CreateConn(ctx, client.GRPC, "hello")
	if err != nil {
		r.log.Errorf("Failed to create grpc connection: %v", err)
		return "", err
	}

	conn := connWrapper.Value().(gogrpc.ClientConnInterface)

	helloClient := sayhellopb.NewSayHelloClient(conn)
	ret, err := helloClient.Hello(ctx, &sayhellopb.HelloRequest{Greeting: in})
	if err != nil {
		r.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return ret.Reply, nil
}
