package auth

import (
	"context"
	"errors"
	"time"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

// "tiny-forum/internal/model/po

// 验证密码是否有效
// - 用户新密码
// - 用户id
func (r *authRepository) ValidateResetToken(ctx context.Context, token string) (bool, error) {
	var refreshToken po.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		Where("jti LIKE ?", "reset_%").
		Where("deleted_at IS NULL").
		First(&refreshToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // token 不存在
		}
		return false, err
	}

	// 检查是否过期
	if time.Now().After(refreshToken.ExpiresAt) {
		return false, nil // token 已过期
	}

	// 可选：检查是否已使用
	if refreshToken.IsUsed {
		return false, nil // token 已使用
	}

	return true, nil
}
