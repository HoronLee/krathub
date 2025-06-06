package service

import (
	"context"
	"fmt"
	v1 "krathub/api/auth/v1"
	"krathub/internal/biz"
	"krathub/internal/data/model"
	"krathub/pkg/helper"
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
		Data: &v1.UserInfo{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Token: "待实现",
	}, nil
}

// LoginByPassword user login by password.
func (s *AuthService) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordRequest) (*v1.LoginByPasswordReply, error) {
	// 参数校验
	user := &model.User{}
	if helper.IsEmail(req.LoginId) {
		user.Email = req.LoginId
	} else if helper.IsPhone(req.LoginId) {
		user.Phone = &req.LoginId
	} else {
		return nil, fmt.Errorf("login_id must be email or phone")
	}
	user.Password = req.Password
	// 调用 biz 层
	token, err := s.uc.LoginByPassword(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}
	return &v1.LoginByPasswordReply{
		Token: token,
	}, nil
}
