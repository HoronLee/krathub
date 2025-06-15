package service

import (
	"context"
	"fmt"
	authv1 "krathub/api/auth/v1"
	"krathub/internal/biz"
	"krathub/internal/data/model"

	"github.com/fatedier/golib/log"
	"github.com/go-kratos/kratos/v2/metadata"
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

// SayHello implements helloworld.GreeterServer.
func (s *AuthService) SayHello(ctx context.Context, in *authv1.HelloRequest) (*authv1.HelloReply, error) {
	var username, role string
	if md, ok := metadata.FromServerContext(ctx); ok {
		username = md.Get("username")
		role = md.Get("role")
		log.Debugf("User %s with role %s is logging in", username, role)
	} else {
		log.Debugf("No metadata found in context")
	}
	log.Infof("User %s with role %s is logging in", username, role)

	return &authv1.HelloReply{Message: "Hello World"}, nil
}
