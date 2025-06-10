package service

import (
	"context"
	userv1 "krathub/api/user/v1"
	"krathub/internal/biz"
	"krathub/internal/data/model"
)

type UserService struct {
	userv1.UnimplementedUserServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}
func (s *UserService) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserReply, error) {
	success, err := s.uc.DeleteUser(ctx, &model.User{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &userv1.DeleteUserReply{Success: success}, err
}
