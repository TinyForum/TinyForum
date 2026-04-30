package user

import (
	"context"
	"errors"
	"fmt"
	"tiny-forum/internal/infra/validator"
	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// SetBlocked 管理员封禁/解封用户
func (s *userService) SetBlocked(targetID uint, operatorID uint, isBlocked bool) error {
	ctx := context.Background()

	targetUser, err := s.repo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound.WithMessagef( "ID: %d", targetID)
		}
		return fmt.Errorf("查询目标用户失败: %w", err)
	}

	operator, err := s.repo.FindByID(operatorID)
	if err != nil {
		return fmt.Errorf("查询操作者信息失败: %w", err)
	}

	if targetID == operatorID {
		return apperrors.ErrCannotModifySelf.WithMessage( "不能封禁自己的账号")
	}
	if targetUser.Role == model.RoleSuperAdmin {
		return apperrors.ErrCannotChangeOwnerRole.WithMessage( "不能封禁超级管理员")
	}
	if operator.Role != model.RoleSuperAdmin && targetUser.Role == model.RoleAdmin {
		return apperrors.ErrInsufficientPermission.WithMessage("只有超级管理员才能封禁其他管理员")
	}
	if targetUser.IsBlocked == isBlocked {
		return nil
	}

	if err := s.repo.UpdateBlocked(ctx, targetID, isBlocked); err != nil {
		return fmt.Errorf("更新用户封禁状态失败: %w", err)
	}

	action := "unblock_user"
	if isBlocked {
		action = "block_user"
	}
	s.logAudit(ctx, operatorID, targetID, action,
		fmt.Sprintf("管理员 %d %s用户 %d (%s)", operatorID, map[bool]string{true: "封禁", false: "解封"}[isBlocked], targetID, targetUser.Username))
	return nil
}

// SetActive 管理员设置用户激活状态
func (s *userService) SetActive(targetID uint, operatorID uint, isActive bool) error {
	ctx := context.Background()

	targetUser, err := s.repo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound.WithMessagef( "ID: %d", targetID)
		}
		return fmt.Errorf("查询目标用户失败: %w", err)
	}

	if targetID == operatorID && !isActive {
		return apperrors.ErrCannotModifySelf.WithMessage( "不能停用自己的账号")
	}
	if targetUser.Role == model.RoleSuperAdmin && !isActive {
		return apperrors.ErrCannotChangeOwnerRole.WithMessage( "不能停用超级管理员")
	}
	if targetUser.IsActive == isActive {
		return nil
	}

	if err := s.repo.UpdateActive(ctx, targetID, isActive); err != nil {
		return fmt.Errorf("更新用户激活状态失败: %w", err)
	}
	if !isActive {
		_ = s.repo.InvalidateUserTokens(ctx, targetID)
	}

	action := "activate_user"
	if !isActive {
		action = "deactivate_user"
	}
	s.logAudit(ctx, operatorID, targetID, action,
		fmt.Sprintf("管理员 %d %s用户 %d (%s)", operatorID, map[bool]string{true: "激活", false: "停用"}[isActive], targetID, targetUser.Username))
	return nil
}

// SetRole 变更用户角色
func (s *userService) SetRole(operatorID, targetID uint, newRole string) error {
	ctx := context.Background()
	targetRole := model.UserRole(newRole)
	if !model.IsValidRole(targetRole) {
		return fmt.Errorf("%w: %s", apperrors.ErrInvalidRole, newRole)
	}

	operator, err := s.repo.FindByID(operatorID)
	if err != nil {
		return fmt.Errorf("操作者不存在: %w", err)
	}
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return err
	}
	if target.Role == targetRole {
		return nil
	}

	if err := s.roleChecker.Check(ctx, validator.RoleChangeRequest{
		Operator: operator,
		Target:   target,
		NewRole:  targetRole,
	}); err != nil {
		return err
	}

	return s.repo.UpdateFields(targetID, map[string]interface{}{"role": newRole})
}
