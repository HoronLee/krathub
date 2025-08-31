package service

import (
	"context"

	hellov1 "github.com/horonlee/krathub/api/auth/v1"
	"github.com/horonlee/krathub/internal/biz"

	"github.com/fatedier/golib/log"
)

// HelloService is a auth service.
type HelloService struct {
	hellov1.UnimplementedCallHelloServer

	uc *biz.HelloUsecase
}

// NewHelloService new a auth service.
func NewHelloService(uc *biz.HelloUsecase) *HelloService {
	return &HelloService{uc: uc}
}

// SayHello 实现 authv1.AuthServer 接口的 SayHello 方法
func (s *HelloService) Hello(ctx context.Context, req *hellov1.HelloRequest) (*hellov1.HelloResponse, error) {
	log.Debugf("Received SayHello request with greeting: %v", req.Req)
	// 调用 biz 层
	res, err := s.uc.Hello(ctx, req.Req)
	if err != nil {
		return nil, err
	}
	// 拼装返回响应
	return &hellov1.HelloResponse{
		Rep: res,
	}, nil
}
