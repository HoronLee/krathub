package service

import (
	"context"
	v1 "krathub/api/user/v1"
	"krathub/internal/biz"
	"krathub/internal/data/model"
)

type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}
func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	success, err := s.uc.DeleteUser(ctx, &model.User{
		ID: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &v1.DeleteUserReply{Success: success}, err
}
