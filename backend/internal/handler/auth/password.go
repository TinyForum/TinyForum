// internal/handler/auth/password.go
package auth

import (
	"tiny-forum/internal/dto"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

func (c *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, apperrors.ErrInternalError)
		return
	}

	locale := ctx.GetHeader("Accept-Language")
	if locale == "" {
		locale = "en"
	}

	// 修改调用方式，匹配 service 签名
	if err := c.authSvc.ForgotPassword(ctx.Request.Context(), req.Email, locale); err != nil {
		response.Success(ctx, gin.H{
			"message": "If your email is registered, you will receive a password reset link",
		})
		return
	}

	response.Success(ctx, gin.H{
		"message": "If your email is registered, you will receive a password reset link",
	})
}

func (c *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, apperrors.ErrInternalError)
		return
	}

	// 修改调用方式
	if err := c.authSvc.ResetPassword(ctx, &req); err != nil {
		response.Error(ctx, apperrors.ErrInternalError)
		return
	}

	response.Success(ctx, gin.H{
		"message": "Password has been reset successfully",
	})
}

func (c *AuthHandler) ValidateResetToken(ctx *gin.Context) {
	token := ctx.Query("token")

	if token == "" {
		response.Error(ctx, apperrors.ErrValidationFailed)
		return
	}

	// 修改调用方式，传递 context 而不是 gin.Context
	valid, err := c.authSvc.ValidateResetToken(ctx.Request.Context(), token)
	if err != nil {
		response.Error(ctx, apperrors.ErrValidationFailed)
		return
	}

	response.Success(ctx, gin.H{
		"valid": valid,
	})
}
