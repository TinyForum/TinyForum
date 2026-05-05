package user

import (
	"strconv"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} vo.BasicResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {

		response.ValidationFailed(c, []response.ValidationError{
			{Field: "id", Message: "无效的用户ID格式"},
		})
		return
	}
	viewerID := getViewerID(c)
	profile, err := h.userSvc.GetUserProfile(uint(targetID), viewerID)
	if err != nil {
		response.HandleError(c, err)
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
// @Param body body do.UpdateProfileInput true "资料"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input do.UpdateProfileInput
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

// Me godoc
// @Summary 获取当前用户信息
// @Tags 验证管理
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} vo.BasicResponse
// @Router /auth/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userSvc.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}
