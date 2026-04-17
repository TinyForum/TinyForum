package token

import (
	"context"
	"time"
	"tiny-forum/internal/model"

	"github.com/google/uuid"
)

// Create 创建 RefreshToken
func (r *TokenRepository) Create(ctx context.Context, userID uint, token string, expiresAt time.Time, userAgent, ip string) (*model.RefreshToken, error) {
	rt := &model.RefreshToken{
		UserID:    userID,
		Token:     token,
		JTI:       uuid.New().String(),
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: expiresAt,
	}

	if err := r.db.WithContext(ctx).Create(rt).Error; err != nil {
		return nil, err
	}
	return rt, nil
}

// FindByToken 根据 Token 查找（仅返回未过期的）
func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// FindByJTI 根据 JTI 查找
func (r *TokenRepository) FindByJTI(ctx context.Context, jti string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("jti = ?", jti).
		First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}
