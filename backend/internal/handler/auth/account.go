package auth

import (
	authService "tiny-forum/internal/service/auth"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// DeleteAccount godoc
// @Summary 用户注销账户（软删除）
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body auth.DeleteAccountInput true "注销请求"
// @Success 200 {object} response.Response
// @Router /auth/account [delete]
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	ctx := c.Request.Context()

	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	// 可选：验证密码或确认码
	var input authService.DeleteAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果没有额外验证字段，可以忽略绑定错误
		input = authService.DeleteAccountInput{}
	}

	err := h.authSvc.DeleteAccount(ctx, userID.(uint), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "账户已成功删除",
	})
}

// GetDeletionStatus godoc
// @Summary 获取注销状态
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/account/deletion [get]
func (h *AuthHandler) DeletionStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	// 获取用户删除状态
	status, err := h.authSvc.GetDeletionStatus(ctx, userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, status)
}

// CancelDeletion godoc
// @Summary 取消注销账户
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/account/restore [post]
func (h *AuthHandler) CancelDeletion(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	err := h.authSvc.CancelDeletion(ctx, userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "已取消注销，账户已恢复",
	})
}

// ConfirmDeletion godoc
// @Summary 确认永久删除账户
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/account/permanent [delete]
func (h *AuthHandler) ConfirmDeletion(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	var input authService.DeleteAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}

	if input.Confirm != "DELETE" {
		response.HandleError(c, apperrors.ErrInvalidConfirmation)
		return
	}

	err := h.authSvc.ConfirmDeletion(ctx, userID.(uint))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "账户已永久删除",
	})
}

// ChangePassword 修改密码（登录后）
// PUT /api/v1/auth/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8,max=32"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debugf("ChangePassword bind error: %v", err)
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}

	userID := c.GetUint("user_id")

	// 直接调用 service，所有业务逻辑都在 service 层处理
	msg, err := h.authSvc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, msg)
}
