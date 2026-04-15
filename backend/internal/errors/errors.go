package apperrors

import "errors"

var (
	// 内容
	ErrPostNotFound = errors.New("帖子不存在")
	// 角色
	ErrInvalidRole           = errors.New("无法更改到此角色类型")
	ErrCannotChangeOwnerRole = errors.New("不能修改超级管理员的角色")
	ErrCannotModifySelf      = errors.New("不能修改自己的角色")
	// 权限
	ErrInsufficientPermission = errors.New("权限不足")
	// 板块
	ErrBoardNotFound = errors.New("板块不存在")
	// 问答
	ErrAcceptForbidden  = errors.New("只有发帖人才能采纳")
	ErrAnswerNotFound   = errors.New("回答不存在")
	ErrQuestionNotFound = errors.New("问题不存在")
	// 用户信息
	ErrUserNotFound    = errors.New("用户不存在")
	ErrUserExist       = errors.New("用户已存在")
	ErrInvalidEmail    = errors.New("无效的邮箱")
	ErrInvalidPhone    = errors.New("无效的手机号")
	ErrInvalidPassword = errors.New("无效的密码")
	ErrInvalidUsername = errors.New("无效的用户名")
	ErrInvalidAvatar   = errors.New("无效的头像")
	ErrInvalidNickname = errors.New("无效的昵称")
	ErrInvalidUserID   = errors.New("无效的用户ID")
	// 积分
	ErrFailedToQueryScore = errors.New("查询积分失败")
)
