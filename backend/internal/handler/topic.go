package handler

import (
	"strconv"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topicSvc *service.TopicService
}

func NewTopicHandler(topicSvc *service.TopicService) *TopicHandler {
	return &TopicHandler{topicSvc: topicSvc}
}

// Create 创建专题
// @Summary 创建专题
// @Description 创建一个新的专题（需要管理员权限）
// @Tags 专题管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.CreateTopicInput true "专题信息"
// @Success 200 {object} response.Response{data=model.Topic} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /topics [post]
func (h *TopicHandler) Create(c *gin.Context) {
	var input service.CreateTopicInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	topic, err := h.topicSvc.Create(userID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, topic)
}

// Update 更新专题
// @Summary 更新专题
// @Description 更新指定专题的信息（需要管理员权限）
// @Tags 专题管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "专题ID"
// @Param body body service.CreateTopicInput true "专题信息"
// @Success 200 {object} response.Response{data=model.Topic} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误或无效的专题ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "专题不存在"
// @Router /topics/{id} [put]
func (h *TopicHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
		return
	}

	var input service.CreateTopicInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	topic, err := h.topicSvc.Update(uint(id), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, topic)
}

// Delete 删除专题
// @Summary 删除专题
// @Description 删除指定专题（需要管理员权限）
// @Tags 专题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "专题ID"
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "无效的专题ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "专题不存在"
// @Router /topics/{id} [delete]
func (h *TopicHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
		return
	}

	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin"

	if err := h.topicSvc.Delete(uint(id), userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// GetByID 获取专题详情
// @Summary 获取专题详情
// @Description 根据ID获取专题详细信息
// @Tags 专题管理
// @Produce json
// @Param id path int true "专题ID"
// @Success 200 {object} response.Response{data=model.Topic} "获取成功"
// @Failure 400 {object} response.Response "无效的专题ID"
// @Failure 404 {object} response.Response "专题不存在"
// @Router /topics/{id} [get]
func (h *TopicHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
		return
	}

	topic, err := h.topicSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, topic)
}

// List 获取专题列表
// @Summary 获取专题列表
// @Description 分页获取所有专题列表
// @Tags 专题管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Topic}} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /topics [get]
func (h *TopicHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	topics, total, err := h.topicSvc.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, topics, total, page, pageSize)
}

// GetByCreator 获取用户创建的专题
// @Summary 获取用户创建的专题
// @Description 获取指定用户创建的所有专题
// @Tags 专题管理
// @Produce json
// @Param creator_id path int true "创建者用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Topic}} "获取成功"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /topics/creator/{creator_id} [get]
func (h *TopicHandler) GetByCreator(c *gin.Context) {
	creatorID, err := strconv.ParseUint(c.Param("creator_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	topics, total, err := h.topicSvc.GetByCreator(uint(creatorID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, topics, total, page, pageSize)
}

// AddPost 添加帖子到专题
// @Summary 添加帖子到专题
// @Description 将指定帖子添加到专题中（需要管理员权限）
// @Tags 专题管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "专题ID"
// @Param body body AddPostToTopicRequest true "帖子信息"
// @Success 200 {object} response.Response{data=object} "添加成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "专题或帖子不存在"
// @Router /topics/{id}/posts [post]
func (h *TopicHandler) AddPost(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
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
	addInput := service.AddPostToTopicInput{
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

// RemovePost 从专题移除帖子
// @Summary 从专题移除帖子
// @Description 将指定帖子从专题中移除（需要管理员权限）
// @Tags 专题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "专题ID"
// @Param post_id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "移除成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "专题或帖子不存在"
// @Router /topics/{id}/posts/{post_id} [delete]
func (h *TopicHandler) RemovePost(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
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

// GetTopicPosts 获取专题帖子列表
// @Summary 获取专题帖子列表
// @Description 分页获取指定专题下的所有帖子
// @Tags 专题管理
// @Produce json
// @Param id path int true "专题ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 400 {object} response.Response "无效的专题ID"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /topics/{id}/posts [get]
func (h *TopicHandler) GetTopicPosts(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
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

// Follow 关注专题
// @Summary 关注专题
// @Description 关注指定专题，接收专题更新通知
// @Tags 专题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "专题ID"
// @Success 200 {object} response.Response{data=object} "关注成功"
// @Failure 400 {object} response.Response "无效的专题ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "专题不存在"
// @Router /topics/{id}/follow [post]
func (h *TopicHandler) Follow(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.topicSvc.Follow(userID, uint(topicID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "关注成功"})
}

// Unfollow 取消关注专题
// @Summary 取消关注专题
// @Description 取消关注指定专题
// @Tags 专题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "专题ID"
// @Success 200 {object} response.Response{data=object} "取消关注成功"
// @Failure 400 {object} response.Response "无效的专题ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "专题不存在或未关注"
// @Router /topics/{id}/follow [delete]
func (h *TopicHandler) Unfollow(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.topicSvc.Unfollow(userID, uint(topicID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "取消关注成功"})
}

// IsFollowing 检查是否关注专题
// @Summary 检查是否关注专题
// @Description 检查当前用户是否已关注指定专题
// @Tags 专题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "专题ID"
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的专题ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /topics/{id}/follow/status [get]
func (h *TopicHandler) IsFollowing(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
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

// GetFollowers 获取专题关注者列表
// @Summary 获取专题关注者列表
// @Description 分页获取关注指定专题的用户列表
// @Tags 专题管理
// @Produce json
// @Param id path int true "专题ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.User}} "获取成功"
// @Failure 400 {object} response.Response "无效的专题ID"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /topics/{id}/followers [get]
func (h *TopicHandler) GetFollowers(c *gin.Context) {
	topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的专题ID")
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

// AddPostToTopicRequest 添加帖子到专题请求参数
type AddPostToTopicRequest struct {
	PostID    uint `json:"post_id" binding:"required" example:"123"` // 帖子ID
	SortOrder int  `json:"sort_order" example:"0"`                   // 排序顺序
}
