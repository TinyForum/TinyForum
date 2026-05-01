package topic

import (
	"strconv"

	topicService "tiny-forum/internal/service/topic"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Create 创建话题
// @Summary 创建话题
// @Description 创建一个新的话题（需要管理员权限）
// @Tags 话题管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body topic.CreateTopicInput true "话题信息"
// @Success 200 {object} response.Response{data=po.Topic} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /topics [post]
func (h *TopicHandler) Create(c *gin.Context) {
	var input topicService.CreateTopicInput
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

// Update 更新话题
// @Summary 更新话题
// @Description 更新指定话题的信息（需要管理员权限）
// @Tags 话题管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Param body body topic.CreateTopicInput true "话题信息"
// @Success 200 {object} response.Response{data=po.Topic} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误或无效的话题ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "话题不存在"
// @Router /topics/{id} [put]
func (h *TopicHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	var input topicService.CreateTopicInput
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

// Delete 删除话题
// @Summary 删除话题
// @Description 删除指定话题（需要管理员权限）
// @Tags 话题管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "话题ID"
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "无效的话题ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "话题不存在"
// @Router /topics/{id} [delete]
func (h *TopicHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
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

// GetByID 获取话题详情
// @Summary 获取话题详情
// @Description 根据ID获取话题详细信息
// @Tags 话题管理
// @Produce json
// @Param id path int true "话题ID"
// @Success 200 {object} response.Response{data=po.Topic} "获取成功"
// @Failure 400 {object} response.Response "无效的话题ID"
// @Failure 404 {object} response.Response "话题不存在"
// @Router /topics/{id} [get]
func (h *TopicHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的话题ID")
		return
	}

	topic, err := h.topicSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, topic)
}

// List 获取话题列表
// @Summary 获取话题列表
// @Description 分页获取所有话题列表
// @Tags 话题管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]po.Topic}} "获取成功"
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

// GetByCreator 获取用户创建的话题
// @Summary 获取用户创建的话题
// @Description 获取指定用户创建的所有话题
// @Tags 话题管理
// @Produce json
// @Param creator_id path int true "创建者用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]po.Topic}} "获取成功"
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
