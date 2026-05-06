package response

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"tiny-forum/internal/model/common"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// ========== 响应结构体 ==========

// ValidationError 字段校验错误详情（发送给客户端）
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// ========== 响应选项 ==========

// Option 响应选项函数，用于在构造响应时附加额外字段
type Option func(*common.BasicResponse)

// WithTraceID 设置追踪ID
func WithTraceID(traceID string) Option {
	return func(r *common.BasicResponse) { r.TraceID = traceID }
}

// WithMessage 覆盖默认消息
func WithMessage(msg string) Option {
	return func(r *common.BasicResponse) { r.Message = msg }
}

// ========== 成功响应 ==========

// Success 返回成功响应 (HTTP 200)
// Success 函数用于处理成功的响应
// 它接收一个 gin.Context 上下文、任意类型的数据和可选的配置选项
// 然后构造一个响应对象并返回 HTTP 200 状态码
func Success(c *gin.Context, data any, opts ...Option) {
	// 创建一个新的响应对象，默认状态码为 0，消息为 "success"
	resp := newResp(c, 0, "success")
	// 设置响应状态码为 0（成功）
	resp.Code = 0
	// 将传入的数据设置到响应对象中
	resp.Data = data
	// 应用传入的可选配置选项
	applyOpts(&resp, opts)
	// 将响应对象以 JSON 格式返回，HTTP 状态码为 200（OK）
	c.JSON(http.StatusOK, resp)
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	Success(c, data, WithMessage(msg))
}

// SuccessPage 分页成功响应
func SuccessPage[T any](c *gin.Context, list []T, total int64, page, pageSize int, opts ...Option) {
	hasMore := int64(page*pageSize) < total
	result := common.PageResult[T]{ // 确保 PageResult 在作用域内（vo.PageResult 或 common.PageResult）
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
		HasMore:  hasMore,
	}
	Success(c, result, opts...)
}

// Created 创建资源成功响应 (HTTP 201)，同时写入 Location 头
func Created(c *gin.Context, data interface{}, location string) {
	if location != "" {
		c.Header("Location", location)
	}
	resp := newResp(c, 0, "created successfully")
	resp.Data = data
	c.JSON(http.StatusCreated, resp)
}

// NoContent 无内容响应 (HTTP 204)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// ========== 错误响应（内部基础实现）==========

// fail 写入错误响应并终止后续 handler（内部调用）
func fail(c *gin.Context, httpStatus int, bizCode int, msg string, data interface{}, opts ...Option) {
	resp := newResp(c, bizCode, msg)
	if data != nil {
		resp.Data = data
	}
	applyOpts(&resp, opts)
	c.JSON(httpStatus, resp)
	c.Abort()
}

// ========== 语义化错误响应（直接调用版）==========

// BadRequest 参数错误 (HTTP 400)
func BadRequest(c *gin.Context, msg string) {
	fail(c, http.StatusBadRequest, apperrors.CodeInvalidRequest, msg, nil)
}

// Unauthorized 未授权 (HTTP 401)
func Unauthorized(c *gin.Context, msg string) {
	fail(c, http.StatusUnauthorized, apperrors.CodeUnauthorized, msg, nil)
}

// Forbidden 权限不足 (HTTP 403)
func Forbidden(c *gin.Context, msg string) {
	fail(c, http.StatusForbidden, apperrors.CodeForbidden, msg, nil)
}

// NotFound 资源不存在 (HTTP 404)
func NotFound(c *gin.Context, msg string) {
	fail(c, http.StatusNotFound, apperrors.CodeNotFound, msg, nil)
}

// Conflict 资源冲突 (HTTP 409)
func Conflict(c *gin.Context, msg string) {
	fail(c, http.StatusConflict, apperrors.CodeInvalidRequest, msg, nil)
}

// TooManyRequests 请求过于频繁 (HTTP 429)
func TooManyRequests(c *gin.Context, msg string) {
	fail(c, http.StatusTooManyRequests, apperrors.CodeTooManyRequests, msg, nil)
}

// InternalError 内部错误 (HTTP 500)
func InternalError(c *gin.Context, msg string) {
	if msg == "" {
		msg = "系统繁忙，请稍后再试"
	}
	fail(c, http.StatusInternalServerError, apperrors.CodeInternalError, msg, nil)
}

func ValidationFailed(c *gin.Context, errs []ValidationError) {
	fail(c, http.StatusBadRequest, apperrors.CodeValidation, "参数校验失败", errs)
}

// ========== 统一错误处理入口 ==========

