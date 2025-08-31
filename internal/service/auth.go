package service

import (
	"context"

	callhellov1 "github.com/horonlee/krathub/api/callhello/v1"
	"github.com/horonlee/krathub/internal/biz"

	"github.com/fatedier/golib/log"
)

// CallHelloService is a auth service.
type CallHelloService struct {
	callhellov1.UnimplementedCallHelloServer

	uc *biz.CallHelloUsecase
}

// NewHelloService new a auth service.
func NewHelloService(uc *biz.CallHelloUsecase) *CallHelloService {
	return &CallHelloService{uc: uc}
}

// SayHello 实现 authv1.AuthServer 接口的 SayHello 方法
func (s *CallHelloService) Hello(ctx context.Context, req *callhellov1.HelloRequest) (*callhellov1.HelloResponse, error) {
	log.Debugf("Received SayHello request with greeting: %v", req.Req)
	// 调用 biz 层
	res, err := s.uc.Hello(ctx, req.Req)
	if err != nil {
		return nil, err
	}
	// 拼装返回响应
	return &callhellov1.HelloResponse{
		Rep: res,
	}, nil
}
