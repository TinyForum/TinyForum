package service

import (
	"fmt"
	"tiny-forum/internal/model"

	apperrors "tiny-forum/internal/errors"
)

// RoleChangeRequest 封装一次角色变更的上下文
type RoleChangeRequest struct {
	Operator *model.User
	Target   *model.User
	NewRole  model.UserRole
}

// RoleChangeChecker 细粒度角色变更权限校验器
type RoleChangeChecker struct{}

// Check 执行全部权限校验，返回第一个不满足的错误
func (RoleChangeChecker) Check(req RoleChangeRequest) error {
	checks := []func(RoleChangeRequest) error{
		checkNotSelf,
		checkNewRoleValid,
		checkCanOperateTarget,
		checkCanAssignRole,
		checkAssignPermission,
	}
	for _, fn := range checks {
		if err := fn(req); err != nil {
			return err
		}
	}
	return nil
}

// ── 各个独立校验规则 ─────────────────────────────────────────────────────────

// checkNotSelf 不允许修改自己的角色
func checkNotSelf(req RoleChangeRequest) error {
	if req.Operator.ID == req.Target.ID {
		return apperrors.ErrCannotModifySelf
	}
	return nil
}

// checkNewRoleValid 目标角色必须是合法值
func checkNewRoleValid(req RoleChangeRequest) error {
	if !model.IsValidRole(req.NewRole) {
		return fmt.Errorf("%w: %s", apperrors.ErrInvalidRole, req.NewRole)
	}
	return nil
}

// checkCanOperateTarget 操作者是否有权限操作该目标（基于目标当前角色）
func checkCanOperateTarget(req RoleChangeRequest) error {
	if !model.CanOperateTarget(req.Operator.Role, req.Target.Role) {
		return fmt.Errorf("%w: 无权操作角色为 %s 的用户",
			apperrors.ErrInsufficientPermission, req.Target.Role)
	}
	return nil
}

// checkCanAssignRole 操作者是否可以分配目标角色（基于矩阵）
func checkCanAssignRole(req RoleChangeRequest) error {
	if !model.CanAssignRole(req.Operator.Role, req.NewRole) {
		return fmt.Errorf("%w: 无权将用户角色设置为 %s",
			apperrors.ErrInsufficientPermission, req.NewRole)
	}
	return nil
}

// checkAssignPermission 操作者是否拥有对应角色的 role.assign.* 权限
// 这是最细粒度的一层：即使矩阵允许，也需要具体权限节点
func checkAssignPermission(req RoleChangeRequest) error {
	perm, ok := roleAssignPermission[req.NewRole]
	if !ok {
		// 没有对应权限节点，直接拒绝
		return fmt.Errorf("%w: 目标角色 %s 不支持分配",
			apperrors.ErrInsufficientPermission, req.NewRole)
	}
	if !model.HasPermission(req.Operator.Role, perm) {
		return fmt.Errorf("%w: 缺少权限 %s",
			apperrors.ErrInsufficientPermission, perm)
	}
	return nil
}

// roleAssignPermission 目标角色→所需 role.assign.* 权限节点
var roleAssignPermission = map[model.UserRole]model.Permission{
	model.RoleUser:       model.PermAssignRoleUser,
	model.RoleMember:     model.PermAssignRoleMember,
	model.RoleModerator:  model.PermAssignRoleModerator,
	model.RoleReviewer:   model.PermAssignRoleReviewer,
	model.RoleBot:        model.PermAssignRoleBot,
	model.RoleAdmin:      model.PermAssignRoleAdmin,
	model.RoleSuperAdmin: model.PermAssignRoleSuperAdmin,
}
