package user

import (
	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"
)

// RoleChangeRequest 角色变更请求
type RoleChangeRequest struct {
	Operator *model.User
	Target   *model.User
	NewRole  model.UserRole
}

// RoleChangeChecker 角色变更权限检查器
type RoleChangeChecker struct{}

// Check 执行权限检查
func (c RoleChangeChecker) Check(req RoleChangeRequest) error {
	operatorRole := req.Operator.Role
	targetRole := req.Target.Role

	// 不能修改自己的角色
	if req.Operator.ID == req.Target.ID {
		return apperrors.ErrCannotModifySelf
	}
	// 不能修改超级管理员的角色
	if targetRole == model.RoleSuperAdmin {
		return apperrors.ErrCannotChangeOwnerRole
	}
	// 只有超级管理员可以修改管理员角色
	if targetRole == model.RoleAdmin && operatorRole != model.RoleSuperAdmin {
		return apperrors.ErrInsufficientPermission
	}
	// 超级管理员可以修改任何角色
	if operatorRole == model.RoleSuperAdmin {
		return nil
	}
	// 普通管理员只能将普通用户提升为版主等，但不能提升为管理员
	if operatorRole == model.RoleAdmin {
		if req.NewRole == model.RoleAdmin || req.NewRole == model.RoleSuperAdmin {
			return apperrors.ErrInsufficientPermission
		}
		return nil
	}
	// 普通用户不能修改任何角色
	return apperrors.ErrInsufficientPermission
}
