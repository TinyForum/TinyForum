package user

import (
	"context"
	"errors"
	"fmt"
	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// List 管理员获取用户列表（分页）
func (s *UserService) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

// DeleteUser 管理员删除用户（软删除）
func (s *UserService) DeleteUser(operatorID uint, targetID uint) error {
	ctx := context.Background()

	targetUser, err := s.repo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
		}
		return fmt.Errorf("查询目标用户失败: %w", err)
	}

	operator, err := s.repo.FindByID(operatorID)
	if err != nil {
		return fmt.Errorf("查询操作者信息失败: %w", err)
	}

	if targetID == operatorID {
		return apperrors.Wrap(apperrors.ErrCannotModifySelf, "不能删除自己的账号")
	}
	if targetUser.Role == model.RoleSuperAdmin {
		return apperrors.Wrap(apperrors.ErrCannotChangeOwnerRole, "不能删除超级管理员")
	}
	if operator.Role != model.RoleSuperAdmin && targetUser.Role == model.RoleAdmin {
		return apperrors.Wrap(apperrors.ErrInsufficientPermission, "只有超级管理员才能删除其他管理员")
	}

	if err := s.repo.SoftDelete(ctx, targetID); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	s.logAudit(ctx, operatorID, targetID, "delete_user",
		fmt.Sprintf("管理员 %d 删除了用户 %d (%s)", operatorID, targetID, targetUser.Username))
	return nil
}
