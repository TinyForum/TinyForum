package apperrors

import (
	"fmt"
	"net/http"
)

// ========== 错误码常量 ==========
// 错误码格式：模块码(2位) + 子模块码(2位) + 序号(2位)
// 例如：200101 表示用户模块(20) -> 基础信息(01) -> 用户不存在(01)
// 为便于阅读，这里使用十进制连续编号，但按功能区间预留号码
const (
	// ------------------------------------------------------------
	// 通用错误 (10000-10099)
	// ------------------------------------------------------------
	CodeUnknown         = 10000
	CodeValidation      = 10001
	CodeUnauthorized    = 10002
	CodeForbidden       = 10003
	CodeNotFound        = 10004
	CodeTooManyRequests = 10005
	CodeInternalError   = 10006
	// 10007-10099 预留

	// ------------------------------------------------------------
	// 用户模块 (20000-20999)
	// ------------------------------------------------------------
	// 基础信息 (20001-20019)
	CodeUserNotFound = 20001
	CodeUserExist    = 20002
	// 20003-20009 预留
	CodeInvalidEmail    = 20010
	CodeInvalidPhone    = 20011
	CodeInvalidPassword = 20012
	CodeInvalidUsername = 20013
	CodeInvalidAvatar   = 20014
	CodeInvalidNickname = 20015
	CodeInvalidUserID   = 20016
	// 20017-20019 预留

	// 角色与权限 (20020-20039)
	CodeInvalidRole       = 20020
	CodeCannotModifySelf  = 20021
	CodeCannotChangeOwner = 20022
	// 20023-20039 预留

	// 关注关系 (20040-20059)
	CodeFollowSelf    = 20040
	CodeAlreadyFollow = 20041
	CodeNotFollow     = 20042
	// 20043-20059 预留

	// 积分 (20060-20079)
	CodeScoreNotEnough = 20060
	// 20061-20079 预留

	// 封禁 (20080-20099) 注：复用部分角色/权限码，但单独列出方便扩展
	// 20080-20099 预留（用户封禁、注销等）

	// ------------------------------------------------------------
	// 内容模块 (30000-30999)
	// ------------------------------------------------------------
	// 帖子 (30001-30019)
	CodePostNotFound = 30001
	CodePostLocked   = 30002
	CodePostPinned   = 30003
	CodePostDeleted  = 30004
	// 30005-30019 预留

	// 板块 (30020-30039)
	CodeBoardNotFound = 30020
	// 30021-30039 预留

	// 问答 (30040-30059)
	CodeQuestionNotFound = 30040 // 问题不存在
	CodeAnswerNotFound   = 30041 // 回答不存在
	// 30043-30059 预留

	// 评论 (30060-30079)
	CodeCommentNotFound = 30060
	CodeCommentDeleted  = 30061
	// 30062-30079 预留

	// 主题、标签、通知 (30080-30099)
	CodeTopicNotFound        = 30080
	CodeTagNotFound          = 30081
	CodeNotificationNotFound = 30082
	// 30083-30099 预留

	// 点赞收藏 (30100-30119)
	CodeLikeAlready     = 30100
	CodeLikeNotExist    = 30101
	CodeCollectAlready  = 30102
	CodeCollectNotExist = 30103
	// 30104-30119 预留

	// 密码校验 (30120-30139) —— 虽然属于用户模块，但放在内容模块区？根据原代码这里归在内容模块，保留
	CodePasswordTooShort  = 30120
	CodePasswordSameAsOld = 30121
	// 30122-30139 预留

	// 30140-30199 预留（内容模块后续扩展）

	// ------------------------------------------------------------
	// 权限模块 (40000-40999)
	// ------------------------------------------------------------
	CodeInsufficientPermission = 40001
	CodeAcceptForbidden        = 40002
	CodeModeratorApplyExist    = 40003
	CodeModeratorApplyNotFound = 40004
	CodeAlreadyModerator       = 40005
	// 40006-40099 预留

	// ------------------------------------------------------------
	// 积分模块 (50000-50999) —— 独立于用户模块的积分错误
	// ------------------------------------------------------------
	CodeFailedToQueryScore  = 50001
	CodeScoreRecordNotFound = 50002
	// 50003-50099 预留

	// ------------------------------------------------------------
	// 公告模块 (60000-60999)
	// ------------------------------------------------------------
	CodeAnnouncementNotFound    = 60001
	CodeAnnouncementInvalidTime = 60002
	// 60003-60099 预留

	// ------------------------------------------------------------
	// 统计模块 (70000-70999)
	// ------------------------------------------------------------
	CodeStatsNotFound = 70001 // 统计数据不存在
	// 70002-70099 预留

	// ------------------------------------------------------------
	// 时间线模块 (80000-80999)
	// ------------------------------------------------------------
	CodeTimelineEmpty = 80001
	// 80002-80099 预留

	// ------------------------------------------------------------
	// 文件上传模块 (90000-90999)
	// ------------------------------------------------------------
	CodeFileTooLarge    = 90001
	CodeFileTypeInvalid = 90002
	CodeUploadFailed    = 90003
	// 90004-90099 预留
)

