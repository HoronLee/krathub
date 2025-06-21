package biz

import (
	"context"
	authv1 "krathub/api/auth/v1"
	userv1 "krathub/api/user/v1"
	"krathub/internal/conf"
	"krathub/internal/data/model"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
)

type UserRepo interface {
	SaveUser(context.Context, *model.User) (*model.User, error)
	GetUserById(context.Context, int64) (*model.User, error)
	DeleteUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)
}

type UserUsecase struct {
	repo  UserRepo
	log   *log.Helper
	cfg   *conf.App
	aRepo AuthRepo // 引入auth包的repo接口
}

func NewUserUsecase(repo UserRepo, logger log.Logger, cfg *conf.App, aRepo AuthRepo) *UserUsecase {
	uc := &UserUsecase{
		repo:  repo,
		log:   log.NewHelper(logger),
		cfg:   cfg,
		aRepo: aRepo,
	}
	return uc
}

func (uc *UserUsecase) CurrentUserInfo(ctx context.Context) (*model.User, error) {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, authv1.ErrorMissingToken("missing transport context")
	}
	authHeader := tr.RequestHeader().Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	if tokenString == "" {
		return nil, authv1.ErrorMissingToken("missing Authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(uc.cfg.Jwt.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, authv1.ErrorUnauthorized("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, authv1.ErrorInvalidTokenType("invalid claims")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return nil, authv1.ErrorInvalidTokenType("id not found in token")
	}
	userID := int64(idFloat)

	user, err := uc.repo.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := uc.repo.GetUserById(ctx, user.ID)
	if err != nil {
		return nil, userv1.ErrorUserNotFound("user not found: %v", err)
	}

	existingUsers, err := uc.aRepo.ListUserByUserName(ctx, user.Name)
	if err != nil {
		return nil, authv1.ErrorUserNotFound("failed to check username: %v", err)
	}
	if len(existingUsers) > 0 {
		return nil, authv1.ErrorUserAlreadyExists("username already exists")
	}

	updatedUser, err := uc.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, userv1.ErrorUpdateUserFailed("failed to update user: %v", err)
	}
	return updatedUser, nil
}

func (uc *UserUsecase) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := uc.checkUserExists(ctx, user); err != nil {
		return nil, err
	}

	savedUser, err := uc.repo.SaveUser(ctx, user)
	if err != nil {
		return nil, userv1.ErrorSaveUserFailed("failed to save user: %v", err)
	}
	return savedUser, nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, user *model.User) (success bool, err error) {
	_, err = uc.repo.DeleteUser(ctx, user)
	if err != nil {
		return false, userv1.ErrorDeleteUserFailed("failed to delete user: %v", err)
	}
	return true, nil
}

func (uc *UserUsecase) checkUserExists(ctx context.Context, user *model.User) error {
	if users, err := uc.aRepo.ListUserByUserName(ctx, user.Name); err != nil {
		return authv1.ErrorUserNotFound("failed to check username: %v", err)
	} else if len(users) > 0 {
		return authv1.ErrorUserAlreadyExists("username already exists")
	}

	if emails, err := uc.aRepo.ListUserByEmail(ctx, user.Email); err != nil {
		return authv1.ErrorUserNotFound("failed to check email: %v", err)
	} else if len(emails) > 0 {
		return authv1.ErrorUserAlreadyExists("email already exists")
	}
	return nil
}
