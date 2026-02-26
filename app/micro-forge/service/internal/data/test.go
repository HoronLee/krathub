package data

import (
	"context"

	sayhellopb "github.com/horonlee/micro-forge/api/gen/go/sayhello/service/v1"
	"github.com/horonlee/micro-forge/app/micro-forge/service/internal/biz"
	"github.com/horonlee/micro-forge/pkg/logger"
	"github.com/horonlee/micro-forge/pkg/transport/client"
)

type testRepo struct {
	data *Data
	log  *logger.Helper
}

func NewTestRepo(data *Data, l logger.Logger) biz.TestRepo {
	return &testRepo{
		data: data,
		log:  logger.NewHelper(l, logger.WithModule("test/data/micro-forge-service")),
	}
}

func (r *testRepo) Hello(ctx context.Context, in string) (string, error) {
	r.log.Debugf("Saying hello with greeting: %s", in)

	conn, err := client.GetGRPCConn(ctx, r.data.client, "sayhello.service")
	if err != nil {
		r.log.Errorf("Failed to create grpc connection: %v", err)
		return "", err
	}

	helloClient := sayhellopb.NewSayHelloServiceClient(conn)
	ret, err := helloClient.Hello(ctx, &sayhellopb.HelloRequest{Greeting: in})
	if err != nil {
		r.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return ret.GetReply(), nil
}
