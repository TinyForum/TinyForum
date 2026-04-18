// internal/auth/password_reset_service_impl.go
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"tiny-forum/internal/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 修改 ForgotPassword 方法，使用 tokenRepo 保存重置令牌
func (s *authService) ForgotPassword(ctx context.Context, email, locale string) error {
	// 查找用户
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // 用户不存在时返回成功（安全考虑）
		}
		return err
	}

	// 检查用户状态
	if user.IsBlocked || !user.IsActive {
		return nil // 用户被锁定或未激活时不发送邮件（安全考虑）
	}

	// 生成重置令牌
	token, err := s.generateResetToken()
	if err != nil {
		return err
	}

	// 保存重置令牌（通常设置过期时间，如1小时）
	if err := s.tokenRepo.SaveResetToken(ctx, user.ID, token, 1*time.Hour); err != nil {
		return err
	}

	// 构建应用基础URL
	appURL := s.buildAppURL()

	// 异步发送邮件
	go func() {
		// 使用 context.Background() 避免请求上下文被取消
		_ = s.emailSvc.SendResetPasswordEmail(email, token, user.Username, locale, appURL)
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
		host = "localhost"
	}

	port := api.Port
	if port == 0 {
		port = server.Port
	}

	// 构建完整URL
	baseURL := fmt.Sprintf("%s://%s", protocol, host)

	// 非标准端口才添加端口号
	if (protocol == "http" && port != 80) || (protocol == "https" && port != 443) {
		baseURL = fmt.Sprintf("%s:%d", baseURL, port)
	}

	// 添加API前缀
	if api.Prefix != "" {
		baseURL = baseURL + api.Prefix
	}

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

func (s *authService) ValidateResetToken(ctx context.Context, token string) (bool, error) {
	user, err := s.userRepo.FindByResetToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if user.ResetPasswordSentAt == nil {
		return false, nil
	}
	tokenExpiry := s.cfg.Private.JWT.Expire
	if time.Since(*user.ResetPasswordSentAt) > tokenExpiry {
		return false, nil
	}

	return true, nil
}

func (s *authService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
