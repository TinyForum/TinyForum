package validator

import (
	"context"
	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"
)

type roleValidator struct {
	// 角色等级映射
	roleLevel map[model.UserRole]int
	// 角色权限映射
	rolePermissions map[model.UserRole]map[model.Permission]bool
}

type RoleChangeChecker struct {
	validator RoleValidator
}

func NewRoleValidator() RoleValidator {
	// 定义角色等级（数字越大等级越高）
	roleLevel := map[model.UserRole]int{
		model.RoleUser:       1,
		model.RoleMember:     2,
		model.RoleBot:        2,
		model.RoleReviewer:   3,
		model.RoleModerator:  4,
		model.RoleAdmin:      5,
		model.RoleSuperAdmin: 6,
	}

	// 定义每个角色拥有的权限
	rolePermissions := map[model.UserRole]map[model.Permission]bool{
		model.RoleUser: {
			// 普通用户无分配权限
		},
		model.RoleMember: {
			// 成员用户无分配权限
		},
		model.RoleBot: {
			// Bot用户无分配权限
		},
		model.RoleReviewer: {
			// 审核员可以分配普通用户和成员角色
			model.PermAssignRoleUser:   true,
			model.PermAssignRoleMember: true,
		},
		model.RoleModerator: {
			// 版主可以分配普通用户、成员和审核员角色
			model.PermAssignRoleUser:     true,
			model.PermAssignRoleMember:   true,
			model.PermAssignRoleReviewer: true,
		},
		model.RoleAdmin: {
			// 管理员可以分配普通用户、成员、审核员、版主和Bot角色
			model.PermAssignRoleUser:      true,
			model.PermAssignRoleMember:    true,
			model.PermAssignRoleReviewer:  true,
			model.PermAssignRoleModerator: true,
			model.PermAssignRoleBot:       true,
		},
		model.RoleSuperAdmin: {
			// 超级管理员可以分配所有角色
			model.PermAssignRoleUser:       true,
			model.PermAssignRoleMember:     true,
			model.PermAssignRoleReviewer:   true,
			model.PermAssignRoleModerator:  true,
			model.PermAssignRoleBot:        true,
			model.PermAssignRoleAdmin:      true,
			model.PermAssignRoleSuperAdmin: true,
		},
	}

	return &roleValidator{
		roleLevel:       roleLevel,
		rolePermissions: rolePermissions,
	}
}

// 角色更改请求
type RoleChangeRequest struct {
	Operator *model.User
	Target   *model.User
	NewRole  model.UserRole
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
		return apperrors.Wrapf(apperrors.ErrInvalidRole, "无效角色: %s", req.NewRole)
	}
	return nil
}

func (c *RoleChangeChecker) checkCanOperateTarget(_ context.Context, req RoleChangeRequest) error {
	if !c.validator.CanOperateTarget(req.Operator.Role, req.Target.Role) {
		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
			"无权操作角色为 %s 的用户", req.Target.Role)
	}
	return nil
}

func (c *RoleChangeChecker) checkCanAssignRole(_ context.Context, req RoleChangeRequest) error {
	if !c.validator.CanAssignRole(req.Operator.Role, req.NewRole) {
		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
			"无权将用户角色设置为 %s", req.NewRole)
	}
	return nil
}

func (c *RoleChangeChecker) checkAssignPermission(_ context.Context, req RoleChangeRequest) error {
	perm, ok := roleAssignPermission[req.NewRole]
	if !ok {
		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
			"目标角色 %s 不支持分配", req.NewRole)
	}
	if !c.validator.HasPermission(req.Operator.Role, perm) {
		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
			"缺少权限 %s", perm)
	}
	return nil
}

var roleAssignPermission = map[model.UserRole]model.Permission{
	model.RoleUser:       model.PermAssignRoleUser,
	model.RoleMember:     model.PermAssignRoleMember,
	model.RoleModerator:  model.PermAssignRoleModerator,
	model.RoleReviewer:   model.PermAssignRoleReviewer,
	model.RoleBot:        model.PermAssignRoleBot,
	model.RoleAdmin:      model.PermAssignRoleAdmin,
	model.RoleSuperAdmin: model.PermAssignRoleSuperAdmin,
}

// IsValidRole 检查角色是否有效
func (v *roleValidator) IsValidRole(role model.UserRole) bool {
	_, exists := v.roleLevel[role]
	return exists
}

// CanOperateTarget 检查操作者是否有权限操作目标用户
// 规则：操作者等级必须高于目标用户等级
func (v *roleValidator) CanOperateTarget(operator, target model.UserRole) bool {
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
func (v *roleValidator) CanAssignRole(operator, target model.UserRole) bool {
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
func (v *roleValidator) HasPermission(role model.UserRole, perm model.Permission) bool {
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
func (v *roleValidator) GetRoleLevel(role model.UserRole) (int, bool) {
	level, exists := v.roleLevel[role]
	return level, exists
}

// GetRolePermissions 获取角色的所有权限
func (v *roleValidator) GetRolePermissions(role model.UserRole) []model.Permission {
	perms, exists := v.rolePermissions[role]
	if !exists {
		return []model.Permission{}
	}

	result := make([]model.Permission, 0, len(perms))
	for perm := range perms {
		result = append(result, perm)
	}
	return result
}

// CanModifyRole 综合检查是否可以修改角色（包含所有规则）
func (v *roleValidator) CanModifyRole(operator, currentRole, newRole model.UserRole) bool {
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
