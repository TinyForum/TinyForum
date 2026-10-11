package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// forgotPasswordEmailLimit 同一邮箱 1 小时内允许的最大请求次数
	forgotPasswordEmailLimit = 300
	// forgotPasswordIPLimit 同一 IP 1 分钟内允许的最大请求次数
	forgotPasswordIPLimit = 500
	// forgotPasswordGlobalLimit 全局 1 分钟内允许的最大请求次数
	forgotPasswordGlobalLimit = 1000
)

// FIXME: 在配置 / 管理后台中进行控制
// 修改 ForgotPassword 方法，使用 tokenRepo 保存重置令牌
func (s *authService) ForgotPassword(ctx context.Context, email, ip, userAgent, locale string) error {
	// 1. 邮箱限流：1小时内最多 forgotPasswordEmailLimit 次请求
	emailKey := fmt.Sprintf("forgot_password:email:%s", email)

	emailCount, err := s.redis.Incr(ctx, emailKey).Result()
	if err != nil {
		logger.Error("Failed to increment email rate limit",
			zap.String("email", email),
			zap.Error(err))
		return nil // 返回 nil 而不是 err，避免信息泄露
	}

	// 第一次设置，添加过期时间
	if emailCount == 1 {
		if err := s.redis.Expire(ctx, emailKey, 1*time.Hour).Err(); err != nil {
			logger.Warn("Failed to set expire for email key",
				zap.String("key", emailKey),
				zap.Error(err))
		}
	}

	// 检查是否超过限制
	if emailCount > forgotPasswordEmailLimit {
		logger.Warn("Email rate limit exceeded",
			zap.String("email", email),
			zap.Int64("count", emailCount))
		return nil // 限流期内，静默返回成功
	}

	// 2. IP 限流：1分钟内最多 forgotPasswordIPLimit 次请求
	ipKey := fmt.Sprintf("forgot_password:ip:%s", ip)
	ipCount, err := s.redis.Incr(ctx, ipKey).Result()
	if err != nil {
		logger.Error("Failed to increment IP rate limit",
			zap.String("ip", ip),
			zap.Error(err))
		return nil
	}

	if ipCount == 1 {
		if err := s.redis.Expire(ctx, ipKey, 1*time.Minute).Err(); err != nil {
			logger.Warn("Failed to set expire for IP key",
				zap.String("key", ipKey),
				zap.Error(err))
		}
	}

	if ipCount > forgotPasswordIPLimit {
		logger.Warn("IP rate limit exceeded",
			zap.String("ip", ip),
			zap.Int64("count", ipCount))
		return nil // IP 限流，静默返回成功
	}

	// 3. 全局限流：整个系统每分钟最多 forgotPasswordGlobalLimit 次请求
	globalKey := "forgot_password:global"
	globalCount, err := s.redis.Incr(ctx, globalKey).Result()
	if err != nil {
		logger.Error("Failed to increment global rate limit", zap.Error(err))

	} else {
		if globalCount == 1 {
			s.redis.Expire(ctx, globalKey, 1*time.Minute)
		}
		if globalCount > forgotPasswordGlobalLimit {
			logger.Warn("Global rate limit exceeded",
				zap.Int64("count", globalCount))

			return nil
		}
	}

	// 3. 查找用户（不泄露信息）
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，记录日志但不返回错误
			logger.Debug("Forgot password: email not found", zap.String("email", email))
			return nil
		}
		logger.Error("Failed to find user", zap.String("error", err.Error()), zap.String("email", email))
		return nil // 返回 nil 而不是 err，避免信息泄露
	}

	// 4. 检查用户状态
	if user.IsBlocked {
		logger.Debug("Forgot password: user is blocked", zap.Uint("user_id", user.ID))
		return nil
	}

	if !user.IsActive {
		logger.Debug("Forgot password: user inactive", zap.Uint("user_id", user.ID))
		return nil
	}

	// 5. 生成重置令牌
	token, err := s.generateResetToken()
	if err != nil {
		logger.Error("Failed to generate reset token", zap.String("error", err.Error()))
		return nil
	}

	// 保存重置令牌（通常设置过期时间，15 分钟）
	tokenExpiresIn := 15 * time.Minute
	if err := s.tokenRepo.SaveResetToken(ctx, user.ID, token, tokenExpiresIn); err != nil {
		logger.Error("Failed to save reset token", zap.String("error", err.Error()))
		return nil
	}

	// 构建应用基础URL
	appURL := s.buildAppURL(s.cfg.Basic.Frontend.Protocol, s.cfg.Basic.Frontend.Host, s.cfg.Basic.Frontend.Port, "")
	apiVersion := s.cfg.Basic.API.Version
	// http://localhost:3000

	logger.Info("password reset requested",
		zap.String("email", email),
		zap.String("ip", ip),
		zap.String("user_agent", userAgent))

	// 7. 异步发送邮件
	go func() {
		// 使用 recover 防止 panic
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Recovered from panic in ForgotPassword goroutine", zap.Any("panic", r))
			}
		}()
		if sendErr := s.emailSvc.SendResetPasswordEmail(email, token, tokenExpiresIn, user.Username, appURL, apiVersion, ip, userAgent, locale); sendErr != nil {
			logger.Errorf("failed to send reset password email: %v", sendErr)
		}

	}()

	return nil
}

