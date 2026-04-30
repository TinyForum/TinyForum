package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tiny-forum/internal/dto"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// FIXME: 在配置 / 管理后台中进行控制
// 修改 ForgotPassword 方法，使用 tokenRepo 保存重置令牌
func (s *authService) ForgotPassword(ctx context.Context, email, ip, userAgent, locale string) error {
	// 1. 邮箱限流：1小时内最多3次请求
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
	if emailCount > 300 {
		logger.Warn("Email rate limit exceeded",
			zap.String("email", email),
			zap.Int64("count", emailCount))
		return nil // 限流期内，静默返回成功
	}

	// 2. IP 限流：1分钟内最多5次请求
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

	if ipCount > 500 {
		logger.Warn("IP rate limit exceeded",
			zap.String("ip", ip),
			zap.Int64("count", ipCount))
		return nil // IP 限流，静默返回成功
	}

	// 3. 全局限流：整个系统每分钟最多100次请求
	globalKey := "forgot_password:global"
	globalCount, err := s.redis.Incr(ctx, globalKey).Result()
	if err != nil {
		logger.Error("Failed to increment global rate limit", zap.Error(err))

	} else {
		if globalCount == 1 {
			s.redis.Expire(ctx, globalKey, 1*time.Minute)
		}
		if globalCount > 1000 {
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
		err = s.emailSvc.SendResetPasswordEmail(email, token, tokenExpiresIn, user.Username, appURL, apiVersion, ip, userAgent, locale)
		if err != nil {
			logger.Errorf("failed to send reset password email", "error", err)
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
	logger.Debugf("built app URL",
		"protocol", protocol,
		"host", host,
		"port", port,
		"base_path", basePath,
		"full_url", fullURL,
	)

	return fullURL
}

// FIXME: 应该从 token 表查
// 重置密码
func (s *authService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	// 查找令牌
	user, err := s.authRepo.GetUserByResetToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired reset token")
		}
		return err
	}

	// 哈希新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码并清除令牌
	user.Password = string(hashedPassword)

	return s.userRepo.Update(ctx, user)
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
	logger.Infof("Service: ResetPasswordWithToken called with token: %s", token)

	// 1. 验证 token 并获取用户
	user, err := s.authRepo.GetUserByResetToken(ctx, token)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			logger.Infof("Token already used or expired (idempotent response): %s", token)
			return nil
		}
		// 其他错误才记录 Error
		logger.Errorf("Failed to get user by token: %v", err)
		return err
	}

	// 2. 检查用户状态
	if user.DeletedAt.Valid {
		logger.Warnf("User deleted: %d", user.ID)
		return apperrors.ErrUserDeleted
	}
	if user.IsBlocked {
		logger.Warnf("User blocked: %d", user.ID)
		return apperrors.ErrUserBlocked
	}

	// 3. 哈希新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Failed to hash password: %v", err)
		return apperrors.ErrInternalError
	}

	// 4. 更新用户密码
	user.Password = string(hashedPassword)
	if err := s.authRepo.Update(ctx, user); err != nil {
		logger.Errorf("Failed to update user: %v", err)
		return apperrors.ErrInternalError
	}

	// 5. 删除已使用的 reset token
	if err := s.authRepo.DeleteResetToken(ctx, token); err != nil {
		// 删除失败不影响主流程，改为 Warn 级别
		logger.Warnf("Failed to delete reset token (may already be deleted): %v", err)
		// 不返回错误
	}

	logger.Infof("Password reset successfully for user: %s", user.Email)
	return nil
}
