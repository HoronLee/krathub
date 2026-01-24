package service

import (
	"context"

	testv1 "github.com/horonlee/krathub/api/gen/go/test/service/v1"
	"github.com/horonlee/krathub/app/krathub/service/internal/biz"
)

// TestService is a test service.
type TestService struct {
	testv1.UnimplementedTestServiceServer

	uc *biz.TestUsecase
}

// NewTestService new a test service.
func NewTestService(uc *biz.TestUsecase) *TestService {
	return &TestService{uc: uc}
}

// Hello calls the hello service
func (s *TestService) Hello(ctx context.Context, req *testv1.HelloRequest) (*testv1.HelloResponse, error) {
	// 调用 biz 层
	res, err := s.uc.Hello(ctx, req.Req)
	if err != nil {
		return nil, err
	}
	// 拼装返回响应
	return &testv1.HelloResponse{
		Rep: res,
	}, nil
}

// Test is a test method.
func (s *TestService) Test(ctx context.Context, req *testv1.TestRequest) (*testv1.TestResponse, error) {
	return &testv1.TestResponse{Message: "公开的测试路由"}, nil
}

// PrivateTest is a private test method.
func (s *TestService) PrivateTest(ctx context.Context, req *testv1.PrivateTestRequest) (*testv1.PrivateTestResponse, error) {
	return &testv1.PrivateTestResponse{Message: "私有的测试路由"}, nil
}
