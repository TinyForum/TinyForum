package validator

import (
	"context"
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
)

// RoleValidator 角色权限校验器接口
type roleValidator struct {
	// 角色等级映射
	roleLevel map[do.UserRole]int
	// 角色权限映射
	rolePermissions map[do.UserRole]map[do.Permission]bool
}

// RoleValidator 角色权限校验器接口
type RoleChangeChecker struct {
	validator RoleValidator
}

// NewRoleValidator 创建角色权限校验器
func NewRoleValidator() RoleValidator {
	// 定义角色等级（数字越大等级越高）
	roleLevel := map[do.UserRole]int{
		do.RoleUser:       1,
		do.RoleMember:     2,
		do.RoleBot:        2,
		do.RoleReviewer:   3,
		do.RoleModerator:  4,
		do.RoleAdmin:      5,
		do.RoleSystemMaintainer: 6,
		do.RoleSuperAdmin: 7,
		
	}

	// 定义每个角色拥有的权限
	rolePermissions := map[do.UserRole]map[do.Permission]bool{
		do.RoleUser: {
			// 普通用户无分配权限
		},
		do.RoleMember: {
			// 成员用户无分配权限
		},
		do.RoleBot: {
			// Bot用户无分配权限
		},
		do.RoleReviewer: {
			// 审核员可以分配普通用户和成员角色
			do.PermAssignRoleUser:   true,
			do.PermAssignRoleMember: true,
		},
		do.RoleModerator: {
			// 版主可以分配普通用户、成员和审核员角色
			do.PermAssignRoleUser:     true,
			do.PermAssignRoleMember:   true,
			do.PermAssignRoleReviewer: true,
		},
		do.RoleAdmin: {
			// 管理员可以分配普通用户、成员、审核员、版主和Bot角色
			do.PermAssignRoleUser:      true,
			do.PermAssignRoleMember:    true,
			do.PermAssignRoleReviewer:  true,
			do.PermAssignRoleModerator: true,
			do.PermAssignRoleBot:       true,
		},
		do.RoleSuperAdmin: {
			// 超级管理员可以分配所有角色
			do.PermAssignRoleUser:       true,
			do.PermAssignRoleMember:     true,
			do.PermAssignRoleReviewer:   true,
			do.PermAssignRoleModerator:  true,
			do.PermAssignRoleBot:        true,
			do.PermAssignRoleAdmin:      true,
			do.PermAssignRoleSuperAdmin: true,
			do.PermAssignRoleSystemMaintainer: true,
		},
	}

	return &roleValidator{
		roleLevel:       roleLevel,
		rolePermissions: rolePermissions,
	}
}

// 角色更改请求
type RoleChangeRequest struct {
	Operator *do.User
	Target   *do.User
	NewRole  do.UserRole
}

