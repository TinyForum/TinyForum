package user

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tiny-forum/internal/model/do"
)

func (r *userRepository) UpdateBlocked(ctx context.Context, userID uint, isBlocked bool) error {
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Update("is_blocked", isBlocked)
	if result.Error != nil {
		return fmt.Errorf("更新用户封禁状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	if isBlocked {
		_ = r.tokenRepo.DeleteByUserID(ctx, userID)
	}
	return nil
}

func (r *userRepository) UpdateActive(ctx context.Context, userID uint, isActive bool) error {
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Update("is_active", isActive)
	if result.Error != nil {
		return fmt.Errorf("更新用户激活状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func (r *userRepository) UpdateRole(ctx context.Context, userID uint, role string) error {
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Update("role", role)
	if result.Error != nil {
		return fmt.Errorf("更新用户角色失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func (r *userRepository) SoftDelete(ctx context.Context, userID uint) error {
	_ = r.tokenRepo.DeleteByUserID(ctx, userID)
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("软删除用户失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在或已被删除")
	}
	return nil
}

func (r *userRepository) HardDelete(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Unscoped().
		Delete(&do.User{}, userID)
	if result.Error != nil {
		return fmt.Errorf("硬删除用户失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword)
	if result.Error != nil {
		return fmt.Errorf("更新密码失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	_ = r.tokenRepo.DeleteByUserID(ctx, userID)
	return nil
}

func (r *userRepository) InvalidateUserTokens(ctx context.Context, userID uint) error {
	return r.tokenRepo.DeleteByUserID(ctx, userID)
}

func (r *userRepository) RestoreDeleted(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("恢复用户失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在或未被删除")
	}
	return nil
}

func (r *userRepository) SetTempPasswordFlag(ctx context.Context, userID uint, isTemp bool, expireAt time.Time) error {
	updates := map[string]interface{}{
		"is_temp_password":     isTemp,
		"temp_password_expire": expireAt,
	}
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("设置临时密码标记失败: %w", result.Error)
	}
	return nil
}

func (r *userRepository) ClearTempPasswordFlag(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_temp_password":     false,
			"temp_password_expire": nil,
		}).Error
}

func (r *userRepository) BatchUpdateBlocked(ctx context.Context, userIDs []uint, isBlocked bool) (int64, error) {
	if len(userIDs) == 0 {
		return 0, nil
	}
	result := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("id IN ?", userIDs).
		Update("is_blocked", isBlocked)
	return result.RowsAffected, result.Error
}
