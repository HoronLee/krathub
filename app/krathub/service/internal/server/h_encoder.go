package server

import (
	"encoding/json"
	"net/http"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

const (
	// successCode 是统一成功响应的业务码（与 HTTP 状态码分离）。
	successCode = 0
	// successMessage 是统一成功响应提示。
	successMessage = "OK"
	// defaultInternalReason 是非 Kratos 错误时的默认 reason。
	defaultInternalReason = "INTERNAL_ERROR"

	// 下面这些 Header 名称用于尽量从请求中提取链路/请求标识，方便排障。
	headerXRequestID      = "X-Request-Id"
	headerXRequestIDUpper = "X-Request-ID"
	headerTraceParent     = "TraceParent"
	headerXTraceID        = "X-Trace-Id"
)

// HTTPResponse 是统一 HTTP 出口结构体。
//
// 约定：
// 1. 成功时：code=0, message=OK, data=业务数据；
// 2. 失败时：code=HTTP状态码, reason=业务原因, message=错误信息；
// 3. trace_id 尽可能从请求头透传，便于定位问题。
type HTTPResponse struct {
	Code    int         `json:"code"`
	Reason  string      `json:"reason,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// EncodeResponse 是 Kratos HTTP 的成功响应编码器。
// 框架在 handler 返回 (resp, nil) 时会调用它。
func EncodeResponse(w http.ResponseWriter, r *http.Request, v interface{}) error {
	traceID := getTraceID(r)
	resp := HTTPResponse{
		Code:    successCode,
		Message: successMessage,
		Data:    v,
		TraceID: traceID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// EncodeError 是 Kratos HTTP 的错误响应编码器。
// 框架在 handler 返回 error 时会调用它。
func EncodeError(w http.ResponseWriter, r *http.Request, err error) {
	traceID := getTraceID(r)

	// 1) 优先把 error 转换为 Kratos 标准错误，保留 code/reason/message。
	se := kerrors.FromError(err)
	// 2) 如果是普通 error（不是 Kratos error），统一映射到 500。
	if se == nil {
		se = kerrors.New(http.StatusInternalServerError, defaultInternalReason, err.Error())
	}
	// 3) 兜底：有些错误可能没有 reason，这里补默认值，避免前端判空。
	if se.Reason == "" {
		se.Reason = defaultInternalReason
	}

	resp := HTTPResponse{
		Code:    int(se.Code),
		Reason:  se.Reason,
		Message: se.Message,
		TraceID: traceID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(int(se.Code))
	_ = json.NewEncoder(w).Encode(resp)
}

// getTraceID 按优先级提取请求追踪 ID。
// 这里不依赖特定链路系统，尽量兼容常见 Header。
func getTraceID(r *http.Request) string {
	if r == nil {
		return ""
	}
	if value := r.Header.Get(headerXRequestID); value != "" {
		return value
	}
	if value := r.Header.Get(headerXRequestIDUpper); value != "" {
		return value
	}
	if value := r.Header.Get(headerXTraceID); value != "" {
		return value
	}
	if value := r.Header.Get(headerTraceParent); value != "" {
		return value
	}
	return ""
}

var (
	// 编译期接口断言：确保方法签名满足 Kratos 需要的函数类型。
	_ khttp.EncodeResponseFunc = EncodeResponse
	_ khttp.EncodeErrorFunc    = EncodeError
)
