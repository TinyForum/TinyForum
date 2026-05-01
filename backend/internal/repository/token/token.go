package token

import (
	"context"
	"errors"
	"time"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

// ---- Password Reset Token ----

// SaveResetToken 保存重置密码 token

// FindUserByResetToken 通过 reset token 查找对应 user ID（自动过滤已过期、已使用）
func (r *tokenRepository) FindUserByResetToken(ctx context.Context, token string) (uint, error) {
	var resetToken po.RefreshToken
	err := r.db.WithContext(ctx).
		Where("token = ? AND expires_at > ? AND used = ?", token, time.Now(), false).
		First(&resetToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // token 无效，不返回错误（避免暴露信息）
		}
		return 0, err
	}

	return resetToken.UserID, nil
}

// DeleteResetToken 使用后立即删除 reset token（单次使用保证）
func (r *tokenRepository) DeleteResetToken(ctx context.Context, token string) error {
	// 方案1: 物理删除
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&po.RefreshToken{}).Error

	// 方案2: 标记为已使用（推荐，便于审计）
	// return r.db.WithContext(ctx).
	// 	Model(&model.PasswordResetToken{}).
	// 	Where("token = ?", token).
	// 	Updates(map[string]interface{}{
	// 		"used":      true,
	// 		"used_at":   time.Now(),
	// 		"updated_at": time.Now(),
	// 	}).Error
}

// MarkResetTokenAsUsed 标记 token 为已使用（软删除，防止重放）
func (r *tokenRepository) MarkResetTokenAsUsed(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&po.RefreshToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"used":       true,
			"used_at":    time.Now(),
			"updated_at": time.Now(),
		}).Error
}

// ---- JWT Token Blacklist ----

// RevokeToken 将 JWT 的 JTI 加入黑名单（用于注销）
func (r *tokenRepository) RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error {
	if jti == "" {
		return nil
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil // 已过期，无需加入黑名单
	}

	// 存入 Redis，TTL 等于 JWT 剩余有效期
	return r.redis.Set(ctx, "revoked:"+jti, 1, ttl).Err()
}

// IsTokenRevoked 检查 JTI 是否已被吊销
func (r *tokenRepository) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}

	// 先查 Redis 黑名单
	exists, err := r.redis.Exists(ctx, "revoked:"+jti).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	}

	// 可选：再查数据库（如果 Redis 没有，且需要持久化黑名单）
	// var revokedToken model.RevokedToken
	// err = r.db.WithContext(ctx).
	//     Where("jti = ?", jti).
	//     First(&revokedToken).Error
	// if err == nil {
	//     return true, nil
	// }

	return false, nil
}

// ---- Session Management ----

// RevokeAllUserTokens 吊销某用户的所有 token（改密/删号后全端下线）
func (r *tokenRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	// 1. 获取用户所有有效的 refresh token
	var tokens []po.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Find(&tokens).Error

	if err != nil {
		return err
	}

	// 2. 将所有 token 的 JTI 加入 Redis 黑名单
	pipe := r.redis.Pipeline()
	for _, token := range tokens {
		ttl := time.Until(token.ExpiresAt)
		if ttl > 0 {
			pipe.Set(ctx, "revoked:"+token.JTI, 1, ttl)
		}
	}

	// 3. 批量执行 Redis 命令
	if _, err := pipe.Exec(ctx); err != nil {
		// 记录错误但不中断，继续删除数据库记录
		_ = err
	}

	// 4. 删除数据库中的记录（或标记为已撤销）
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&po.RefreshToken{}).Error
}

// GetUserActiveTokens 获取用户所有活跃 token（用于设备管理）
func (r *tokenRepository) GetUserActiveTokens(ctx context.Context, userID uint) ([]po.RefreshToken, error) {
	var tokens []po.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&tokens).Error

	return tokens, err
}

// RevokeTokenByJTI 根据 JTI 吊销单个 token
func (r *tokenRepository) RevokeTokenByJTI(ctx context.Context, jti string) error {
	// 先查询 token 信息
	var token po.RefreshToken
	err := r.db.WithContext(ctx).
		Where("jti = ?", jti).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	// 加入 Redis 黑名单
	ttl := time.Until(token.ExpiresAt)
	if ttl > 0 {
		if err := r.redis.Set(ctx, "revoked:"+jti, 1, ttl).Err(); err != nil {
			return err
		}
	}

	// 删除数据库记录
	return r.db.WithContext(ctx).
		Where("jti = ?", jti).
		Delete(&po.RefreshToken{}).Error
}

// CleanExpiredBlacklist 清理过期的黑名单记录（可选，Redis 会自动过期）
func (r *tokenRepository) CleanExpiredBlacklist(ctx context.Context) error {
	// Redis 的 TTL 机制会自动清理，这个方法可以空实现
	// 或者手动扫描并删除（通常不需要）
	return nil
}
