package service

import (
	"context"
	"fmt"
	v1 "krathub/api/auth/v1"
	"krathub/internal/biz"
	"krathub/internal/data/model"
)

// AuthService is a auth service.
type AuthService struct {
	v1.UnimplementedAuthServer

	uc *biz.AuthUsecase
}

// NewAuthService new a auth service.
func NewAuthService(uc *biz.AuthUsecase) *AuthService {
	return &AuthService{uc: uc}
}

func (s *AuthService) SignupByEmail(ctx context.Context, req *v1.SignupByEmailRequest) (*v1.SignupByEmailReply, error) {
	fmt.Printf("[service] SignupByEmail: %v\n", req)
	// 调用 biz 层
	user, err := s.uc.SignupByEmail(ctx, &model.User{})
	// 拼装返回结果
	return &v1.SignupByEmailReply{
		Data: &v1.UserInfo{
			Id:    user.ID,
			Email: user.Email,
		},
	}, err
}

func (s *AuthService) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordRequest) (*v1.LoginByPasswordReply, error) {
	fmt.Printf("[service] LoginByPassword: %v\n", req)
	// 调用 biz 层
	// 拼装返回结果
	return &v1.LoginByPasswordReply{}, nil
}

func (s *AuthService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	fmt.Printf("[service] hello\n")
	return &v1.HelloReply{Message: "Hello " + in.Name}, nil
}
