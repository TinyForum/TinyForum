package user

import (
	"errors"
	"strconv"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetScore 查询用户积分（用户查自己，管理员可查任意用户）
// @Summary 查询积分
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Router /users/score [get]
func (h *UserHandler) GetScore(c *gin.Context) {
	viewerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}
	viewerUint, ok := viewerID.(uint)
	if !ok {
		response.BadRequest(c, "无效的用户身份信息")
		return
	}
	viewerRole, _ := c.Get("role")
	viewerRoleStr, _ := viewerRole.(string)

	var targetID uint
	idParam := c.Param("id")
	if idParam != "" {
		userID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的用户ID")
			return
		}
		targetID = uint(userID)
		if viewerUint != targetID && viewerRoleStr != "admin" && viewerRoleStr != "super_admin" {
			response.Forbidden(c, "权限不足，只能查询自己的积分")
			return
		}
	} else {
		targetID = viewerUint
	}

	score, err := h.userSvc.GetScoreById(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, "查询积分失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{
		"score":   score,
		"user_id": targetID,
	})
}
