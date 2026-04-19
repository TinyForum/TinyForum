package post

import (
	"strconv"

	postRepo "tiny-forum/internal/repository/post"
	postService "tiny-forum/internal/service/post"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Create 创建帖子
// @Summary 创建帖子
// @Description 创建新的帖子（支持普通帖和问答帖）
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body post.CreatePostInput true "帖子信息"
// @Success 200 {object} response.Response{data=model.Post} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	// ctx := c.Request.Context()
	authorID := c.GetUint("user_id")

	var input postService.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	post, err := h.postSvc.Create(c, authorID, input)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, post)
}

// GetByID 获取帖子详情
// @Summary 获取帖子详情
// @Description 根据ID获取帖子的详细信息，包括点赞状态
// @Tags 帖子管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParams(c, []response.ValidationError{
			{Field: "id", Message: "无效的帖子ID格式"},
		})
		return
	}

	viewerID, exists := c.Get("user_id")
	var viewerUint uint
	if exists {
		if v, ok := viewerID.(uint); ok {
			viewerUint = v
		}
	}

	post, liked, err := h.postSvc.GetByID(uint(postID), viewerUint)
	if err != nil {
		response.Error(c, apperrors.Wrapf(apperrors.ErrPostNotFound, "ID: %d", postID))
		return
	}

	response.Success(c, gin.H{
		"post":  post,
		"liked": liked,
	})
}

// List 获取帖子列表
// @Summary 获取帖子列表
// @Description 分页获取帖子列表，支持多种筛选和排序
// @Tags 帖子管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Param sort_by query string false "排序方式" Enums(created_at, updated_at, like_count, comment_count) default(created_at)
// @Param type query string false "帖子类型" Enums(post, question)
// @Param author_id query int false "作者ID"
// @Param tag_id query int false "标签ID"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts [get]
func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	sortBy := c.Query("sort_by")
	postType := c.Query("type")

	var authorID uint
	if v := c.Query("author_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		authorID = uint(id)
	}

	var tagID uint
	if v := c.Query("tag_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		tagID = uint(id)
	}

	opts := postRepo.PostListOptions{
		AuthorID: authorID,
		TagID:    tagID,
		PostType: postType,
		Keyword:  keyword,
		SortBy:   sortBy,
	}

	posts, total, err := h.postSvc.List(page, pageSize, opts)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// Update 更新帖子
// @Summary 更新帖子
// @Description 更新自己的帖子（管理员可以更新任何帖子）
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Param body body post.UpdatePostInput true "帖子信息"
// @Success 200 {object} response.Response{data=model.Post} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /posts/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin"

	var input postService.UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	post, err := h.postSvc.Update(uint(postID), userID, isAdmin, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, post)
}

// Delete 删除帖子
// @Summary 删除帖子
// @Description 删除自己的帖子（管理员可以删除任何帖子）
// @Tags 帖子管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin"

	if err := h.postSvc.Delete(uint(postID), userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}
