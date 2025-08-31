package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/horonlee/krathub/internal/biz"
)

// sayHelloRepo 是 SayHelloRepo 接口的实现
type sayHelloRepo struct {
	data *Data
	log  *log.Helper
}

// NewSayHelloRepo new a SayHello repo.
func NewSayHelloRepo(data *Data, logger log.Logger) biz.SayHelloRepo {
	return &sayHelloRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Hello 模拟了简单的数据库查询
func (r *sayHelloRepo) Hello(ctx context.Context, in string) (string, error) {
	r.log.Debugf("[data] Saying hello with greeting: %s", in)
	return in, nil
}
