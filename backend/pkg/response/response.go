package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 统一业务错误码定义（规范语义，避免直接暴露 HTTP 状态码）
const (
	CodeSuccess      = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeForbidden    = 40300
	CodeNotFound     = 40400
	CodeInternal     = 50000
	// 密码相关
	CodePasswordIncorrect    = 50001
	CodePasswordIncorrectMsg = 50002
)

// Response 优化后的响应结构
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`            // 响应时间戳，便于排查
	RequestID string      `json:"request_id,omitempty"` // 链路追踪ID（可选，需配合中间件）
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ValidationError 字段校验错误详情
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ========== 核心响应方法 ==========

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Fail 通用失败响应（建议内部使用）
func Fail(c *gin.Context, httpStatus int, bizCode int, msg string, data ...interface{}) {
	resp := Response{
		Code:      bizCode,
		Message:   msg,
		Timestamp: time.Now().Unix(),
		RequestID: getRequestID(c),
	}
	// 支持携带额外数据（如校验失败的具体字段）
	if len(data) > 0 {
		resp.Data = data[0]
	}
	c.JSON(httpStatus, resp)
	// 关键优化：终止后续中间件执行
	c.Abort()
}

// ========== 语义化错误快捷方法 ==========

// BadRequest 参数错误（400）
func BadRequest(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, CodeBadRequest, msg)
}

// Unauthorized 未授权/未登录（401）
func Unauthorized(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, CodeUnauthorized, msg)
}

// Forbidden 权限不足（403）
func Forbidden(c *gin.Context, msg string) {
	Fail(c, http.StatusForbidden, CodeForbidden, msg)
}

// NotFound 资源不存在（404）
func NotFound(c *gin.Context, msg string) {
	Fail(c, http.StatusNotFound, CodeNotFound, msg)
}

// InternalError 服务器内部错误（500）
func InternalError(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, CodeInternal, msg)
}

// InvalidParams 表单校验失败专用（返回详细字段错误）
func InvalidParams(c *gin.Context, errors []ValidationError) {
	Fail(c, http.StatusBadRequest, CodeBadRequest, "参数校验失败", gin.H{"errors": errors})
}

// getRequestID 尝试从 gin.Context 获取请求 ID（需要在中间件中设置）
func getRequestID(c *gin.Context) string {
	if id, exists := c.Get("RequestID"); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

// ========== 与 apperrors 包的桥接（修复版） ==========

// AppError 处理 apperrors.AppError 类型的错误
func AppError(c *gin.Context, err error) {
	// 类型断言获取 AppError
	type appError interface {
		HTTPStatus() int
		GetCode() int
		GetMessage() string
		GetDetail() interface{}
	}

	if appErr, ok := err.(appError); ok {
		httpStatus := appErr.HTTPStatus()
		bizCode := appErr.GetCode()
		msg := appErr.GetMessage() // 直接获取 Message，不调用 Error()
		detail := appErr.GetDetail()

		if detail == nil {
			Fail(c, httpStatus, bizCode, msg)
		} else {
			Fail(c, httpStatus, bizCode, msg, detail)
		}
		return
	}

	// 如果不是 AppError，当作普通错误处理
	InternalError(c, err.Error())
}

// Error 通用错误处理（根据 error 类型自动分流）
func Error(c *gin.Context, err error) {
	if err == nil {
		return
	}
	AppError(c, err)
}
