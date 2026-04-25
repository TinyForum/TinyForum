package validator

// 验证器
import (
	"tiny-forum/internal/model"
)

type RoleValidator interface {
	IsValidRole(role model.UserRole) bool
	CanOperateTarget(operator, target model.UserRole) bool
	CanAssignRole(operator, target model.UserRole) bool
	HasPermission(role model.UserRole, perm model.Permission) bool
}

// 创建角色变更校验器
func NewRoleChangeChecker(validator RoleValidator) *RoleChangeChecker {
	return &RoleChangeChecker{validator: validator}
}