// Check 执行全部权限校验，返回第一个不满足的错误
func (c *RoleChangeChecker) Check(ctx context.Context, req RoleChangeRequest) error {
	checks := []func(context.Context, RoleChangeRequest) error{
		c.checkNotSelf,
		c.checkNewRoleValid,
		c.checkCanOperateTarget,
		c.checkCanAssignRole,
		c.checkAssignPermission,
	}
	for _, fn := range checks {
		if err := fn(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func (c *RoleChangeChecker) checkNotSelf(_ context.Context, req RoleChangeRequest) error {
	if req.Operator.ID == req.Target.ID {
		return apperrors.ErrCannotModifySelf
	}
	return nil
}

func (c *RoleChangeChecker) checkNewRoleValid(_ context.Context, req RoleChangeRequest) error {
	if !c.validator.IsValidRole(req.NewRole) {
		return apperrors.ErrInvalidRole.WithMessagef("无效角色: %s", req.NewRole)
	}
	return nil
}

func (c *RoleChangeChecker) checkCanOperateTarget(_ context.Context, req RoleChangeRequest) error {
	if !c.validator.CanOperateTarget(req.Operator.Role, req.Target.Role) {
		return apperrors.ErrInsufficientPermission.WithMessagef("无权操作角色为 %v 的用户", req.Target.Role)
	}
	return nil
}

func (c *RoleChangeChecker) checkCanAssignRole(_ context.Context, req RoleChangeRequest) error {
	if !c.validator.CanAssignRole(req.Operator.Role, req.NewRole) {
		return apperrors.ErrInsufficientPermission.WithMessagef(
			"无权将用户角色设置为 %s", req.NewRole)
	}
	return nil
}

// checkAssignPermission 检查操作者是否有权限分配目标角色
func (c *RoleChangeChecker) checkAssignPermission(_ context.Context, req RoleChangeRequest) error {
	perm, ok := roleAssignPermission[req.NewRole]
	if !ok {
		return apperrors.ErrInsufficientPermission.WithMessagef(
			"目标角色 %s 不支持分配", req.NewRole)
	}
	if !c.validator.HasPermission(req.Operator.Role, perm) {
		return apperrors.ErrInsufficientPermission.WithMessagef(
			"缺少权限 %s", perm)
	}
	return nil
}

var roleAssignPermission = map[do.UserRole]do.Permission{
	do.RoleUser:       do.PermAssignRoleUser,
	do.RoleMember:     do.PermAssignRoleMember,
	do.RoleModerator:  do.PermAssignRoleModerator,
	do.RoleReviewer:   do.PermAssignRoleReviewer,
	do.RoleBot:        do.PermAssignRoleBot,
	do.RoleAdmin:      do.PermAssignRoleAdmin,
	do.RoleSuperAdmin: do.PermAssignRoleSuperAdmin,
	do.RoleSystemMaintainer: do.PermAssignRoleSystemMaintainer,
}

// IsValidRole 检查角色是否有效
func (v *roleValidator) IsValidRole(role do.UserRole) bool {
	_, exists := v.roleLevel[role]
	return exists
}

// CanOperateTarget 检查操作者是否有权限操作目标用户
// 规则：操作者等级必须高于目标用户等级
func (v *roleValidator) CanOperateTarget(operator, target do.UserRole) bool {
	operatorLevel, operatorExists := v.roleLevel[operator]
	targetLevel, targetExists := v.roleLevel[target]

	if !operatorExists || !targetExists {
		return false
	}

	// 操作者等级必须大于目标等级
	return operatorLevel > targetLevel
}

// CanAssignRole 检查操作者是否有权限分配目标角色
// 规则：操作者等级必须大于等于目标角色等级（不能给自己分配同级或更高级别）
func (v *roleValidator) CanAssignRole(operator, target do.UserRole) bool {
	operatorLevel, operatorExists := v.roleLevel[operator]
	targetLevel, targetExists := v.roleLevel[target]

	if !operatorExists || !targetExists {
		return false
	}

	// 操作者等级必须大于等于目标角色等级
	// 注意：超级管理员可以分配管理员，但管理员不能分配超级管理员
	return operatorLevel >= targetLevel
}

// HasPermission 检查角色是否拥有特定权限
func (v *roleValidator) HasPermission(role do.UserRole, perm do.Permission) bool {
	perms, exists := v.rolePermissions[role]
	if !exists {
		return false
	}

	// 检查是否拥有该权限
	has, ok := perms[perm]
	return ok && has
}

// 可选：添加一些辅助方法

// GetRoleLevel 获取角色等级
func (v *roleValidator) GetRoleLevel(role do.UserRole) (int, bool) {
	level, exists := v.roleLevel[role]
	return level, exists
}

// GetRolePermissions 获取角色的所有权限
func (v *roleValidator) GetRolePermissions(role do.UserRole) []do.Permission {
	perms, exists := v.rolePermissions[role]
	if !exists {
		return []do.Permission{}
	}

	result := make([]do.Permission, 0, len(perms))
	for perm := range perms {
		result = append(result, perm)
	}
	return result
}

// CanModifyRole 综合检查是否可以修改角色（包含所有规则）
func (v *roleValidator) CanModifyRole(operator, currentRole, newRole do.UserRole) bool {
	// 检查是否可以操作目标
	if !v.CanOperateTarget(operator, currentRole) {
		return false
	}

	// 检查是否可以分配新角色
	if !v.CanAssignRole(operator, newRole) {
		return false
	}

	// 检查是否有分配该角色的具体权限
	perm, ok := roleAssignPermission[newRole]
	if !ok {
		return false
	}

	return v.HasPermission(operator, perm)
}
