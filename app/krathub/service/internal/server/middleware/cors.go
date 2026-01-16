package middleware

import (
	"github.com/horonlee/krathub/api/gen/go/conf/v1"
	"github.com/horonlee/krathub/pkg/middleware/cors"
)

// CORS 从配置文件创建 CORS 选项
func CORS(corsConfig *conf.CORS) cors.Options {
	if corsConfig == nil || !corsConfig.GetEnable() {
		return cors.Options{} // 返回空配置表示禁用 CORS
	}
	options := cors.DefaultOptions()
	if len(corsConfig.GetAllowedOrigins()) > 0 {
		options.AllowedOrigins = corsConfig.GetAllowedOrigins()
	}
	if len(corsConfig.GetAllowedMethods()) > 0 {
		options.AllowedMethods = corsConfig.GetAllowedMethods()
	}
	if len(corsConfig.GetAllowedHeaders()) > 0 {
		options.AllowedHeaders = corsConfig.GetAllowedHeaders()
	}
	if len(corsConfig.GetExposedHeaders()) > 0 {
		options.ExposedHeaders = corsConfig.GetExposedHeaders()
	}
	// Since AllowCredentials is a bool (not *bool), we use the value directly
	options.AllowCredentials = corsConfig.GetAllowCredentials()
	if corsConfig.MaxAge != nil {
		options.MaxAge = corsConfig.MaxAge.AsDuration()
	}
	return options
}
