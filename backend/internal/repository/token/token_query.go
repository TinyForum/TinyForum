package token

import (
	"context"
	"time"
	"tiny-forum/internal/model"
)

// ListByUserID 获取用户的所有登录设备（仅未过期的）
func (r *tokenRepository) ListByUserID(ctx context.Context, userID uint) ([]model.RefreshToken, error) {
	var tokens []model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// CountByUserID 统计用户的有效 Token 数量
func (r *tokenRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, err
}
