// internal/handler/auth/password.go
package auth

import (
	// "log"
	"tiny-forum/internal/dto"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ForgotPassword 忘记密码
// @Summary 忘记密码
// @Description 发送密码重置链接到用户邮箱（出于安全考虑，无论邮箱是否存在都返回成功）
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "邮箱信息"
// @Success 200 {object} response.Response{data=object,message=string} "成功（邮箱存在与否均返回此消息）"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /auth/password/forgot [post]
func (c *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, apperrors.ErrInvalidRequest)
		return
	}

	locale := parseAcceptLanguage(ctx.GetHeader("Accept-Language"))
	if locale == "" {
		locale = "en"
	}

	_ = c.authSvc.ForgotPassword(ctx.Request.Context(), req.Email, ctx.ClientIP(),ctx.Request.UserAgent(),locale)

	message := getUnifiedMessage(locale)
	response.Success(ctx, &dto.ForgotPasswordResponse{
		Message: message,
	})
}

// 辅助函数
func getUnifiedMessage(locale string) string {
	if locale == "zh-CN" || locale == "zh" {
		return "如果您的邮箱已注册，您将收到密码重置链接"
	}
	return "If your email is registered, you will receive a password reset link"
}

// ResetPassword 重置密码
// ConfirmDeletion godoc
// @Summary 重置密码
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/password/reset [post]
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

// ValidateResetToken 验证重置密码 token
// ConfirmDeletion godoc
// @Summary 忘记密码
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/password/validate-token [get]
func (c *AuthHandler) ValidateResetToken(ctx *gin.Context) {
	token := ctx.Query("token")

	if token == "" {
		response.Error(ctx, apperrors.ErrValidationFailed)
		return
	}

	valid, err := c.authSvc.ValidateResetToken(ctx.Request.Context(), token)
	if err != nil {
		response.Error(ctx, apperrors.ErrValidationFailed)
		return
	}

	response.Success(ctx, gin.H{
		"valid": valid,
	})
}
