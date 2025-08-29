package middleware

import "github.com/horonlee/krathub/internal/conf"

// MiddlewareManager 中间件管理器
type MiddlewareManager struct {
	appConf *conf.App
}

// NewMiddlewareManager 创建中间件管理器
func NewMiddlewareManager(appConf *conf.App) *MiddlewareManager {
	return &MiddlewareManager{
		appConf: appConf,
	}
}
