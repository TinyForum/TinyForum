package apperrors

import (
	"fmt"
	"net/http"
)

// ========== 错误码常量（与 response 包呼应） ==========
const (
	// 通用错误 100xx
	CodeUnknown      = 10000
	CodeValidation   = 10001
	CodeUnauthorized = 10002
	CodeForbidden    = 10003
	CodeNotFound     = 10004

	// 用户模块 200xx
	CodeUserNotFound      = 20001
	CodeUserExist         = 20002
	CodeInvalidEmail      = 20003
	CodeInvalidPhone      = 20004
	CodeInvalidPassword   = 20005
	CodeInvalidUsername   = 20006
	CodeInvalidAvatar     = 20007
	CodeInvalidNickname   = 20008
	CodeInvalidUserID     = 20009
	CodeInvalidRole       = 20010
	CodeCannotModifySelf  = 20011
	CodeCannotChangeOwner = 20012

	// 内容模块 300xx
	CodePostNotFound     = 30001
	CodeBoardNotFound    = 30002
	CodeQuestionNotFound = 30003
	CodeAnswerNotFound   = 30004

	// 权限模块 400xx
	CodeInsufficientPermission = 40001
	CodeAcceptForbidden        = 40002

	// 积分模块 500xx
	CodeFailedToQueryScore = 50001
	// 密码

)

// ========== 结构化错误类型 ==========

// AppError 应用错误结构
type AppError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
	Err     error       `json:"-"` // 原始错误，不暴露给前端
}

// Unwrap 支持 errors.Is 和 errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// ========== 预定义错误实例（直接替换原有的 var） ==========

var (
	// 通用
	ErrUnknown = &AppError{Code: CodeUnknown, Message: "未知错误"}

	// 内容
	ErrPostNotFound = &AppError{Code: CodePostNotFound, Message: "帖子不存在"}

	// 角色
	ErrInvalidRole           = &AppError{Code: CodeInvalidRole, Message: "无法更改到此角色类型"}
	ErrCannotChangeOwnerRole = &AppError{Code: CodeCannotChangeOwner, Message: "不能修改超级管理员的角色"}
	ErrCannotModifySelf      = &AppError{Code: CodeCannotModifySelf, Message: "不能修改自己的角色"}

	// 权限
	ErrInsufficientPermission = &AppError{Code: CodeInsufficientPermission, Message: "权限不足"}

	// 板块
	ErrBoardNotFound = &AppError{Code: CodeBoardNotFound, Message: "板块不存在"}

	// 问答
	ErrAcceptForbidden  = &AppError{Code: CodeAcceptForbidden, Message: "只有发帖人才能采纳"}
	ErrAnswerNotFound   = &AppError{Code: CodeAnswerNotFound, Message: "回答不存在"}
	ErrQuestionNotFound = &AppError{Code: CodeQuestionNotFound, Message: "问题不存在"}

	// 用户信息
	ErrUserNotFound    = &AppError{Code: CodeUserNotFound, Message: "用户不存在"}
	ErrUserExist       = &AppError{Code: CodeUserExist, Message: "用户已存在"}
	ErrInvalidEmail    = &AppError{Code: CodeInvalidEmail, Message: "无效的邮箱"}
	ErrInvalidPhone    = &AppError{Code: CodeInvalidPhone, Message: "无效的手机号"}
	ErrInvalidPassword = &AppError{Code: CodeInvalidPassword, Message: "无效的密码"}
	ErrInvalidUsername = &AppError{Code: CodeInvalidUsername, Message: "无效的用户名"}
	ErrInvalidAvatar   = &AppError{Code: CodeInvalidAvatar, Message: "无效的头像"}
	ErrInvalidNickname = &AppError{Code: CodeInvalidNickname, Message: "无效的昵称"}
	ErrInvalidUserID   = &AppError{Code: CodeInvalidUserID, Message: "无效的用户ID"}

	// 积分
	ErrFailedToQueryScore = &AppError{Code: CodeFailedToQueryScore, Message: "查询积分失败"}

	// 封禁相关
	ErrCannotBlockSelf       = &AppError{Code: CodeCannotModifySelf, Message: "不能封禁自己的账号"}
	ErrCannotBlockSuperAdmin = &AppError{Code: CodeCannotChangeOwner, Message: "不能封禁超级管理员"}
	ErrCannotBlockAdmin      = &AppError{Code: CodeInsufficientPermission, Message: "只有超级管理员才能封禁其他管理员"}

	// 密码相关
	ErrPasswordNotMatch  = &AppError{Code: CodeUnauthorized, Message: "密码不匹配"}
	ErrPasswordTooShort  = &AppError{Code: CodeInvalidPassword, Message: "密码长度至少为6位"}
	ErrPasswordTooLong   = &AppError{Code: CodeInvalidPassword, Message: "密码长度不能超过32位"}
	ErrPasswordSameAsOld = &AppError{Code: CodeInvalidPassword, Message: "新密码不能与旧密码相同"}
)

// ========== 辅助函数：追加上下文信息 ==========

// Wrap 包装错误，追加上下文信息
// 示例：apperrors.Wrap(apperrors.ErrUserNotFound, "用户ID: 123")
func Wrap(err *AppError, context string) *AppError {
	return &AppError{
		Code:    err.Code,
		Message: fmt.Sprintf("%s (%s)", err.Message, context),
		Err:     err,
	}
}

// Wrapf 格式化上下文
func Wrapf(err *AppError, format string, args ...interface{}) *AppError {
	return Wrap(err, fmt.Sprintf(format, args...))
}

// WithDetail 附加详情数据（如字段校验错误列表）
func (e *AppError) WithDetail(detail interface{}) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Detail:  detail,
		Err:     e.Err,
	}
}

// ========== HTTP 状态码映射（供 response 包调用） ==========

// GetCode 获取错误码（避免字段冲突）
func (e *AppError) GetCode() int {
	return e.Code
}

// GetMessage 获取错误信息
func (e *AppError) GetMessage() string {
	return e.Message
}

// GetDetail 获取详情
func (e *AppError) GetDetail() interface{} {
	return e.Detail
}

// HTTPStatus 根据错误码返回对应的 HTTP 状态码
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden, CodeInsufficientPermission, CodeAcceptForbidden,
		CodeCannotModifySelf, CodeCannotChangeOwner:
		return http.StatusForbidden
	case CodeUserNotFound, CodePostNotFound, CodeBoardNotFound,
		CodeQuestionNotFound, CodeAnswerNotFound:
		return http.StatusNotFound
	case CodeInvalidEmail, CodeInvalidPhone, CodeInvalidPassword,
		CodeInvalidUsername, CodeInvalidAvatar, CodeInvalidNickname,
		CodeInvalidUserID, CodeInvalidRole, CodeUserExist, CodeValidation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Error 实现 error 接口（仅用于日志输出）
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}
