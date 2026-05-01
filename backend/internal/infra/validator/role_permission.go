package validator

import (
	"context"
	"tiny-forum/internal/model/po"
	apperrors "tiny-forum/pkg/errors"
)

type roleValidator struct {
	// 角色等级映射
	roleLevel map[po.UserRole]int
	// 角色权限映射
	rolePermissions map[po.UserRole]map[po.Permission]bool
}

type RoleChangeChecker struct {
	validator RoleValidator
}

func NewRoleValidator() RoleValidator {
	// 定义角色等级（数字越大等级越高）
	roleLevel := map[po.UserRole]int{
		po.RoleUser:       1,
		po.RoleMember:     2,
		po.RoleBot:        2,
		po.RoleReviewer:   3,
		po.RoleModerator:  4,
		po.RoleAdmin:      5,
		po.RoleSuperAdmin: 6,
	}

	// 定义每个角色拥有的权限
	rolePermissions := map[po.UserRole]map[po.Permission]bool{
		po.RoleUser: {
			// 普通用户无分配权限
		},
		po.RoleMember: {
			// 成员用户无分配权限
		},
		po.RoleBot: {
			// Bot用户无分配权限
		},
		po.RoleReviewer: {
			// 审核员可以分配普通用户和成员角色
			po.PermAssignRoleUser:   true,
			po.PermAssignRoleMember: true,
		},
		po.RoleModerator: {
			// 版主可以分配普通用户、成员和审核员角色
			po.PermAssignRoleUser:     true,
			po.PermAssignRoleMember:   true,
			po.PermAssignRoleReviewer: true,
		},
		po.RoleAdmin: {
			// 管理员可以分配普通用户、成员、审核员、版主和Bot角色
			po.PermAssignRoleUser:      true,
			po.PermAssignRoleMember:    true,
			po.PermAssignRoleReviewer:  true,
			po.PermAssignRoleModerator: true,
			po.PermAssignRoleBot:       true,
		},
		po.RoleSuperAdmin: {
			// 超级管理员可以分配所有角色
			po.PermAssignRoleUser:       true,
			po.PermAssignRoleMember:     true,
			po.PermAssignRoleReviewer:   true,
			po.PermAssignRoleModerator:  true,
			po.PermAssignRoleBot:        true,
			po.PermAssignRoleAdmin:      true,
			po.PermAssignRoleSuperAdmin: true,
		},
	}

	return &roleValidator{
		roleLevel:       roleLevel,
		rolePermissions: rolePermissions,
	}
}

// 角色更改请求
type RoleChangeRequest struct {
	Operator *po.User
	Target   *po.User
	NewRole  po.UserRole
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

var roleAssignPermission = map[po.UserRole]po.Permission{
	po.RoleUser:       po.PermAssignRoleUser,
	po.RoleMember:     po.PermAssignRoleMember,
	po.RoleModerator:  po.PermAssignRoleModerator,
	po.RoleReviewer:   po.PermAssignRoleReviewer,
	po.RoleBot:        po.PermAssignRoleBot,
	po.RoleAdmin:      po.PermAssignRoleAdmin,
	po.RoleSuperAdmin: po.PermAssignRoleSuperAdmin,
}

// IsValidRole 检查角色是否有效
func (v *roleValidator) IsValidRole(role po.UserRole) bool {
	_, exists := v.roleLevel[role]
	return exists
}

// CanOperateTarget 检查操作者是否有权限操作目标用户
// 规则：操作者等级必须高于目标用户等级
func (v *roleValidator) CanOperateTarget(operator, target po.UserRole) bool {
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
func (v *roleValidator) CanAssignRole(operator, target po.UserRole) bool {
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
func (v *roleValidator) HasPermission(role po.UserRole, perm po.Permission) bool {
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
func (v *roleValidator) GetRoleLevel(role po.UserRole) (int, bool) {
	level, exists := v.roleLevel[role]
	return level, exists
}

// GetRolePermissions 获取角色的所有权限
func (v *roleValidator) GetRolePermissions(role po.UserRole) []po.Permission {
	perms, exists := v.rolePermissions[role]
	if !exists {
		return []po.Permission{}
	}

	result := make([]po.Permission, 0, len(perms))
	for perm := range perms {
		result = append(result, perm)
	}
	return result
}

// CanModifyRole 综合检查是否可以修改角色（包含所有规则）
func (v *roleValidator) CanModifyRole(operator, currentRole, newRole po.UserRole) bool {
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
