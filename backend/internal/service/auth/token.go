package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"gorm.io/gorm"
)

// 验证 token
func (s *authService) ValidateResetToken(ctx context.Context, token string) (bool, error) {
	isVaildToken,err := s.authRepo.ValidateResetToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return isVaildToken, nil
}

func (s *authService) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *authService) GetUserEmailByResetToken(ctx context.Context, token string) (string, error)     {
	return s.authRepo.GetUserEmailByResetToken(ctx, token)
}
