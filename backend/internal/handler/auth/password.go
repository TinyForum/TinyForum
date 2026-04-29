// internal/handler/auth/password.go
package auth

import (
	"tiny-forum/internal/dto"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
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
		response.HandleError(ctx, apperrors.ErrInvalidRequest)
		return
	}
	locale := parseAcceptLanguage(ctx.GetHeader("Accept-Language"))
	if locale == "" {
		locale = "en"
	}
	// 忽略错误，保护用户隐私
	_ = c.authSvc.ForgotPassword(ctx.Request.Context(), req.Email, ctx.ClientIP(), ctx.Request.UserAgent(), locale)
	message := getUnifiedMessage(locale)
	response.Success(ctx, &dto.ForgotPasswordResponse{
		Message: message,
	})
}

// ResetPassword 登录用户重置密码
// ConfirmDeletion godoc
// @Summary 重置密码
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/password/reset [put]
func (c *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.HandleError(ctx, apperrors.ErrInternalError)
		return
	}

	if err := c.authSvc.ResetPassword(ctx, &req); err != nil {
		response.HandleError(ctx, apperrors.ErrInternalError)
		return
	}

	response.Success(ctx, gin.H{
		"message": "Password has been reset successfully",
	})
}

// internal/handler/auth/password.go

// ResetPasswordWithToken 通过 token 重置密码（用户未登录状态）
// @Summary 通过token重置密码
// @Description 使用邮件中的token重置密码
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordWithTokenRequest true "重置密码请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/password/reset [put]
func (h *AuthHandler) ResetPasswordWithToken(ctx *gin.Context) {
	logger.Info("=== ResetPasswordWithToken called ===")

	var req dto.ResetPasswordWithTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Failed to bind request: %v", err)
		response.BadRequest(ctx, "Invalid request format")
		return
	}

	// 验证密码
	if len(req.Password) < 6 {
		response.BadRequest(ctx, "Password must be at least 6 characters")
		return
	}

	// 调用服务层重置密码
	err := h.authSvc.ResetPasswordWithToken(ctx.Request.Context(), req.Token, req.Password)
	if err != nil {
		logger.Errorf("Reset password failed: %v", err)
		response.HandleError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"message": "Password has been reset successfully",
		"success": true,
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
func (h *AuthHandler) ValidateResetToken(ctx *gin.Context) {
	logger.Info("=== [handler] request validate-token ===")

	token := ctx.Query("token")
	if token == "" {
		response.BadRequest(ctx, "token parameter is required")
		return
	}

	valid, err := h.authSvc.ValidateResetToken(ctx.Request.Context(), token)
	if err != nil {
		logger.Errorf("validate reset token failed: %v", err)
		response.InternalError(ctx, "failed to validate token")
		return
	}

	if !valid {
		response.BadRequest(ctx, "token is invalid or has expired")
		return
	}

	// 验证成功
	response.Success(ctx, gin.H{
		"valid": valid,
	})
}

// // TODO: 修复重置密码页面和更改密码页面

// // ShowResetPage 显示重置密码页面
// // ValidateResetToken 验证重置密码 token
// // ConfirmDeletion godoc
// // @Summary 重置密码
// // @Tags 验证管理
// // @Accept json
// // @Produce json
// // @Success 200 {object} response.Response
// // @Router /auth/password/reset [get]
// func (h *AuthHandler) ShowResetPage(ctx *gin.Context) {
// 	logger.Info("=== [handler] Request ResetPage ===")
// token := ctx.Query("token")
// 	if token == "" {
// 		response.Error(ctx, apperrors.ErrValidationFailed)
// 		return
// 	}

// 	valid, err := h.authSvc.ValidateResetToken(ctx.Request.Context(), token)
// 	if err != nil {
// 		response.Error(ctx, apperrors.ErrValidationFailed)
// 		return
// 	}
// 	if !valid {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
//             "message": "重置链接已过期或无效，请重新申请",
//             "action":  "reset_password",
//             "link":    "/auth/password/forgot",
//         })
// 		return
// 	}
// 	userEmail, err := h.authSvc.GetUserEmailByResetToken(ctx,token)

// 	if err != nil {
// 		response.BadRequest(ctx, "验证时出现错误")
// 		return
// 	}

// 	ctx.HTML(http.StatusOK, "reset_password.html", gin.H{
// 		"token":   token,
// 		"email":   maskEmail(userEmail), // 显示部分邮箱，如 u***r@example.com
// 		"expires": "10分钟",
// 	})
// }
