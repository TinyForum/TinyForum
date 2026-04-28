package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"
)

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
