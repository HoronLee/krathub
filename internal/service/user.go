package service

import (
	"context"
	"fmt"
	authv1 "krathub/api/auth/v1"
	userv1 "krathub/api/user/v1"

	"krathub/internal/biz"
	"krathub/internal/consts"
	"krathub/internal/data/model"
)

type UserService struct {
	userv1.UnimplementedUserServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) CurrentUserInfo(ctx context.Context, req *userv1.CurrentUserInfoRequest) (*userv1.CurrentUserInfoReply, error) {
	user, err := s.uc.CurrentUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &userv1.CurrentUserInfoReply{
		Id:   user.ID,
		Name: user.Name,
		Role: user.Role,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserReply, error) {
	currentUser, err := s.uc.CurrentUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	switch currentUser.Role {
	case consts.User.String():
		if currentUser.ID != req.Id {
			return nil, authv1.ErrorUnauthorized("you only can update your own information")
		}
		if req.Role != "" && req.Role != consts.User.String() {
			return nil, authv1.ErrorUnauthorized("you do not have permission to change your role")
		}
	case consts.Admin.String():
		if req.Role != "" && req.Role >= consts.Admin.String() {
			return nil, authv1.ErrorUnauthorized("admin cannot assign role higher than admin")
		}
	case consts.Operator.String():
		if req.Role != "" && req.Role > consts.Operator.String() {
			return nil, authv1.ErrorUnauthorized("operator cannot assign role higher than operator")
		}
	}

	user := &model.User{
		ID:       req.Id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    &req.Phone,
		Avatar:   &req.Avatar,
		Bio:      &req.Bio,
		Location: &req.Location,
		Website:  &req.Website,
		Role:     req.Role,
	}
	_, err = s.uc.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &userv1.UpdateUserReply{
		Success: "true",
	}, nil
}

// SaveUser 保存用户
func (s *UserService) SaveUser(ctx context.Context, req *userv1.SaveUserRequest) (*userv1.SaveUserReply, error) {
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    &req.Phone,
		Avatar:   &req.Avatar,
		Bio:      &req.Bio,
		Location: &req.Location,
		Website:  &req.Website,
		Role:     req.Role,
	}
	user, err := s.uc.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &userv1.SaveUserReply{Id: fmt.Sprintf("%d", user.ID)}, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserReply, error) {
	success, err := s.uc.DeleteUser(ctx, &model.User{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &userv1.DeleteUserReply{Success: success}, err
}
