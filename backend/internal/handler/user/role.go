package user

import (
	"errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetCurrentUserRoleResponse struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

// GetCurrentUserRole 获取当前登录用户的角色
// @Summary 获取当前用户角色
// @Description 从数据库查询当前登录用户的角色信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=GetCurrentUserRoleResponse} "操作成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /users/me/role [get]
func (h *UserHandler) GetCurrentUserRole(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权操作")
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		response.InternalError(c, "用户身份解析失败")
		return
	}
	role, err := h.userSvc.GetUserRoleById(userIDUint)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, "查询用户角色失败")
		return
	}
	response.Success(c, GetCurrentUserRoleResponse{
		UserID: userIDUint,
		Role:   role,
	})
}
