package user

import (
	"strconv"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetProfile 获取用户资料（支持用户名或数字ID）
// @Summary 获取用户资料
// @Tags 用户管理
// @Produce json
// @Param id path string true "用户名或用户ID"
// @Success 200 {object} common.BasicResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	idParam := c.Param("id")
	viewerID := getViewerID(c)

	if targetID, err := strconv.ParseUint(idParam, 10, 64); err == nil {
		profile, err := h.userSvc.GetUserProfile(uint(targetID), viewerID)
		if err != nil {
			response.HandleError(c, err)
			return
		}
		response.Success(c, profile)
		return
	}

	profile, err := h.userSvc.GetUserProfileByUsername(idParam, viewerID)
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
// @Router /users/me/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input do.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.HandleError(c, err)
		return
	}
	if err := h.userSvc.UpdateProfile(userID, input); err != nil {
		response.HandleError(c, err)
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
// @Success 200 {object} common.BasicResponse
// @Router /auth/me [get]
// Deprecated: 无路由引用，当前用户信息由前端 auth 状态提供
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userSvc.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}
