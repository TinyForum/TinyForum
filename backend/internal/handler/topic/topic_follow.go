package topic

import (
	"strconv"

	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Follow 关注话题
// @Summary 关注话题
// @Description 关注指定话题，接收话题更新通知
// @Tags 话题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Success 200 {object} vo.BasicResponse  "关注成功"
// @Failure 400 {object} vo.BasicResponse"无效的话题ID"
// @Failure 401 {object} vo.BasicResponse"未授权"
// @Failure 404 {object} vo.BasicResponse"话题不存在"
// @Router /topics/{id}/follow [post]
func (h *TopicHandler) Follow(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.topicSvc.Follow(userID, uint(topicID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "关注成功"})
}

// Unfollow 取消关注话题
// @Summary 取消关注话题
// @Description 取消关注指定话题
// @Tags 话题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Success 200 {object} vo.BasicResponse  "取消关注成功"
// @Failure 400 {object} vo.BasicResponse"无效的话题ID"
// @Failure 401 {object} vo.BasicResponse"未授权"
// @Failure 404 {object} vo.BasicResponse"话题不存在或未关注"
// @Router /topics/{id}/follow [delete]
func (h *TopicHandler) Unfollow(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.topicSvc.Unfollow(userID, uint(topicID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "取消关注成功"})
}

// IsFollowing 检查是否关注话题
// @Summary 检查是否关注话题
// @Description 检查当前用户是否已关注指定话题
// @Tags 话题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Success 200 {object} vo.BasicResponse  "获取成功"
// @Failure 400 {object} vo.BasicResponse"无效的话题ID"
// @Failure 401 {object} vo.BasicResponse"未授权"
// @Failure 500 {object} vo.BasicResponse"服务器内部错误"
// @Router /topics/{id}/follow/status [get]
func (h *TopicHandler) IsFollowing(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	userID := c.GetUint("user_id")

	isFollowing, err := h.topicSvc.IsFollowing(userID, uint(topicID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"is_following": isFollowing})
}

// GetFollowers 获取话题关注者列表
// @Summary 获取话题关注者列表
// @Description 分页获取关注指定话题的用户列表
// @Tags 话题管理
// @Produce json
// @Param id path int true "话题ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} vo.BasicResponse "获取成功"
// @Failure 400 {object} vo.BasicResponse"无效的话题ID"
// @Failure 500 {object} vo.BasicResponse"服务器内部错误"
// @Router /topics/{id}/followers [get]
func (h *TopicHandler) GetFollowers(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	followers, total, err := h.topicSvc.GetFollowers(uint(topicID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, followers, total, page, pageSize)
}
