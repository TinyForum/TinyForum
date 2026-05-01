package user

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"time"
	"tiny-forum/internal/model/po"
	apperrors "tiny-forum/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	tempPasswordLength = 12
	tempPasswordTTL    = 30 * time.Minute
)

// ResetUserPasswordWithTemp 管理员重置用户密码（生成临时密码）
func (s *userService) ResetUserPasswordWithTemp(operatorID uint, targetID uint) (string, error) {
	ctx := context.Background()

	targetUser, err := s.repo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.ErrUserNotFound.WithMessagef("ID: %d", targetID)
		}
		return "", fmt.Errorf("查询目标用户失败: %w", err)
	}

	operator, err := s.repo.FindByID(operatorID)
	if err != nil {
		return "", fmt.Errorf("查询操作者信息失败: %w", err)
	}

	if targetUser.Role == po.RoleSuperAdmin && operator.Role != po.RoleSuperAdmin {
		return "", apperrors.ErrInsufficientPermission.WithMessage("只有超级管理员才能重置其他超级管理员的密码")
	}
	if operator.Role != po.RoleSuperAdmin && targetUser.Role == po.RoleAdmin {
		return "", apperrors.ErrInsufficientPermission.WithMessage("只有超级管理员才能重置其他管理员的密码")
	}

	tempPassword, err := generateSecurePassword(tempPasswordLength)
	if err != nil {
		return "", fmt.Errorf("生成临时密码失败: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	if err := s.repo.UpdatePassword(ctx, targetID, string(hashedPassword)); err != nil {
		return "", fmt.Errorf("更新密码失败: %w", err)
	}

	expiresAt := time.Now().Add(tempPasswordTTL)
	_ = s.repo.SetTempPasswordFlag(ctx, targetID, true, expiresAt)

	s.logAudit(ctx, operatorID, targetID, "reset_password_temp",
		fmt.Sprintf("管理员 %d 为用户 %d (%s) 生成了临时密码", operatorID, targetID, targetUser.Username))

	return tempPassword, nil
}

// ResetUserPassword 管理员重置用户密码（指定新密码）
func (s *userService) ResetUserPassword(operatorID uint, targetID uint, newPassword string) error {
	ctx := context.Background()

	targetUser, err := s.repo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound.WithMessagef("ID: %d", targetID)
		}
		return fmt.Errorf("查询目标用户失败: %w", err)
	}

	operator, err := s.repo.FindByID(operatorID)
	if err != nil {
		return fmt.Errorf("查询操作者信息失败: %w", err)
	}

	if targetUser.Role == po.RoleSuperAdmin && operator.Role != po.RoleSuperAdmin {
		return apperrors.ErrInsufficientPermission.WithMessage("只有超级管理员才能重置其他超级管理员的密码")
	}
	if operator.Role != po.RoleSuperAdmin && targetUser.Role == po.RoleAdmin {
		return apperrors.ErrInsufficientPermission.WithMessage("只有超级管理员才能重置其他管理员的密码")
	}

	if err := s.validatePasswordStrength(newPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	if err := s.repo.UpdatePassword(ctx, targetID, string(hashedPassword)); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	s.logAudit(ctx, operatorID, targetID, "reset_password",
		fmt.Sprintf("管理员 %d 重置了用户 %d (%s) 的密码", operatorID, targetID, targetUser.Username))
	return nil
}

// private

// generateSecurePassword 生成安全随机密码
func generateSecurePassword(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}
	return string(result), nil
}

// validatePasswordStrength 密码强度校验
func (s *userService) validatePasswordStrength(password string) error {
	if len(password) < 6 {
		return apperrors.ErrInvalidPassword.WithMessagef("密码长度至少 6 位")
	}
	if len(password) > 32 {
		return apperrors.ErrInvalidPassword.WithMessagef("密码长度不能超过 32 位")
	}
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasDigit || !hasLetter {
		return apperrors.ErrInvalidPassword.WithMessagef("密码必须包含数字和字母")
	}
	return nil
}