// ========== 结构化错误类型 ==========

type AppError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
	Err     error       `json:"-"`
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// ========== 预定义错误实例 ==========

var (
	// 通用
	ErrUnknown               = &AppError{Code: CodeUnknown, Message: "未知错误"}
	ErrValidation            = &AppError{Code: CodeValidation, Message: "参数验证失败"}
	ErrUnauthorized          = &AppError{Code: CodeUnauthorized, Message: "未授权，请先登录"}
	ErrForbidden             = &AppError{Code: CodeForbidden, Message: "权限不足"}
	ErrNotFound              = &AppError{Code: CodeNotFound, Message: "资源不存在"}
	ErrTooManyRequests       = &AppError{Code: CodeTooManyRequests, Message: "请求过于频繁，请稍后再试"}
	ErrInternalError         = &AppError{Code: CodeInternalError, Message: "服务器内部错误"}
	ErrCodePasswordTooShort  = &AppError{Code: CodePasswordTooShort, Message: "密码长度至少为6位"}
	ErrCodePasswordSameAsOld = &AppError{Code: CodePasswordSameAsOld, Message: "新密码不能与旧密码相同"}
	// 用户模块
	ErrUserNotFound          = &AppError{Code: CodeUserNotFound, Message: "用户不存在"}
	ErrUserExist             = &AppError{Code: CodeUserExist, Message: "用户已存在"}
	ErrInvalidEmail          = &AppError{Code: CodeInvalidEmail, Message: "无效的邮箱地址"}
	ErrInvalidPhone          = &AppError{Code: CodeInvalidPhone, Message: "无效的手机号码"}
	ErrInvalidPassword       = &AppError{Code: CodeInvalidPassword, Message: "无效的密码"}
	ErrInvalidUsername       = &AppError{Code: CodeInvalidUsername, Message: "无效的用户名"}
	ErrInvalidAvatar         = &AppError{Code: CodeInvalidAvatar, Message: "无效的头像链接"}
	ErrInvalidNickname       = &AppError{Code: CodeInvalidNickname, Message: "无效的昵称"}
	ErrInvalidUserID         = &AppError{Code: CodeInvalidUserID, Message: "无效的用户ID"}
	ErrInvalidRole           = &AppError{Code: CodeInvalidRole, Message: "无法更改到此角色类型"}
	ErrCannotModifySelf      = &AppError{Code: CodeCannotModifySelf, Message: "不能修改自己的信息"}
	ErrCannotChangeOwnerRole = &AppError{Code: CodeCannotChangeOwner, Message: "不能修改超级管理员的角色"}
	ErrFollowSelf            = &AppError{Code: CodeFollowSelf, Message: "不能关注自己"}
	ErrAlreadyFollow         = &AppError{Code: CodeAlreadyFollow, Message: "已经关注了该用户"}
	ErrNotFollow             = &AppError{Code: CodeNotFollow, Message: "尚未关注该用户"}
	ErrScoreNotEnough        = &AppError{Code: CodeScoreNotEnough, Message: "积分不足"}

	// 封禁相关
	ErrCannotBlockSelf       = &AppError{Code: CodeCannotModifySelf, Message: "不能封禁自己的账号"}
	ErrCannotBlockSuperAdmin = &AppError{Code: CodeCannotChangeOwner, Message: "不能封禁超级管理员"}
	ErrCannotBlockAdmin      = &AppError{Code: CodeInsufficientPermission, Message: "只有超级管理员才能封禁其他管理员"}

	// 验证相关
	ErrPasswordNotMatch  = &AppError{Code: CodeUnauthorized, Message: "密码不匹配"}
	ErrPasswordTooShort  = &AppError{Code: CodeInvalidPassword, Message: "密码长度至少为6位"}
	ErrPasswordTooLong   = &AppError{Code: CodeInvalidPassword, Message: "密码长度不能超过32位"}
	ErrPasswordSameAsOld = &AppError{Code: CodeInvalidPassword, Message: "新密码不能与旧密码相同"}
	ErrorTokenExpired    = &AppError{Code: CodeUnauthorized, Message: "Token已过期"}
	ErrInvalidToken      = &AppError{Code: CodeUnauthorized, Message: "无效的Token"}
	ErrRequiredToken     = &AppError{Code: CodeValidation, Message: "需要Token"}
	ErrRequiredCaptcha   = &AppError{Code: CodeValidation, Message: "需要验证码"}
	ErrInvalidCaptcha    = &AppError{Code: CodeValidation, Message: "验证码错误"}
	ErrValidationFailed  = &AppError{Code: CodeValidation, Message: "验证失败"}

	// 内容模块 - 帖子
	ErrPostNotFound = &AppError{Code: CodePostNotFound, Message: "帖子不存在"}
	ErrPostLocked   = &AppError{Code: CodePostLocked, Message: "帖子已被锁定，无法操作"}
	ErrPostPinned   = &AppError{Code: CodePostPinned, Message: "帖子置顶状态冲突"}
	ErrPostDeleted  = &AppError{Code: CodePostDeleted, Message: "帖子已被删除"}

	// 内容模块 - 板块
	ErrBoardNotFound = &AppError{Code: CodeBoardNotFound, Message: "板块不存在"}

	// 内容模块 - 问答
	ErrAcceptForbidden  = &AppError{Code: CodeAcceptForbidden, Message: "只有发帖人才能采纳答案"}
	ErrAnswerNotFound   = &AppError{Code: CodeAnswerNotFound, Message: "回答不存在"}
	ErrQuestionNotFound = &AppError{Code: CodeQuestionNotFound, Message: "问题不存在"}

	// 内容模块 - 评论
	ErrCommentNotFound = &AppError{Code: CodeCommentNotFound, Message: "评论不存在"}
	ErrCommentDeleted  = &AppError{Code: CodeCommentDeleted, Message: "评论已被删除"}

	// 内容模块 - 主题（topic）
	ErrTopicNotFound = &AppError{Code: CodeTopicNotFound, Message: "主题不存在"}

	// 内容模块 - 标签
	ErrTagNotFound = &AppError{Code: CodeTagNotFound, Message: "标签不存在"}

	// 内容模块 - 通知
	ErrNotificationNotFound = &AppError{Code: CodeNotificationNotFound, Message: "通知不存在"}

	// 点赞/收藏
	ErrLikeAlready     = &AppError{Code: CodeLikeAlready, Message: "已经点过赞了"}
	ErrLikeNotExist    = &AppError{Code: CodeLikeNotExist, Message: "尚未点赞"}
	ErrCollectAlready  = &AppError{Code: CodeCollectAlready, Message: "已经收藏过了"}
	ErrCollectNotExist = &AppError{Code: CodeCollectNotExist, Message: "尚未收藏"}

	// 权限模块 - 版主申请
	ErrModeratorApplyExist    = &AppError{Code: CodeModeratorApplyExist, Message: "已经提交过版主申请，请勿重复提交"}
	ErrModeratorApplyNotFound = &AppError{Code: CodeModeratorApplyNotFound, Message: "版主申请不存在"}
	ErrAlreadyModerator       = &AppError{Code: CodeAlreadyModerator, Message: "已经是版主，无需重复申请"}
	ErrInsufficientPermission = &AppError{Code: CodeInsufficientPermission, Message: "权限不足"}

	// 积分模块
	ErrFailedToQueryScore  = &AppError{Code: CodeFailedToQueryScore, Message: "查询积分失败"}
	ErrScoreRecordNotFound = &AppError{Code: CodeScoreRecordNotFound, Message: "积分记录不存在"}

	// 公告模块
	ErrAnnouncementNotFound = &AppError{Code: CodeAnnouncementNotFound, Message: "公告不存在"}
	ErrInvalidPublishTime   = &AppError{Code: CodeAnnouncementInvalidTime, Message: "发布时间无效"}
	ErrExpiredTimeInvalid   = &AppError{Code: CodeAnnouncementInvalidTime, Message: "过期时间必须晚于发布时间"}

	// 统计模块
	ErrStatsNotFound = &AppError{Code: CodeStatsNotFound, Message: "统计数据不存在"}

	// 时间线模块
	ErrTimelineEmpty = &AppError{Code: CodeTimelineEmpty, Message: "时间线暂无内容"}

	// 文件上传模块
	ErrFileTooLarge    = &AppError{Code: CodeFileTooLarge, Message: "文件过大"}
	ErrFileTypeInvalid = &AppError{Code: CodeFileTypeInvalid, Message: "不支持的文件类型"}
	ErrUploadFailed    = &AppError{Code: CodeUploadFailed, Message: "文件上传失败"}
)

