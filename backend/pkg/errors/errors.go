package apperrors

import (
	"fmt"
	"net/http"
)

// ========== 错误码常量 ==========
// 规则：模块前缀 + 4位序号
// 通用(1xxxx) 用户(2xxxx) 内容(3xxxx) 权限(4xxxx) 积分(5xxxx) 公告(6xxxx) 统计(7xxxx) 时间线(8xxxx) 文件(9xxxx)
const (
	// 通用 (10000-10099)
	CodeUnknown             = 10000
	CodeValidation          = 10001
	CodeUnauthorized        = 10002
	CodeForbidden           = 10003
	CodeNotFound            = 10004
	CodeTooManyRequests     = 10005
	CodeInternalError       = 10006
	CodeInvalidRequest      = 10007
	CodeSystemBusy          = 10008
	CodeInvalidConfirmation = 10009

	// 用户模块 (20000-20999)
	CodeUserNotFound           = 20001
	CodeUserExist              = 20002
	CodeInvalidEmail           = 20010
	CodeInvalidPhone           = 20011
	CodeInvalidPassword        = 20012
	CodeInvalidUsername        = 20013
	CodeInvalidAvatar          = 20014
	CodeInvalidNickname        = 20015
	CodeInvalidUserID          = 20016
	CodeInvalidCurrentPassword = 20017
	CodeInvalidRole            = 20020
	CodeCannotModifySelf       = 20021
	CodeCannotChangeOwner      = 20022
	CodeFollowSelf             = 20040
	CodeAlreadyFollow          = 20041
	CodeNotFollow              = 20042
	CodeScoreNotEnough         = 20060
	CodeUserBlocked            = 20080
	CodeUserDeleted            = 20081

	// 内容模块 (30000-30999)
	CodePostNotFound         = 30001
	CodePostLocked           = 30002
	CodePostPinned           = 30003
	CodePostDeleted          = 30004
	CodeBoardNotFound        = 30020
	CodeQuestionNotFound     = 30040
	CodeAnswerNotFound       = 30041
	CodeCommentNotFound      = 30060
	CodeCommentDeleted       = 30061
	CodeTopicNotFound        = 30080
	CodeTagNotFound          = 30081
	CodeNotificationNotFound = 30082
	CodeLikeAlready          = 30100
	CodeLikeNotExist         = 30101
	CodeCollectAlready       = 30102
	CodeCollectNotExist      = 30103
	CodePasswordTooShort     = 30120
	CodePasswordSameAsOld    = 30121

	// 权限模块 (40000-40999)
	CodeInsufficientPermission = 40001
	CodeAcceptForbidden        = 40002
	CodeModeratorApplyExist    = 40003
	CodeModeratorApplyNotFound = 40004
	CodeAlreadyModerator       = 40005

	// 积分模块 (50000-50999)
	CodeFailedToQueryScore  = 50001
	CodeScoreRecordNotFound = 50002

	// 公告模块 (60000-60999)
	CodeAnnouncementNotFound    = 60001
	CodeAnnouncementInvalidTime = 60002

	// 统计模块 (70000-70999)
	CodeStatsNotFound = 70001

	// 时间线模块 (80000-80999)
	CodeTimelineEmpty = 80001

	// 文件上传模块 (90000-90999)
	CodeFileTooLarge    = 90001
	CodeFileTypeInvalid = 90002
	CodeUploadFailed    = 90003
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
	case CodeUnauthorized:
		return http.StatusUnauthorized

	case CodeForbidden, CodeInsufficientPermission, CodeAcceptForbidden,
		CodeCannotModifySelf, CodeCannotChangeOwner,
		CodeAlreadyFollow, CodeAlreadyModerator:
		return http.StatusForbidden

	case CodeUserNotFound, CodePostNotFound, CodeBoardNotFound,
		CodeQuestionNotFound, CodeAnswerNotFound, CodeCommentNotFound,
		CodeTopicNotFound, CodeTagNotFound, CodeNotificationNotFound,
		CodeAnnouncementNotFound, CodeStatsNotFound,
		CodeModeratorApplyNotFound, CodeScoreRecordNotFound,
		CodeLikeNotExist, CodeCollectNotExist, CodeNotFound:
		return http.StatusNotFound

	case CodeTooManyRequests:
		return http.StatusTooManyRequests

	case CodeValidation, CodeInvalidRequest,
		CodeInvalidEmail, CodeInvalidPhone, CodeInvalidPassword,
		CodeInvalidUsername, CodeInvalidAvatar, CodeInvalidNickname,
		CodeInvalidUserID, CodeInvalidRole, CodeInvalidConfirmation,
		CodeUserExist, CodeFollowSelf, CodeScoreNotEnough,
		CodeLikeAlready, CodeCollectAlready,
		CodePostLocked, CodePostDeleted, CodeCommentDeleted,
		CodePasswordTooShort, CodePasswordSameAsOld,
		CodeFileTooLarge, CodeFileTypeInvalid:
		return http.StatusBadRequest

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

// ========== 预定义错误实例 ==========
// 注意：这些是"模板"，业务层若需要附加 Detail 或 Cause，
// 请使用链式方法（如 ErrUserNotFound.WithDetail(...).WithCause(err)），
// 不要直接修改这些全局变量。

var (
	// 通用
	ErrUnknown             = New(CodeUnknown, "未知错误")
	ErrValidation          = New(CodeValidation, "参数验证失败")
	ErrUnauthorized        = New(CodeUnauthorized, "未授权，请先登录")
	ErrForbidden           = New(CodeForbidden, "权限不足")
	ErrNotFound            = New(CodeNotFound, "资源不存在")
	ErrTooManyRequests     = New(CodeTooManyRequests, "请求过于频繁，请稍后再试")
	ErrInternalError       = New(CodeInternalError, "服务器内部错误")
	ErrInvalidRequest      = New(CodeInvalidRequest, "无效的请求")
	ErrSystemBusy          = New(CodeSystemBusy, "系统繁忙，请稍后再试")
	ErrInvalidConfirmation = New(CodeInvalidConfirmation, "无效的确认信息")

	// 用户模块
	ErrUserNotFound           = New(CodeUserNotFound, "用户不存在") // 用户不存在
	ErrUserExist              = New(CodeUserExist, "用户已存在")
	ErrInvalidEmail           = New(CodeInvalidEmail, "无效的邮箱地址")
	ErrInvalidPhone           = New(CodeInvalidPhone, "无效的手机号码")
	ErrInvalidPassword        = New(CodeInvalidPassword, "无效的密码")
	ErrInvalidCurrentPassword = New(CodeInvalidCurrentPassword, "当前密码不正确")
	ErrInvalidUsername        = New(CodeInvalidUsername, "无效的用户名")
	ErrInvalidAvatar          = New(CodeInvalidAvatar, "无效的头像链接")
	ErrInvalidNickname        = New(CodeInvalidNickname, "无效的昵称")
	ErrInvalidUserID          = New(CodeInvalidUserID, "无效的用户ID")
	ErrInvalidRole            = New(CodeInvalidRole, "无法更改到此角色类型")
	ErrCannotModifySelf       = New(CodeCannotModifySelf, "不能修改自己的信息")
	ErrCannotChangeOwnerRole  = New(CodeCannotChangeOwner, "不能修改超级管理员的角色")
	ErrFollowSelf             = New(CodeFollowSelf, "不能关注自己")
	ErrAlreadyFollow          = New(CodeAlreadyFollow, "已经关注了该用户")
	ErrNotFollow              = New(CodeNotFollow, "尚未关注该用户")
	ErrScoreNotEnough         = New(CodeScoreNotEnough, "积分不足")
	ErrUserBlocked            = New(CodeUserBlocked, "用户已被封禁")
	ErrUserDeleted            = New(CodeUserDeleted, "用户已被删除")
	ErrCannotBlockSelf        = New(CodeCannotModifySelf, "不能封禁自己的账号")
	ErrCannotModifySuperAdmin = New(CodeCannotChangeOwner, "不能修改超级管理员") //
	ErrCannotBlockAdmin       = New(CodeInsufficientPermission, "只有超级管理员才能封禁其他管理员")

	// 密码校验
	ErrPasswordNotMatch  = New(CodeUnauthorized, "密码不匹配")
	ErrPasswordTooShort  = New(CodePasswordTooShort, "密码长度至少为6位")
	ErrPasswordTooLong   = New(CodeInvalidPassword, "密码长度不能超过32位")
	ErrPasswordSameAsOld = New(CodePasswordSameAsOld, "新密码不能与旧密码相同")
	ErrWeakPassword      = New(CodeInvalidPassword, "密码强度太弱，请使用更长且更复杂的密码")

	// Token / 验证码
	ErrTokenExpired    = New(CodeUnauthorized, "Token已过期")
	ErrInvalidToken    = New(CodeUnauthorized, "无效的Token")
	ErrRequiredToken   = New(CodeValidation, "需要Token")
	ErrRequiredCaptcha = New(CodeValidation, "需要验证码")
	ErrInvalidCaptcha  = New(CodeValidation, "验证码错误")

	// 内容模块 - 帖子
	ErrPostNotFound = New(CodePostNotFound, "帖子不存在")
	ErrPostLocked   = New(CodePostLocked, "帖子已被锁定，无法操作")
	ErrPostPinned   = New(CodePostPinned, "帖子置顶状态冲突")
	ErrPostDeleted  = New(CodePostDeleted, "帖子已被删除")

	// 内容模块 - 板块
	ErrBoardNotFound = New(CodeBoardNotFound, "板块不存在")

	// 内容模块 - 问答
	ErrAcceptForbidden  = New(CodeAcceptForbidden, "只有发帖人才能采纳答案")
	ErrAnswerNotFound   = New(CodeAnswerNotFound, "回答不存在")
	ErrQuestionNotFound = New(CodeQuestionNotFound, "问题不存在")

	// 内容模块 - 评论
	ErrCommentNotFound = New(CodeCommentNotFound, "评论不存在")
	ErrCommentDeleted  = New(CodeCommentDeleted, "评论已被删除")

	// 内容模块 - 主题 / 标签 / 通知
	ErrTopicNotFound        = New(CodeTopicNotFound, "主题不存在")
	ErrTagNotFound          = New(CodeTagNotFound, "标签不存在")
	ErrNotificationNotFound = New(CodeNotificationNotFound, "通知不存在")

	// 点赞 / 收藏
	ErrLikeAlready     = New(CodeLikeAlready, "已经点过赞了")
	ErrLikeNotExist    = New(CodeLikeNotExist, "尚未点赞")
	ErrCollectAlready  = New(CodeCollectAlready, "已经收藏过了")
	ErrCollectNotExist = New(CodeCollectNotExist, "尚未收藏")

	// 权限模块 - 版主申请
	ErrModeratorApplyExist    = New(CodeModeratorApplyExist, "已经提交过版主申请，请勿重复提交")
	ErrModeratorApplyNotFound = New(CodeModeratorApplyNotFound, "版主申请不存在")
	ErrAlreadyModerator       = New(CodeAlreadyModerator, "已经是版主，无需重复申请")
	ErrInsufficientPermission = New(CodeInsufficientPermission, "权限不足") // 权限不足

	// 积分模块
	ErrFailedToQueryScore  = New(CodeFailedToQueryScore, "查询积分失败")
	ErrScoreRecordNotFound = New(CodeScoreRecordNotFound, "积分记录不存在")

	// 公告模块
	ErrAnnouncementNotFound = New(CodeAnnouncementNotFound, "公告不存在")
	ErrInvalidPublishTime   = New(CodeAnnouncementInvalidTime, "发布时间无效")
	ErrExpiredTimeInvalid   = New(CodeAnnouncementInvalidTime, "过期时间必须晚于发布时间")

	// 统计模块
	ErrStatsNotFound = New(CodeStatsNotFound, "统计数据不存在")

	// 时间线模块
	ErrTimelineEmpty = New(CodeTimelineEmpty, "时间线暂无内容")

	// 文件上传模块
	ErrFileTooLarge    = New(CodeFileTooLarge, "文件过大")
	ErrFileTypeInvalid = New(CodeFileTypeInvalid, "不支持的文件类型")
	ErrUploadFailed    = New(CodeUploadFailed, "文件上传失败")
)
