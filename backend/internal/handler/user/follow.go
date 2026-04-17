package user

import (
	"strconv"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Follow 关注用户
// @Summary 关注用户
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Router /users/{id}/follow [post]
func (h *UserHandler) Follow(c *gin.Context) {
	followerID := c.GetUint("user_id")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.userSvc.Follow(followerID, uint(targetID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "关注成功"})
}

// Unfollow 取消关注
// @Summary 取消关注用户
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Router /users/{id}/follow [delete]
func (h *UserHandler) Unfollow(c *gin.Context) {
	followerID := c.GetUint("user_id")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.userSvc.Unfollow(followerID, uint(targetID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "已取消关注"})
}

// GetFollowers 获取粉丝列表
// @Summary 获取粉丝列表
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Router /users/{id}/followers [get]
func (h *UserHandler) GetFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	followers, total, err := h.userSvc.GetFollowers(uint(userID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, followers, total, page, pageSize)
}

// GetFollowing 获取关注列表
// @Summary 获取关注列表
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Router /users/{id}/following [get]
func (h *UserHandler) GetFollowing(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	following, total, err := h.userSvc.GetFollowing(uint(userID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, following, total, page, pageSize)
}
