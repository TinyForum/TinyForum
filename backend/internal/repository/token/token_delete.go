package token

import (
	"context"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// DeleteByUserID 删除用户的所有 Token（全局登出、密码重置、封禁时调用）
func (r *tokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.RefreshToken{}).Error
}

// DeleteByToken 删除单个 Token（单设备登出）
func (r *tokenRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&model.RefreshToken{}).Error
}

// DeleteByJTI 根据 JTI 精确撤销
func (r *tokenRepository) DeleteByJTI(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).
		Where("jti = ?", jti).
		Delete(&model.RefreshToken{}).Error
}

// DeleteExpired 清理过期的 Token（建议用定时任务每天执行）
func (r *tokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at <= ?", time.Now()).
		Delete(&model.RefreshToken{})
	return result.RowsAffected, result.Error
}

// DeleteByUserIDExcept 删除用户的所有 Token，但保留指定的一个（切换设备时用）
func (r *tokenRepository) DeleteByUserIDExcept(ctx context.Context, userID uint, keepToken string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND token != ?", userID, keepToken).
		Delete(&model.RefreshToken{}).Error
}

func (r *tokenRepository) DeleteByUserIDWithTx(ctx context.Context, tx *gorm.DB, userID uint) error {
	return tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}
