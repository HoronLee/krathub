package service

import (
	"context"

	sayhellov1 "github.com/horonlee/krathub/api/sayhello/v1"
	"github.com/horonlee/krathub/internal/biz"

	"github.com/fatedier/golib/log"
)

// SayHelloService is a auth service.
type SayHelloService struct {
	sayhellov1.UnimplementedSayHelloServer

	uc *biz.SayHelloUsecase
}

// NewSayHelloService new a auth service.
func NewSayHelloService(uc *biz.SayHelloUsecase) *SayHelloService {
	return &SayHelloService{uc: uc}
}

func (s *SayHelloService) Hello(ctx context.Context, req *sayhellov1.HelloRequest) (*sayhellov1.HelloResponse, error) {
	log.Debugf("Received SayHello request with greeting: %v", req.Greeting)
	// 调用 biz 层
	res, err := s.uc.Hello(ctx, req.Greeting)
	if err != nil {
		return nil, err
	}
	// 拼装返回响应
	return &sayhellov1.HelloResponse{
		Reply: res,
	}, nil
}
