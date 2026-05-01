package validator

// 验证器
import (
	"tiny-forum/internal/model/po"
)

type RoleValidator interface {
	IsValidRole(role po.UserRole) bool
	CanOperateTarget(operator, target po.UserRole) bool
	CanAssignRole(operator, target po.UserRole) bool
	HasPermission(role po.UserRole, perm po.Permission) bool
}

// 创建角色变更校验器
func NewRoleChangeChecker(validator RoleValidator) *RoleChangeChecker {
	return &RoleChangeChecker{validator: validator}
}
