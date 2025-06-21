package service

import (
	"context"
	"fmt"
	authv1 "krathub/api/auth/v1"
	"krathub/internal/biz"
	"krathub/internal/data/model"

	"github.com/fatedier/golib/log"
)

// AuthService is a auth service.
type AuthService struct {
	authv1.UnimplementedAuthServer

	uc *biz.AuthUsecase
}

// NewAuthService new a auth service.
func NewAuthService(uc *biz.AuthUsecase) *AuthService {
	return &AuthService{uc: uc}
}

func (s *AuthService) SignupByEmail(ctx context.Context, req *authv1.SignupByEmailRequest) (*authv1.SignupByEmailReply, error) {
	// 参数校验
	if req.Password != req.PasswordConfirm {
		return nil, fmt.Errorf("password and confirm password do not match")
	}
	// 调用 biz 层
	user, err := s.uc.SignupByEmail(ctx, &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	// 拼装返回结果
	return &authv1.SignupByEmailReply{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

// LoginByEmailPassword user login by email and password.
func (s *AuthService) LoginByEmailPassword(ctx context.Context, req *authv1.LoginByEmailPasswordRequest) (*authv1.LoginByEmailPasswordReply, error) {
	user := &model.User{
		Email:    req.LoginId,
		Password: req.Password,
	}
	token, err := s.uc.LoginByEmailPassword(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("login by email password failed: %w", err)
	}
	return &authv1.LoginByEmailPasswordReply{
		Token: token,
	}, nil
}

// SayHello 实现 authv1.AuthServer 接口的 SayHello 方法
func (s *AuthService) SayHello(ctx context.Context, req *authv1.HelloRequest) (*authv1.HelloResponse, error) {
	log.Debugf("Received SayHello request with greeting: %v", req.Greeting)
	// 调用 biz 层
	res, err := s.uc.SayHello(ctx, req.Greeting)
	if err != nil {
		return nil, err
	}
	// 拼装返回响应
	return &authv1.HelloResponse{
		Reply: res,
	}, nil
}
