package topic

import (
	"strconv"

	topicService "tiny-forum/internal/service/topic"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AddPost 添加帖子到话题
// @Summary 添加帖子到话题
// @Description 将指定帖子添加到话题中（需要管理员权限）
// @Tags 话题管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Param body body AddPostToTopicRequest true "帖子信息"
// @Success 200 {object} response.Response{data=object} "添加成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "话题或帖子不存在"
// @Router /topics/{id}/posts [post]
func (h *TopicHandler) AddPost(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	var input struct {
		PostID    uint `json:"post_id" binding:"required"`
		SortOrder int  `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	addInput := topicService.AddPostToTopicInput{
		TopicID:   uint(topicID),
		PostID:    input.PostID,
		SortOrder: input.SortOrder,
	}

	if err := h.topicSvc.AddPostToTopic(addInput, userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "添加帖子成功"})
}

// RemovePost 从话题移除帖子
// @Summary 从话题移除帖子
// @Description 将指定帖子从话题中移除（需要管理员权限）
// @Tags 话题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Param post_id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "移除成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "话题或帖子不存在"
// @Router /topics/{id}/posts/{post_id} [delete]
func (h *TopicHandler) RemovePost(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.topicSvc.RemovePostFromTopic(uint(topicID), uint(postID), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "移除帖子成功"})
}

// GetTopicPosts 获取话题帖子列表
// @Summary 获取话题帖子列表
// @Description 分页获取指定话题下的所有帖子
// @Tags 话题管理
// @Produce json
// @Param id path int true "话题ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 400 {object} response.Response "无效的话题ID"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /topics/{id}/posts [get]
func (h *TopicHandler) GetTopicPosts(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := h.topicSvc.GetTopicPosts(uint(topicID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// AddPostToTopicRequest 添加帖子到话题请求参数
type AddPostToTopicRequest struct {
	PostID    uint `json:"post_id" binding:"required" example:"123"` // 帖子ID
	SortOrder int  `json:"sort_order" example:"0"`                   // 排序顺序
}
