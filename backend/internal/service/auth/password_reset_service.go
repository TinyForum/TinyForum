// internal/auth/password_reset_service.go
package auth

import (
	"context"
	"tiny-forum/internal/dto"
)

// PasswordResetService 定义密码重置服务接口
type PasswordResetService interface {
	ForgotPassword(ctx context.Context, email, locale string) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
	ValidateResetToken(ctx context.Context, token string) (bool, error)
}
