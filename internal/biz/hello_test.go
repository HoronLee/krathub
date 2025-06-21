package biz

import (
	"context"
	"krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// HelloRepo 业务用例
type HelloRepo interface {
	SayHello(ctx context.Context, in string) (string, error)
}

// HelloUsecase 业务依赖
type HelloUsecase struct {
	repo HelloRepo
	log  *log.Helper
	cfg  *conf.App
}

// NewHelloUsecase 创建一个新的 HelloUsecase 实例
func NewHelloUsecase(repo HelloRepo, logger log.Logger, cfg *conf.App) *HelloUsecase {
	uc := &HelloUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
		cfg:  cfg,
	}
	return uc
}

// SayHello 实现业务逻辑的 SayHello 方法
func (uc *HelloUsecase) SayHello(ctx context.Context, in *string) (string, error) {
	uc.log.Debugf("Saying hello with greeting: %s", *in)
	// 调用数据层的 SayHello 方法
	response, err := uc.repo.SayHello(ctx, *in)
	if err != nil {
		uc.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return response, nil
}