// buildAppURL 构建应用基础URL
func (s *authService) buildAppURL(protocol, host string, port int, basePath string) string {
	// 参数标准化和默认值处理
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	if protocol == "" {
		protocol = "http"
	}

	host = strings.TrimSpace(host)
	if host == "" {
		host = "localhost"
	}

	// 处理 basePath：去除首尾斜杠，但保留根路径为空字符串
	basePath = strings.TrimSpace(basePath)
	if basePath != "" {
		basePath = strings.Trim(basePath, "/")
	}

	// 构建基础 URL
	var baseURL strings.Builder

	// 添加协议和主机
	baseURL.WriteString(fmt.Sprintf("%s://%s", protocol, host))

	// 添加端口（非标准端口或明确指定时才添加）
	isStandardPort := (protocol == "http" && port == 80) || (protocol == "https" && port == 443)
	hasPort := port > 0

	if hasPort && !isStandardPort {
		baseURL.WriteString(fmt.Sprintf(":%d", port))
	}

	// 添加 basePath（如果存在）
	if basePath != "" {
		baseURL.WriteString("/" + basePath)
	}

	fullURL := baseURL.String()

	// 调试输出
	logger.Debugf("built app URL: protocol=%s host=%s port=%d base_path=%s full_url=%s",
		protocol, host, port, basePath, fullURL)

	return fullURL
}

// consumeResetTokenAndUpdatePassword 在单个事务内原子地消费重置令牌并更新密码。
// 通过行锁 + is_used=false 条件读取，确保令牌只能被消费一次，杜绝重放。
func (s *authService) consumeResetTokenAndUpdatePassword(ctx context.Context, token, hashedPassword string) error {
	return s.txManager.ExecuteInTransaction(ctx, func(tx *gorm.DB) error {
		var rt do.RefreshToken
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token = ?", token).
			Where("jti LIKE ?", "reset_%").
			Where("deleted_at IS NULL").
			Where("is_used = ?", false).
			First(&rt).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrInvalidToken
			}
			return err
		}
		if time.Now().After(rt.ExpiresAt) {
			return apperrors.ErrInvalidToken
		}

		var user do.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, rt.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.ErrUserNotFound
			}
			return err
		}
		if user.DeletedAt.Valid {
			return apperrors.ErrUserDeleted
		}
		if user.IsBlocked {
			return apperrors.ErrUserBlocked
		}

		// 标记令牌已使用（消费）
		if err := tx.Model(&do.RefreshToken{}).Where("id = ?", rt.ID).
			Update("is_used", true).Error; err != nil {
			return err
		}
		// 更新密码
		if err := tx.Model(&do.User{}).Where("id = ?", rt.UserID).
			Update("password", hashedPassword).Error; err != nil {
			return err
		}
		return nil
	})
}

// 重置密码
func (s *authService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	// 密码强度校验：与注册/修改密码保持一致，至少 8 位
	if len(req.Password) < 8 {
		return apperrors.ErrPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 原子消费令牌并更新密码，防止令牌重放
	return s.consumeResetTokenAndUpdatePassword(ctx, req.Token, string(hashedPassword))
}

// 验证密码
// ValidatePassword 验证密码是否正确
func (s *authService) ValidateOldPassword(userID uint, oldPassword string) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	return err == nil, nil
}

func (s *authService) ResetPasswordWithToken(ctx context.Context, token, newPassword string) error {
	// 1. 哈希新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		return apperrors.ErrInternalError
	}

	// 2. 原子消费令牌并更新密码
	if err := s.consumeResetTokenAndUpdatePassword(ctx, token, string(hashedPassword)); err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			// 令牌已使用/过期/无效：按幂等处理，不泄露令牌状态
			logger.Infof("Token already used or expired (idempotent response)")
			return nil
		}
		logger.Errorf("Failed to reset password with token: %v", err)
		return err
	}

	logger.Infof("Password reset successfully with token")
	return nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) (string, error) {
	// 1. 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.ErrUserNotFound
		}
		return "", fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return "", apperrors.ErrInvalidCurrentPassword
	}

	// 3. 禁止与旧密码相同
	if oldPassword == newPassword {
		return "", apperrors.ErrPasswordSameAsOld
	}

	// 4. 密码强度校验
	if err := validatePasswordStrength(newPassword); err != nil {
		return "", err
	}

	// 5. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	// 6. 更新密码
	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return "", fmt.Errorf("更新密码失败: %w", err)
	}

	// 7. 吊销所有 refresh token（强制其他设备下线）
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		// 记录日志但不影响主流程
		logger.Warnf("警告: 吊销 token 失败: %v", err)
	}

	// 8. 返回消息（根据密码强度）
	strength := s.checkPasswordStrength(newPassword)
	if strength == "weak" {
		return "密码修改成功，但密码强度较弱，建议使用更复杂的密码", nil
	}
	return "密码修改成功", nil
}
