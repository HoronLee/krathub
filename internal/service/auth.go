package service

import (
	"context"
	"fmt"
	v1 "krathub/api/auth/v1"
	"krathub/internal/biz"
	"krathub/internal/data/model"
	"krathub/pkg/helper"

	"github.com/fatedier/golib/log"
	"github.com/go-kratos/kratos/v2/metadata"
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
	return &v1.SignupByEmailReply{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

// LoginByPassword user login by password.
func (s *AuthService) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordRequest) (*v1.LoginByPasswordReply, error) {
	user := &model.User{}
	var token string
	var err error
	if helper.IsEmail(req.LoginId) {
		user.Email = req.LoginId
		user.Password = req.Password
		token, err = s.uc.LoginByEmailPassword(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("login by email password failed: %w", err)
		}
	} else if helper.IsPhone(req.LoginId) {
		user.Phone = &req.LoginId
		user.Password = req.Password
		token, err = s.uc.LoginByPhonePassword(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("login by phone password failed: %w", err)
		}
	} else {
		return nil, fmt.Errorf("login_id must be email or phone")
	}
	user.Password = req.Password
	return &v1.LoginByPasswordReply{
		Token: token,
	}, nil
}

// SayHello implements helloworld.GreeterServer.
func (s *AuthService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	var username, role string
	if md, ok := metadata.FromServerContext(ctx); ok {
		username = md.Get("username")
		role = md.Get("role")
		log.Debugf("User %s with role %s is logging in", username, role)
	} else {
		log.Debugf("No metadata found in context")
	}
	log.Infof("User %s with role %s is logging in", username, role)

	return &v1.HelloReply{Message: "Hello World"}, nil
}