// ========== 辅助函数 ==========

// Wrap 包装错误，追加上下文信息
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

// WithDetail 附加详情数据
func (e *AppError) WithDetail(detail interface{}) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Detail:  detail,
		Err:     e.Err,
	}
}

// GetCode 获取错误码
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
		CodeCannotModifySelf, CodeCannotChangeOwner, CodeAlreadyFollow,
		CodeAlreadyModerator:
		return http.StatusForbidden
	case CodeUserNotFound, CodePostNotFound, CodeBoardNotFound,
		CodeQuestionNotFound, CodeAnswerNotFound, CodeCommentNotFound,
		CodeTopicNotFound, CodeTagNotFound, CodeNotificationNotFound,
		CodeAnnouncementNotFound, CodeStatsNotFound, CodeModeratorApplyNotFound,
		CodeScoreRecordNotFound, CodeLikeNotExist, CodeCollectNotExist:
		return http.StatusNotFound
	case CodeInvalidEmail, CodeInvalidPhone, CodeInvalidPassword,
		CodeInvalidUsername, CodeInvalidAvatar, CodeInvalidNickname,
		CodeInvalidUserID, CodeInvalidRole, CodeUserExist, CodeValidation,
		CodeFollowSelf, CodeScoreNotEnough, CodeLikeAlready, CodeCollectAlready,
		CodePostLocked, CodePostDeleted, CodeCommentDeleted,
		CodePasswordTooShort, CodePasswordSameAsOld,
		CodeFileTooLarge, CodeFileTypeInvalid:
		return http.StatusBadRequest
	case CodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Is 支持 errors.Is 按错误码比较（可选）
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
