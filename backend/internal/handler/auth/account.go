package auth

import (
	authService "tiny-forum/internal/service/auth"
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
		response.BadRequest(c, "请求参数错误")
		return
	}

	if input.Confirm != "PERMANENT_DELETE" {
		response.BadRequest(c, "请输入 PERMANENT_DELETE 确认永久删除")
		return
	}

	err := h.authSvc.ConfirmDeletion(ctx, userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "账户已永久删除",
	})
}
