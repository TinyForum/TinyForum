// internal/auth/repository/repository.go
package auth

import (
	"context"
	"errors"
	"time"
	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
)

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindByResetToken 根据 token 查找有效的重置密码记录
func (r *authRepository) FindByResetToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	logger.Info("=== FindByResetToken ===")
	if token == "" {
		return nil, errors.New("token is empty")
	}

	var refreshToken model.RefreshToken

	// 添加安全验证条件
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		Where("expires_at > ?", time.Now()). // 检查是否过期
		Where("deleted_at IS NULL").         // 软删除检查
		Where("jti LIKE ?", "reset_%").      // 只查询密码重置类型的 token
		First(&refreshToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired reset token")
		}
		return nil, err
	}
	if refreshToken.IsUsed {
		return nil, errors.New("token already used")
	}

	return &refreshToken, nil
}

// GetUserByResetToken 直接返回用户信息
// internal/repository/auth/user_repository.go

func (r *authRepository) GetUserByResetToken(ctx context.Context, token string) (*model.User, error) {
	logger.Infof("=== GetUserByResetToken ===")
	logger.Infof("Token: %s", token)

	if token == "" {
		return nil, errors.New("token is empty")
	}

	// 先查询 refresh_token
	var resetToken model.RefreshToken

	// 查询 token
	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		Where("deleted_at IS NULL").
		Where("expires_at > ?", time.Now()).
		Where("jti LIKE ?", "reset_%").
		First(&resetToken).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrTokenInvalid
		}
		logger.Errorf("Database error: %v", err)
		return nil, err
	}

	logger.Infof("Found reset token: UserID=%d, ExpiresAt=%v", resetToken.UserID, resetToken.ExpiresAt)

	// 然后查询用户
	var user model.User
	err = r.db.WithContext(ctx).
		Where("id = ?", resetToken.UserID).
		Where("deleted_at IS NULL").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("User not found: %d", resetToken.UserID)
			return nil, apperrors.ErrUserNotFound
		}
		logger.Errorf("Database error: %v", err)
		return nil, err
	}

	logger.Infof("Found user: ID=%d, Email=%s", user.ID, user.Email)
	return &user, nil
}
func (r *authRepository) GetUserEmailByResetToken(ctx context.Context, token string) (string, error) {
	logger.Infof("=== GetUserEmailByResetToken: %s ===", token)

	if token == "" {
		return "", errors.New("token is empty")
	}

	var email string

	// 联表查询，只获取 email
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("INNER JOIN refresh_tokens ON refresh_tokens.user_id = users.id").
		Where("refresh_tokens.token = ?", token).
		Where("refresh_tokens.expires_at > ?", time.Now()).
		Where("refresh_tokens.deleted_at IS NULL").
		Where("refresh_tokens.jti LIKE ?", "reset_%").
		Where("users.deleted_at IS NULL").
		Where("users.is_blocked = ?", false).
		Select("users.email").
		Scan(&email).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid or expired reset token")
		}
		return "", err
	}

	if email == "" {
		return "", errors.New("user not found")
	}

	logger.Infof("Found email: %s", email)
	return email, nil
}

// MarkTokenAsUsed 标记 token 为已使用（防止重复使用）
func (r *authRepository) MarkTokenAsUsed(ctx context.Context, tokenID uint) error {
	logger.Info("=== MarkTokenAsUsed ===")
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("id = ?", tokenID).
		Updates(map[string]interface{}{
			"used":       true,
			"used_at":    now,
			"updated_at": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}

// DeleteExpiredResetTokens 清理过期的重置 token（可定期执行）
func (r *authRepository) DeleteExpiredResetTokens(ctx context.Context) error {
	// 软删除过期的重置 token
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Where("jti LIKE ?", "reset_%").
		Delete(&model.RefreshToken{})

	if result.Error != nil {
		return result.Error
	}

	// 记录清理数量（可选）
	logger.Debugf("Cleaned up %d expired reset tokens", result.RowsAffected)

	return nil
}

func (r *authRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) Save(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
