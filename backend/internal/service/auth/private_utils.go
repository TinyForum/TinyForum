package auth

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	apperrors "tiny-forum/pkg/errors"
)

// 保留用户名黑名单（防止注册 admin、root 等特权名）
var reservedUsernames = map[string]bool{
	"admin":         true,
	"root":          true,
	"system":        true,
	"moderator":     true,
	"support":       true,
	"help":          true,
	"info":          true,
	"contact":       true,
	"superadmin":    true,
	"administrator": true,
	"staff":         true,
	"official":      true,
}

// 只允许字母、数字、下划线、连字符，长度 3-30，防止用户名包含 XSS/注入特殊字符
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)

// 验证结果

// validatePasswordStrength 统一密码强度校验（注册、修改、重置均使用）
// 规则：至少 8 位，包含字母和数字
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return apperrors.ErrPasswordTooShort
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasLetter || !hasDigit {
		return apperrors.ErrInvalidPassword.WithMessage("密码必须同时包含字母和数字")
	}
	return nil
}

// 严格邮箱校验，过滤换行符等注入字符
func isValidEmail(email string) bool {
	// 拒绝包含换行、回车、空格等危险字符（防邮件头注入）
	if strings.ContainsAny(email, "\r\n\t ") {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email) && len(email) <= 254
}

// 密码强度评级（用于提示，不强制）
func (s *authService) checkPasswordStrength(password string) string {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

	score := 0
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if hasUpper && hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}

	if score <= 2 {
		return "weak"
	}
	if score <= 4 {
		return "medium"
	}
	return "strong"
}

// 生成默认头像URL
func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
}

// 生成重置密码的 token
func (s *authService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
