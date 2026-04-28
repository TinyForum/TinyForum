package auth

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tiny-forum/internal/dto"
	"tiny-forum/pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 修改 ForgotPassword 方法，使用 tokenRepo 保存重置令牌
func (s *authService) ForgotPassword(ctx context.Context, email,ip,userAgent, locale string) error {

// 1. 邮箱限流（使用 SetNX 原子操作）
    emailKey := fmt.Sprintf("forgot_password:email:%s", email)
    
    // SetNX: 只在 key 不存在时设置，返回 true 表示设置成功
    ok, err := s.redis.SetNX(ctx, emailKey, "1", 5*time.Minute).Result()
    if err != nil {
logger.Error("Failed to set rate limit", zap.String("error", err.Error()))

        return nil // 限流失败也应该返回成功，避免信息泄露
    }
    
    if !ok {
        // 限流期内，静默返回成功
		logger.Error("Forgot password rate limited",zap.String("error", err.Error()))

        return nil
    }
	// 2. IP 限流（可选，防止同一 IP 批量枚举）
    ipKey := fmt.Sprintf("forgot_password:ip:%s", ip)
    ok, err = s.redis.SetNX(ctx, ipKey, "1", 1*time.Minute).Result()
    if err == nil && !ok {
        // IP 限流，静默返回
        return nil
    }
    
    // 3. 查找用户（不泄露信息）
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 用户不存在，记录日志但不返回错误
            logger.Debug("Forgot password: email not found", zap.String("email", email))
            return nil
        }
        logger.Error("Failed to find user", zap.String("error", err.Error()),zap.String( "email", email))
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
	if err := s.tokenRepo.SaveResetToken(ctx, user.ID, token, 15*time.Minute); err != nil {
        logger.Error("Failed to save reset token", zap.String("error", err.Error()))
        return nil
    }

	// 构建应用基础URL
	appURL := s.buildAppURL()
	apiVersion := s.cfg.Basic.API.Version

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
		err = s.emailSvc.SendResetPasswordEmail(email, token, user.Username, appURL, apiVersion,ip,userAgent, locale)
		if err != nil {
			logger.Errorf("failed to send reset password email", "error", err)
		}
		
	}()

	return nil
}

// buildAppURL 构建应用基础URL
func (s *authService) buildAppURL() string {
	api := s.cfg.Basic.API
	server := s.cfg.Basic.Server

	// 构建基础URL
	protocol := api.Protocol
	if protocol == "" {
		protocol = "http"
	}

	host := api.Host
	if host == "" {
		host = server.Host
		if host == "" {
			host = "localhost"
		}
	}

	port := api.Port
	if port == 0 {
		if server.Port != 0 {
			port = server.Port
		} else {
			port = 80
		}
	}

	// 构建完整URL
	baseURL := fmt.Sprintf("%s://%s", protocol, host)

	// 优化判断逻辑：非标准端口才添加
	isStandardPort := (protocol == "http" && port == 80) || (protocol == "https" && port == 443)
	if !isStandardPort {
		baseURL = fmt.Sprintf("%s:%d", baseURL, port)
	}

	// 添加API前缀
	if api.Prefix != "" {
		baseURL = baseURL + api.Prefix
	}

	// 调试输出
	logger.Debugf("built app URL",
		"protocol", protocol,
		"host", host,
		"port", port,
		"prefix", api.Prefix,
		"full_url", baseURL,
	)

	return baseURL
}
func (s *authService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	// 查找令牌
	user, err := s.userRepo.FindByResetToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired reset token")
		}
		return err
	}

	// 检查令牌是否过期
	if user.ResetPasswordSentAt == nil {
		return errors.New("invalid reset token")
	}
	tokenExpiry := s.cfg.Private.JWT.Expire
	if time.Since(*user.ResetPasswordSentAt) > tokenExpiry {
		return errors.New("reset token has expired")
	}

	// 哈希新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码并清除令牌
	user.Password = string(hashedPassword)
	user.ResetPasswordToken = ""
	user.ResetPasswordSentAt = nil

	return s.userRepo.Update(ctx, user)
}
