package token

import (
	"context"
	"fmt"
	"time"
	"tiny-forum/internal/model/do"

	"github.com/google/uuid"
)

// Create 创建 RefreshToken
func (r *tokenRepository) Create(ctx context.Context, userID uint, token string, expiresAt time.Time, userAgent, ip string) (*do.RefreshToken, error) {
	rt := &do.RefreshToken{
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
func (r *tokenRepository) FindByToken(ctx context.Context, token string) (*do.RefreshToken, error) {
	var rt do.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// FindByJTI 根据 JTI 查找
func (r *tokenRepository) FindByJTI(ctx context.Context, jti string) (*do.RefreshToken, error) {
	var rt do.RefreshToken
	err := r.db.WithContext(ctx).
		Where("jti = ?", jti).
		First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// FindByUserID 查找用户的所有有效令牌
func (r *tokenRepository) FindByUserID(ctx context.Context, userID uint) ([]*do.RefreshToken, error) {
	var tokens []*do.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Find(&tokens).Error
	return tokens, err
}

// Revoke 撤销令牌（软删除或标记过期）
func (r *tokenRepository) Revoke(ctx context.Context, tokenID uint) error {
	return r.db.WithContext(ctx).
		Model(&do.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("expires_at", time.Now()).Error
}

// RevokeAllByUserID 撤销用户的所有令牌
func (r *tokenRepository) RevokeAllByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&do.RefreshToken{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Update("expires_at", time.Now()).Error
}

// RevokeByJTI 根据 JTI 撤销令牌
func (r *tokenRepository) RevokeByJTI(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).
		Model(&do.RefreshToken{}).
		Where("jti = ? AND expires_at > ?", jti, time.Now()).
		Update("expires_at", time.Now()).Error
}

// CleanExpired 清理过期令牌（物理删除）
func (r *tokenRepository) CleanExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at <= ?", time.Now()).
		Delete(&do.RefreshToken{}).Error
}

func (r *tokenRepository) SaveResetToken(ctx context.Context, userID uint, token string, expiration time.Duration) error {
	// 生成 JTI（用于唯一标识）
	jti := generateJTI()

	// 创建重置令牌（复用 RefreshToken 模型）
	refreshToken := &do.RefreshToken{
		UserID:    userID,
		Token:     token,
		JTI:       jti,
		ExpiresAt: time.Now().Add(expiration),
	}

	return r.db.WithContext(ctx).Create(refreshToken).Error
}

// generateJTI 生成唯一标识
func generateJTI() string {
	return fmt.Sprintf("reset_%d_%s", time.Now().UnixNano(), uuid.New().String()[:8])
}
