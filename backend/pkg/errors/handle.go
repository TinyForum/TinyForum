package apperrors

import (
	"fmt"
	"net/http"
)

// ========== 核心错误类型 ==========

// AppError 是框架统一的结构化错误类型，同时实现了 error 接口。
// 业务层只需返回 *AppError，response 层通过类型断言识别并处理。
type AppError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"` // 可附加给客户端的补充信息
	Err     error       `json:"-"`                // 内部原始错误，不对外暴露
}

// Error 实现 error 接口，格式为 "[code] message: cause"
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is / errors.As 链式解包
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is 按错误码比较，使 errors.Is(err, ErrUserNotFound) 成立
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// HTTPStatus 根据业务错误码映射对应的 HTTP 状态码
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	// 业务错误
	case CodeUnauthorized:
		return http.StatusUnauthorized

	// 权限错误
	case CodeForbidden, CodeInsufficientPermission, CodeAcceptForbidden,
		CodeCannotModifySelf, CodeCannotChangeOwner,
		CodeAlreadyFollow, CodeAlreadyModerator:
		return http.StatusForbidden

		// 资源不存在
	case CodeUserNotFound, CodePostNotFound, CodeBoardNotFound,
		CodeQuestionNotFound, CodeAnswerNotFound, CodeCommentNotFound,
		CodeTopicNotFound, CodeTagNotFound, CodeNotificationNotFound,
		CodeAnnouncementNotFound, CodeStatsNotFound,
		CodeModeratorApplyNotFound, CodeScoreRecordNotFound,
		CodeLikeNotExist, CodeCollectNotExist, CodeNotFound:
		return http.StatusNotFound

		// 请求冲突
	case CodeTooManyRequests:
		return http.StatusTooManyRequests

		// 请求错误
	case CodeValidation, // 通用验证错误
		CodeInvalidRequest,             // 无效请求
		CodeUserEmailOrPasswordInvalid, // 用户名或密码错误
		CodeInvalidEmail,               // 无效邮箱
		CodeInvalidPhone,               // 无效手机号
		CodeInvalidPassword,
		CodeInvalidUsername, CodeInvalidAvatar, CodeInvalidNickname,
		CodeInvalidUserID, CodeInvalidRole, CodeInvalidConfirmation,
		CodeUserExist, CodeFollowSelf, CodeScoreNotEnough,
		CodeLikeAlready, CodeCollectAlready,
		CodePostLocked, CodePostDeleted, CodeCommentDeleted,
		CodePasswordTooShort, CodePasswordSameAsOld,
		CodeFileTooLarge, CodeFileTypeInvalid:
		return http.StatusBadRequest
		//
	default:
		return http.StatusInternalServerError
	}
}

// ========== 构造函数 ==========

// New 创建一个新的 AppError（适合自定义错误，不复用预定义实例）
func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Newf 格式化 message 创建 AppError
func Newf(code int, format string, args ...interface{}) *AppError {
	return &AppError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// ========== 链式方法（不修改原实例，返回新副本）==========

// WithDetail 附加补充信息（将发送给客户端），返回新实例
func (e *AppError) WithDetail(detail interface{}) *AppError {
	return &AppError{Code: e.Code, Message: e.Message, Detail: detail, Err: e.Err}
}

// WithCause 包装底层错误（仅用于日志，不对外暴露），返回新实例
func (e *AppError) WithCause(cause error) *AppError {
	return &AppError{Code: e.Code, Message: e.Message, Detail: e.Detail, Err: cause}
}

// WithMessage 覆盖错误信息，返回新实例（用于为通用错误添加上下文）
func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{Code: e.Code, Message: msg, Detail: e.Detail, Err: e.Err}
}

// WithMessagef 格式化覆盖错误信息，返回新实例
func (e *AppError) WithMessagef(format string, args ...interface{}) *AppError {
	return e.WithMessage(fmt.Sprintf(format, args...))
}

// ========== 独立包装函数（适合对现有 error 进行包装）==========

// Wrap 将任意 error 包装为 AppError，追加上下文说明
func Wrap(base *AppError, cause error) *AppError {
	return &AppError{Code: base.Code, Message: base.Message, Detail: base.Detail, Err: cause}
}

// Wrapf 将任意 error 包装为 AppError，并格式化补充说明到 Message
func Wrapf(base *AppError, cause error, format string, args ...interface{}) *AppError {
	msg := fmt.Sprintf("%s: %s", base.Message, fmt.Sprintf(format, args...))
	return &AppError{Code: base.Code, Message: msg, Detail: base.Detail, Err: cause}
}
