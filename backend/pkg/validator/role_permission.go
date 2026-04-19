// package validator

// import (
// 	"context"
// 	"tiny-forum/internal/model"
// 	apperrors "tiny-forum/pkg/errors"
// )

// type RoleValidator interface {
// 	IsValidRole(role model.UserRole) bool
// 	CanOperateTarget(operator, target model.UserRole) bool
// 	CanAssignRole(operator, target model.UserRole) bool
// 	HasPermission(role model.UserRole, perm model.Permission) bool
// }

// type roleValidator struct {
// 	// 可以添加配置项
// }

// type RoleChangeChecker struct {
// 	validator RoleValidator
// }

// func NewRoleValidator() RoleValidator {
// 	return &roleValidator{}
// }
// func NewRoleChangeChecker(validator RoleValidator) *RoleChangeChecker {
// 	return &RoleChangeChecker{validator: validator}
// }

// type RoleChangeRequest struct {
// 	Operator *model.User
// 	Target   *model.User
// 	NewRole  model.UserRole
// }

// // Check 执行全部权限校验，返回第一个不满足的错误
// func (c *RoleChangeChecker) Check(ctx context.Context, req RoleChangeRequest) error {
// 	checks := []func(context.Context, RoleChangeRequest) error{
// 		c.checkNotSelf,
// 		c.checkNewRoleValid,
// 		c.checkCanOperateTarget,
// 		c.checkCanAssignRole,
// 		c.checkAssignPermission,
// 	}
// 	for _, fn := range checks {
// 		if err := fn(ctx, req); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (c *RoleChangeChecker) checkNotSelf(_ context.Context, req RoleChangeRequest) error {
// 	if req.Operator.ID == req.Target.ID {
// 		return apperrors.ErrCannotModifySelf
// 	}
// 	return nil
// }

// func (c *RoleChangeChecker) checkNewRoleValid(_ context.Context, req RoleChangeRequest) error {
// 	if !c.validator.IsValidRole(req.NewRole) {
// 		return apperrors.Wrapf(apperrors.ErrInvalidRole, "无效角色: %s", req.NewRole)
// 	}
// 	return nil
// }

// func (c *RoleChangeChecker) checkCanOperateTarget(_ context.Context, req RoleChangeRequest) error {
// 	if !c.validator.CanOperateTarget(req.Operator.Role, req.Target.Role) {
// 		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
// 			"无权操作角色为 %s 的用户", req.Target.Role)
// 	}
// 	return nil
// }

// func (c *RoleChangeChecker) checkCanAssignRole(_ context.Context, req RoleChangeRequest) error {
// 	if !c.validator.CanAssignRole(req.Operator.Role, req.NewRole) {
// 		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
// 			"无权将用户角色设置为 %s", req.NewRole)
// 	}
// 	return nil
// }

// func (c *RoleChangeChecker) checkAssignPermission(_ context.Context, req RoleChangeRequest) error {
// 	perm, ok := roleAssignPermission[req.NewRole]
// 	if !ok {
// 		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
// 			"目标角色 %s 不支持分配", req.NewRole)
// 	}
// 	if !c.validator.HasPermission(req.Operator.Role, perm) {
// 		return apperrors.Wrapf(apperrors.ErrInsufficientPermission,
// 			"缺少权限 %s", perm)
// 	}
// 	return nil
// }

// var roleAssignPermission = map[model.UserRole]model.Permission{
// 	model.RoleUser:       model.PermAssignRoleUser,
// 	model.RoleMember:     model.PermAssignRoleMember,
// 	model.RoleModerator:  model.PermAssignRoleModerator,
// 	model.RoleReviewer:   model.PermAssignRoleReviewer,
// 	model.RoleBot:        model.PermAssignRoleBot,
// 	model.RoleAdmin:      model.PermAssignRoleAdmin,
// 	model.RoleSuperAdmin: model.PermAssignRoleSuperAdmin,
// }

package validator
