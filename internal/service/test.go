package service

import (
	"context"

	testv1 "github.com/horonlee/krathub/api/test/v1"
)

// AuthService is a auth service.
type TestService struct {
	testv1.UnimplementedTestServer
}

// NewAuthService new a auth service.
func NewTestService() *TestService {
	return &TestService{}
}

// Test is a test method.
func (s *TestService) Test(ctx context.Context, req *testv1.TestRequest) (*testv1.TestResponse, error) {
	return &testv1.TestResponse{Message: "公开的测试路由"}, nil
}

// PrivateTest is a private test method.
func (s *TestService) PrivateTest(ctx context.Context, req *testv1.PrivateTestRequest) (*testv1.PrivateTestResponse, error) {
	return &testv1.PrivateTestResponse{Message: "私有的测试路由"}, nil
}
