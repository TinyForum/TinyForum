package user

import (
	"strconv"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=user.UserProfileResponse}
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParams(c, []response.ValidationError{
			{Field: "id", Message: "无效的用户ID格式"},
		})
		return
	}
	viewerID := getViewerID(c)
	profile, err := h.userSvc.GetUserProfile(uint(targetID), viewerID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, profile)
}

// UpdateProfile 更新个人资料
// @Summary 更新用户资料
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.UpdateProfileInput true "资料"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input model.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.UpdateProfile(userID, input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	user, _ := h.userSvc.GetProfile(userID)
	response.Success(c, user)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.ChangePasswordInput true "密码"
// @Router /users/password [patch]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	var input model.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}
	if input.OldPassword == "" {
		response.BadRequest(c, "请输入当前密码")
		return
	}
	if input.NewPassword == "" {
		response.BadRequest(c, "请输入新密码")
		return
	}
	message, err := h.userSvc.ChangePassword(userID.(uint), input.OldPassword, input.NewPassword)
	if err != nil {
		response.AppError(c, err)
		return
	}
	response.Success(c, gin.H{"message": message})
}
