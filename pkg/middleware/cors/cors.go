package cors

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Options 包含 CORS 中间件的配置选项
type Options struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultOptions 返回默认的 CORS 配置
func DefaultOptions() Options {
	return Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}
}

// Middleware 创建 CORS 中间件
func Middleware(options Options) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			// 设置响应头
			setCORSHeaders(w, options, origin)
			// 处理预检请求
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// setCORSHeaders 设置 CORS 响应头
func setCORSHeaders(w http.ResponseWriter, options Options, origin string) {
	// 检查源是否被允许
	allowedOrigin := ""
	if isOriginAllowed(origin, options.AllowedOrigins) {
		allowedOrigin = origin
	}
	if allowedOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	}
	// 设置允许的方法
	if len(options.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ", "))
	}
	// 设置允许的头部
	if len(options.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
	}
	// 设置暴露的头部
	if len(options.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(options.ExposedHeaders, ", "))
	}
	// 设置是否允许凭证
	if options.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	// 设置预检请求缓存时间
	if options.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", int64(options.MaxAge.Seconds())))
	}
}

// isOriginAllowed 检查源是否被允许
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
		// 支持通配符匹配，如 *.example.com
		if strings.HasPrefix(allowedOrigin, "*.") {
			suffix := strings.TrimPrefix(allowedOrigin, "*.")
			if strings.HasSuffix(origin, suffix) {
				// 检查是否只有一个点，防止匹配过长
				parts := strings.Split(strings.TrimSuffix(origin, suffix), ".")
				if len(parts) == 2 {
					return true
				}
			}
		}
	}

	return false
}
