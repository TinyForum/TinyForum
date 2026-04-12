package apperrors

import "errors"

var (
	ErrPostNotFound           = errors.New("帖子不存在")
	ErrInvalidRole            = errors.New("无法更改到此角色类型")
	ErrCannotChangeOwnerRole  = errors.New("不能修改超级管理员的角色")
	ErrCannotModifySelf       = errors.New("不能修改自己的角色")
	ErrInsufficientPermission = errors.New("权限不足")
	ErrBoardNotFound          = errors.New("板块不存在")
)
