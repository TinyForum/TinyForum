package auth

import (
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
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
// @Param request body request.DeleteAccountRequest true "注销请求"
// @Success 200 {object} common.BasicResponse
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
	var input request.DeleteAccountRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果没有额外验证字段，可以忽略绑定错误
		input = request.DeleteAccountRequest{}
	}

	isDeelte, err := h.authSvc.DeleteAccount(ctx, userID.(uint), input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	Result := vo.DeleteAccountVO{
		IsDeleted: isDeelte,
	}
	logger.Infof("用户 %d 注销账户（软删除）", userID)

	response.Success(c, Result)
}

// GetDeletionStatus godoc
// @Summary 获取注销状态
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse
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
		response.HandleError(c, err)
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
// @Success 200 {object} common.BasicResponse
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
		response.HandleError(c, err)
		return
	}

	response.Success(c, common.ResponseMessage{
		Message: "已取消注销，账户已恢复",
	})
}

// ConfirmDeletion godoc
// @Summary 确认永久删除账户
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse
// @Router /auth/account/permanent [delete]
func (h *AuthHandler) ConfirmDeletion(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	var input request.DeleteAccountRequest
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

	response.Success(c, common.ResponseMessage{
		Message: "账户已永久删除",
	})
}

// ChangePassword 修改密码（登录后）
// @Summary 修改密码
// @Tags 验证管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse "无效的请求参数"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 500 {object} common.BasicResponse "服务器内部错误"
// @Router /auth/account/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest

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
