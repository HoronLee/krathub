package data

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sayhellov1 "github.com/horonlee/krathub/api/sayhello/v1"
	"github.com/horonlee/krathub/internal/biz"
	"github.com/horonlee/krathub/internal/client"
	"github.com/horonlee/krathub/internal/data/model"
	"github.com/horonlee/krathub/pkg/hash"

	"github.com/go-kratos/kratos/v2/log"
	gogrpc "google.golang.org/grpc"
)

// authRepo 统一的认证仓库实现，同时包含数据库和 grpc 操作
type authRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 数据库操作方法

func (r *authRepo) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	if !hash.BcryptIsHashed(user.Password) {
		bcryptPassword, err := hash.BcryptHash(user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = bcryptPassword
	}
	err := r.data.query.User.WithContext(ctx).Create(user)
	if err != nil {
		r.log.Errorf("SaveUser failed: %v", err)
		return nil, err
	}
	return user, nil
}

func (r *authRepo) GetUserByUserName(ctx context.Context, name string) (*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.Name.Eq(name)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.Email.Eq(email)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepo) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := r.data.query.User.WithContext(ctx).Where(r.data.query.User.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Grpc 操作方法

// Hello 负责调用 hello 服务的 SayHello 方法
func (r *authRepo) Hello(ctx context.Context, in string) (string, error) {
	r.log.Debugf("Saying hello with greeting: %s", in)

	// 使用新的 CreateConn 方法获取连接
	connWrapper, err := r.data.client.CreateConn(ctx, client.GRPC, "hello")
	if err != nil {
		r.log.Errorf("Failed to create grpc connection: %v", err)
		return "", err
	}

	// 获取原始gRPC连接
	conn := connWrapper.Value().(gogrpc.ClientConnInterface)

	// 使用连接创建客户端
	helloClient := sayhellov1.NewSayHelloClient(conn)

	ret, err := helloClient.Hello(ctx, &sayhellov1.HelloRequest{Greeting: &in})
	if err != nil {
		r.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return ret.Reply, nil
}

// TokenStore methods implementation

// SaveRefreshToken 保存Refresh Token到Redis
func (r *authRepo) SaveRefreshToken(ctx context.Context, userID int64, token string, expiration time.Duration) error {
	// 存储refresh token -> user_id的映射
	tokenKey := fmt.Sprintf("refresh_token:%s", token)
	if err := r.data.redis.Set(ctx, tokenKey, strconv.FormatInt(userID, 10), expiration); err != nil {
		r.log.Errorf("Failed to save refresh token: %v", err)
		return err
	}

	// 将token添加到用户的token集合中，用于批量删除
	userTokensKey := fmt.Sprintf("user_tokens:%d", userID)
	if err := r.data.redis.SAdd(ctx, userTokensKey, token); err != nil {
		r.log.Errorf("Failed to add token to user set: %v", err)
		return err
	}

	// 为用户token集合设置过期时间
	if err := r.data.redis.Expire(ctx, userTokensKey, expiration); err != nil {
		r.log.Errorf("Failed to set expiration for user tokens set: %v", err)
		return err
	}

	return nil
}

// GetRefreshToken 获取Refresh Token关联的用户ID
func (r *authRepo) GetRefreshToken(ctx context.Context, token string) (int64, error) {
	tokenKey := fmt.Sprintf("refresh_token:%s", token)
	userIDStr, err := r.data.redis.Get(ctx, tokenKey)
	if err != nil {
		r.log.Errorf("Failed to get refresh token: %v", err)
		return 0, err
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		r.log.Errorf("Failed to parse user ID: %v", err)
		return 0, err
	}

	return userID, nil
}

// DeleteRefreshToken 删除Refresh Token
func (r *authRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	// 首先获取用户ID，以便从用户token集合中删除
	userID, err := r.GetRefreshToken(ctx, token)
	if err != nil {
		// 如果token不存在，也认为删除成功
		r.log.Warnf("Token not found during deletion: %v", err)
		return nil
	}

	// 删除token -> user_id的映射
	tokenKey := fmt.Sprintf("refresh_token:%s", token)
	if err := r.data.redis.Del(ctx, tokenKey); err != nil {
		r.log.Errorf("Failed to delete refresh token: %v", err)
		return err
	}

	// 从用户token集合中删除该token
	userTokensKey := fmt.Sprintf("user_tokens:%d", userID)
	// 获取集合中的所有token
	tokens, err := r.data.redis.SMembers(ctx, userTokensKey)
	if err != nil {
		r.log.Errorf("Failed to get user tokens: %v", err)
		return err
	}

	// 重新创建集合，排除要删除的token
	if err := r.data.redis.Del(ctx, userTokensKey); err != nil {
		r.log.Errorf("Failed to delete user tokens set: %v", err)
		return err
	}

	// 重新添加除了要删除的token之外的所有token
	for _, t := range tokens {
		if t != token {
			if err := r.data.redis.SAdd(ctx, userTokensKey, t); err != nil {
				r.log.Errorf("Failed to re-add token to user set: %v", err)
				return err
			}
		}
	}

	return nil
}

// DeleteUserRefreshTokens 删除用户所有Refresh Token
func (r *authRepo) DeleteUserRefreshTokens(ctx context.Context, userID int64) error {
	userTokensKey := fmt.Sprintf("user_tokens:%d", userID)

	// 获取用户的所有token
	tokens, err := r.data.redis.SMembers(ctx, userTokensKey)
	if err != nil {
		r.log.Errorf("Failed to get user tokens: %v", err)
		return err
	}

	// 删除每个token的映射
	for _, token := range tokens {
		tokenKey := fmt.Sprintf("refresh_token:%s", token)
		if err := r.data.redis.Del(ctx, tokenKey); err != nil {
			r.log.Errorf("Failed to delete token %s: %v", token, err)
			return err
		}
	}

	// 删除用户token集合
	if err := r.data.redis.Del(ctx, userTokensKey); err != nil {
		r.log.Errorf("Failed to delete user tokens set: %v", err)
		return err
	}

	return nil
}
