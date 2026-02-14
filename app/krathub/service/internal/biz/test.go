package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pkglogger "github.com/horonlee/krathub/pkg/logger"
)

type TestRepo interface {
	Hello(ctx context.Context, in string) (string, error)
}

type TestUsecase struct {
	repo TestRepo
	log  *log.Helper
}

func NewTestUsecase(repo TestRepo, logger log.Logger) *TestUsecase {
	return &TestUsecase{
		repo: repo,
		log:  log.NewHelper(pkglogger.With(logger, pkglogger.WithModule("test/biz/krathub-service"))),
	}
}

func (uc *TestUsecase) Hello(ctx context.Context, in string) (string, error) {
	greeting := "World"
	if in != "" {
		greeting = in
	}
	uc.log.Debugf("Saying hello with greeting: %s", greeting)
	return uc.repo.Hello(ctx, greeting)
}
