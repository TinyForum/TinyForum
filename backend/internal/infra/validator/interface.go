package validator

import "tiny-forum/internal/model/do"

// 验证器

type RoleValidator interface {
	IsValidRole(role do.UserRole) bool
	CanOperateTarget(operator, target do.UserRole) bool
	CanAssignRole(operator, target do.UserRole) bool
	HasPermission(role do.UserRole, perm do.Permission) bool
}

// 创建角色变更校验器
func NewRoleChangeChecker(validator RoleValidator) *RoleChangeChecker {
	return &RoleChangeChecker{validator: validator}
}
