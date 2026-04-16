// repository/token.go
package repository

import (
	"context"
	"time"
	"tiny-forum/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

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

// DeleteByUserID 删除用户的所有 Token（全局登出、密码重置、封禁时调用）
func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.RefreshToken{}).Error
}

// DeleteByToken 删除单个 Token（单设备登出）
func (r *TokenRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&model.RefreshToken{}).Error
}

// DeleteByJTI 根据 JTI 精确撤销
func (r *TokenRepository) DeleteByJTI(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).
		Where("jti = ?", jti).
		Delete(&model.RefreshToken{}).Error
}

// DeleteExpired 清理过期的 Token（建议用定时任务每天执行）
func (r *TokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at <= ?", time.Now()).
		Delete(&model.RefreshToken{})
	return result.RowsAffected, result.Error
}

// ListByUserID 获取用户的所有登录设备
func (r *TokenRepository) ListByUserID(ctx context.Context, userID uint) ([]model.RefreshToken, error) {
	var tokens []model.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&tokens).Error
	return tokens, err
}

// CountByUserID 统计用户的有效 Token 数量
func (r *TokenRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, err
}

// DeleteByUserIDExcept 删除用户的所有 Token，但保留指定的一个（切换设备时用）
func (r *TokenRepository) DeleteByUserIDExcept(ctx context.Context, userID uint, keepToken string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND token != ?", userID, keepToken).
		Delete(&model.RefreshToken{}).Error
}
