package token

import (
	"context"
	"errors"
	"time"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// ListByUserID 获取用户的所有登录设备（仅未过期的）
func (r *tokenRepository) ListByUserID(ctx context.Context, userID uint) ([]do.RefreshToken, error) {
	var tokens []do.RefreshToken
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
		Model(&do.RefreshToken{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, err
}

// internal/repository/token/token_crud.go

func (r *tokenRepository) GetRefreshTokenTTL(ctx context.Context, jti string) (time.Duration, error) {
	var refreshToken do.RefreshToken

	err := r.db.WithContext(ctx).
		Where("jti = ? AND expires_at > ?", jti, time.Now()).
		First(&refreshToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // token 不存在或已过期
		}
		return 0, err
	}

	// 计算剩余有效期
	ttl := time.Until(refreshToken.ExpiresAt)
	return ttl, nil
}
