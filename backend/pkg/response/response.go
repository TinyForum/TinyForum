package response

import (
	"context"
	"errors"
	"net/http"
	"time"
	"tiny-forum/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 统一业务错误码定义
const (
	CodeSuccess      = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeForbidden    = 40300
	CodeNotFound     = 40400
	CodeConflict     = 40900
	CodeTooManyReq   = 42900
	CodeInternal     = 50000

	// 用户相关错误码 (10000-19999)
	CodeUserNotFound      = 10001
	CodeUserExists        = 10002
	CodePasswordIncorrect = 10003

	// 参数校验相关 (20000-29999)
	CodeValidationFailed = 20001
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
	TraceID   string      `json:"trace_id,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	HasMore  bool        `json:"has_more"`
}

// ValidationError 字段校验错误详情
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// Option 响应选项函数
type Option func(*Response)

// WithTraceID 设置追踪ID
func WithTraceID(traceID string) Option {
	return func(r *Response) {
		r.TraceID = traceID
	}
}

// WithMessage 自定义消息
func WithMessage(msg string) Option {
	return func(r *Response) {
		r.Message = msg
	}
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, opts ...Option) {
	resp := Response{
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
		TraceID:   getTraceID(c),
	}

	for _, opt := range opts {
		opt(&resp)
	}

	c.JSON(http.StatusOK, resp)
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	Success(c, data, WithMessage(msg))
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	hasMore := int64(page*pageSize) < total
	Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  hasMore,
	})
}

// Created 创建资源成功响应 (201)
func Created(c *gin.Context, data interface{}, location string) {
	c.Header("Location", location)
	c.JSON(http.StatusCreated, Response{
		Code:      CodeSuccess,
		Message:   "created successfully",
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
		TraceID:   getTraceID(c),
	})
}

// NoContent 无内容响应 (204)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Fail 失败响应（内部使用）
func Fail(c *gin.Context, httpStatus int, bizCode int, msg string, opts ...Option) {
	resp := Response{
		Code:      bizCode,
		Message:   msg,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
		TraceID:   getTraceID(c),
	}

	for _, opt := range opts {
		opt(&resp)
	}

	c.JSON(httpStatus, resp)
	c.Abort() // 终止后续处理
}

// ========== 语义化错误响应 ==========

// BadRequest 参数错误 (400)
func BadRequest(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, CodeBadRequest, msg)
}

// Unauthorized 未授权 (401)
func Unauthorized(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, CodeUnauthorized, msg)
}

// Forbidden 权限不足 (403)
func Forbidden(c *gin.Context, msg string) {
	Fail(c, http.StatusForbidden, CodeForbidden, msg)
}

// NotFound 资源不存在 (404)
func NotFound(c *gin.Context, msg string) {
	Fail(c, http.StatusNotFound, CodeNotFound, msg)
}

// Conflict 资源冲突 (409)
func Conflict(c *gin.Context, msg string) {
	Fail(c, http.StatusConflict, CodeConflict, msg)
}

// TooManyRequests 请求过于频繁 (429)
func TooManyRequests(c *gin.Context, msg string) {
	Fail(c, http.StatusTooManyRequests, CodeTooManyReq, msg)
}

// InternalError 内部错误 (500)
func InternalError(c *gin.Context, msg string) {
	if msg == "" {
		msg = "系统繁忙，请稍后再试"
	}
	Fail(c, http.StatusInternalServerError, CodeInternal, msg)
}

// ValidationFailed 参数校验失败
func ValidationFailed(c *gin.Context, errors []ValidationError) {
	Fail(c, http.StatusBadRequest, CodeValidationFailed, "参数校验失败", WithMessage("validation failed"))
	// 通过 Data 字段传递详细错误
	c.Set("validation_errors", errors)
}

// ========== 增强的错误处理 ==========

// AppError 应用错误接口
type AppError interface {
	error
	HTTPStatus() int
	GetCode() int
	GetMessage() string
	GetDetail() interface{}
}

// HandleError 统一错误处理入口（推荐使用）
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 1. 处理应用自定义错误
	if appErr, ok := err.(AppError); ok {
		handleAppError(c, appErr)
		return
	}

	// 2. 处理 validator 校验错误
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		handleValidationErrors(c, validationErrors)
		return
	}

	// 3. 处理标准库错误类型
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		Fail(c, http.StatusGatewayTimeout, 50400, "请求超时")
		return
	case errors.Is(err, context.Canceled):
		// 客户端主动取消，不记录错误日志
		c.Status(http.StatusNoContent)
		c.Abort()
		return
	}

	// 4. 兜底处理：记录错误但返回通用信息
	logger.Error("unhandled error occurred",
		zap.Error(err),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()),
	)

	InternalError(c, "") // 使用默认消息
}

// handleAppError 处理应用错误
func handleAppError(c *gin.Context, err AppError) {
	httpStatus := err.HTTPStatus()
	bizCode := err.GetCode()
	msg := err.GetMessage()
	detail := err.GetDetail()

	// 根据 HTTP 状态码决定日志级别
	switch {
	case httpStatus >= 500:
		logger.Error("app server error",
			zap.Error(err),
			zap.Int("http_status", httpStatus),
			zap.Int("biz_code", bizCode),
		)
	case httpStatus >= 400:
		logger.Warn("app client error",
			zap.Error(err),
			zap.Int("http_status", httpStatus),
			zap.Int("biz_code", bizCode),
			zap.String("path", c.Request.URL.Path),
		)
	}

	if detail == nil {
		Fail(c, httpStatus, bizCode, msg)
	} else {
		Fail(c, httpStatus, bizCode, msg, WithMessage(msg))
		c.Set("error_detail", detail)
	}
}

// handleValidationErrors 处理 validator 校验错误
func handleValidationErrors(c *gin.Context, err validator.ValidationErrors) {
	errors := make([]ValidationError, 0, len(err))
	for _, fe := range err {
		errors = append(errors, ValidationError{
			Field:   fe.Field(),
			Message: getValidationMessage(fe),
			Value:   fe.Value(),
		})
	}
	ValidationFailed(c, errors)

	logger.Warn("validation failed",
		zap.String("path", c.Request.URL.Path),
		zap.Any("errors", errors),
	)
}

// ========== 辅助函数 ==========

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if id := c.GetString("RequestID"); id != "" {
		return id
	}
	if id := c.GetHeader("X-Request-ID"); id != "" {
		return id
	}
	return ""
}

// getTraceID 获取链路追踪ID
func getTraceID(c *gin.Context) string {
	if traceID := c.GetString("TraceID"); traceID != "" {
		return traceID
	}
	if traceID := c.GetHeader("X-Trace-ID"); traceID != "" {
		return traceID
	}
	// 如果没有 traceID，可以使用 requestID 作为 fallback
	return getRequestID(c)
}

// ========== 中间件 ==========

// RecoveryMiddleware 恢复中间件，防止 panic 导致服务崩溃
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Stack("stack"),
				)
				InternalError(c, "系统内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有未处理的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			// 如果响应已经写入，不再处理
			if c.Writer.Written() {
				logger.Warn("response already written, skip error handling", zap.Error(err))
				return
			}
			HandleError(c, err)
		}
	}
}