// HandleError 统一错误处理入口（推荐在 handler 层统一调用）
//
// 处理优先级：
//  1. *apperrors.AppError  —— 业务自定义错误，直接映射 HTTP 状态码
//  2. validator.ValidationErrors —— struct tag 校验失败，展开字段错误
//  3. context.DeadlineExceeded / context.Canceled —— 超时 / 客户端取消
//  4. 其他未知错误 —— 记录日志，返回 500
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 1. 业务错误
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		handleAppError(c, appErr)
		return
	}

	// 2. validator 校验错误
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		handleValidationErrors(c, ve)
		return
	}

	// 3. 标准库 context 错误
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		fail(c, http.StatusGatewayTimeout, apperrors.CodeSystemBusy, "请求超时", nil)
		return
	case errors.Is(err, context.Canceled):
		// 客户端主动取消，不记录错误日志，静默结束
		c.Status(http.StatusNoContent)
		c.Abort()
		return
	}

	// 4. 兜底：记录日志，返回通用 500
	logger.Error("unhandled error",
		zap.Error(err),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("client_ip", c.ClientIP()),
	)
	InternalError(c, "")
}

// ========== 内部处理函数 ==========

// handleAppError 处理 *apperrors.AppError
func handleAppError(c *gin.Context, err *apperrors.AppError) {
	httpStatus := err.HTTPStatus()

	// 按 HTTP 状态级别记录日志
	switch {
	case httpStatus >= 500:
		logger.Error("server error",
			zap.Error(err),
			zap.Int("http_status", httpStatus),
			zap.Int("biz_code", err.Code),
		)
	case httpStatus >= 400:
		logger.Warn("client error",
			zap.Error(err),
			zap.Int("http_status", httpStatus),
			zap.Int("biz_code", err.Code),
			zap.String("path", c.Request.URL.Path),
		)
	}

	// Detail 不为 nil 时一并写入响应体的 data 字段
	fail(c, httpStatus, err.Code, err.Message, err.Detail)
}

// handleValidationErrors 将 validator.ValidationErrors 展开为字段级错误列表
func handleValidationErrors(c *gin.Context, ve validator.ValidationErrors) {
	errs := make([]ValidationError, 0, len(ve))
	for _, fe := range ve {
		errs = append(errs, ValidationError{
			Field:   fe.Field(),
			Message: validationMessage(fe),
			Value:   fe.Value(),
		})
	}

	logger.Warn("validation failed",
		zap.String("path", c.Request.URL.Path),
		zap.Any("errors", errs),
	)

	// 字段错误列表写入 data 字段，方便前端逐字段展示
	fail(c, http.StatusBadRequest, apperrors.CodeValidation, "参数校验失败", errs)
}

// validationMessage 将 validator tag 转为中文提示
func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s 不能为空", fe.Field())
	case "email":
		return fmt.Sprintf("%s 格式不正确", fe.Field())
	case "min":
		return fmt.Sprintf("%s 最小长度为 %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s 最大长度为 %s", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s 长度必须为 %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s 不能小于 %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s 不能大于 %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s 必须是以下值之一: %s", fe.Field(), fe.Param())
	case "url":
		return fmt.Sprintf("%s 不是有效的URL", fe.Field())
	case "numeric":
		return fmt.Sprintf("%s 必须为数字", fe.Field())
	default:
		return fmt.Sprintf("%s 校验失败 (规则: %s)", fe.Field(), fe.Tag())
	}
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

// ErrorHandlerMiddleware 尾部错误处理中间件
// 若 handler 层通过 c.Error(err) 记录错误而未直接写响应，此中间件负责兜底处理。
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// 响应已写入时跳过（避免重复写）
		if c.Writer.Written() {
			logger.Warn("response already written, skip error handling",
				zap.Error(c.Errors.Last().Err),
			)
			return
		}

		HandleError(c, c.Errors.Last().Err)
	}
}

// ========== 内部辅助函数 ==========

// newResp 构造带公共字段的 Response
func newResp(c *gin.Context, code int, msg string) common.BasicResponse {
	return common.BasicResponse{
		Code:      code,
		Message:   msg,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
		TraceID:   getTraceID(c),
	}
}

// applyOpts 函数用于接收一个 BasicResponse 指针和 Option 切片，
// 并遍历 opts 切片中的每个 Option 函数，对传入的 r 进行配置
func applyOpts(r *common.BasicResponse, opts []Option) {
	for _, opt := range opts {
		opt(r)
	}
}

func getRequestID(c *gin.Context) string {
	if id := c.GetString("RequestID"); id != "" {
		return id
	}
	return c.GetHeader("X-Request-ID")
}

func getTraceID(c *gin.Context) string {
	if id := c.GetString("TraceID"); id != "" {
		return id
	}
	if id := c.GetHeader("X-Trace-ID"); id != "" {
		return id
	}
	return getRequestID(c) // 降级用 RequestID
}
