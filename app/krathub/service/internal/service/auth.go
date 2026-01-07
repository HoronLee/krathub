package service

import (
	"context"
	"fmt"

	authv1 "github.com/horonlee/krathub/api/gen/go/auth/service/v1"
	"github.com/horonlee/krathub/app/krathub/service/internal/biz"
	po "github.com/horonlee/krathub/app/krathub/service/internal/data/po"

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
	user, err := s.uc.SignupByEmail(ctx, &po.User{
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
	user := &po.User{
		Email:    req.Email,
		Password: req.Password,
	}
	tokenPair, err := s.uc.LoginByEmailPassword(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("login by email password failed: %w", err)
	}
	return &authv1.LoginByEmailPasswordReply{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// RefreshToken refreshes the access token using a valid refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenReply, error) {
	tokenPair, err := s.uc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token failed: %w", err)
	}
	return &authv1.RefreshTokenReply{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Logout invalidates the refresh token
func (s *AuthService) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutReply, error) {
	err := s.uc.Logout(ctx, req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("logout failed: %w", err)
	}
	return &authv1.LogoutReply{
		Success: true,
	}, nil
}

// SayHello 实现 authv1.AuthServer 接口的 SayHello 方法
func (s *AuthService) Hello(ctx context.Context, req *authv1.HelloRequest) (*authv1.HelloResponse, error) {
	log.Debugf("Received SayHello request with greeting: %v", req.Req)
	// 调用 biz 层
	res, err := s.uc.Hello(ctx, req.Req)
	if err != nil {
		return nil, err
	}
	// 拼装返回响应
	return &authv1.HelloResponse{
		Rep: res,
	}, nil
}
