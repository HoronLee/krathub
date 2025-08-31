package biz

import (
	"context"

	"github.com/horonlee/krathub/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// SayHelloRepo 是biz层对数据层的接口定义
type SayHelloRepo interface {
	Hello(ctx context.Context, greeting string) (string, error)
}

// SayHelloUsecase is a SayHello usecase.
type SayHelloUsecase struct {
	repo SayHelloRepo
	log  *log.Helper
	cfg  *conf.App
}

// NewSayHelloUsecase 实例化 SayHelloUsecase
func NewSayHelloUsecase(repo SayHelloRepo, logger log.Logger, cfg *conf.App) *SayHelloUsecase {
	return &SayHelloUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
		cfg:  cfg,
	}
}

// Hello 业务逻辑示例
func (uc *SayHelloUsecase) Hello(ctx context.Context, in *string) (string, error) {
	uc.log.Debugf("[biz] Saying hello with greeting: %s", *in)
	response, err := uc.repo.Hello(ctx, *in)
	if err != nil {
		uc.log.Errorf("Failed to say hello: %v", err)
		return "", err
	}
	return response, nil
}
