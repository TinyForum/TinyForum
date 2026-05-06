package post

import (
	"strconv"

	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
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
// @Success 200 {object} common.BasicResponse "创建成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	// ctx := c.Request.Context()
	authorID := c.GetUint("user_id")

	var input postService.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if input.BoardID == 0 {
		response.BadRequest(c, "board_id is required")
		return
	}
	if input.Status == "" {
		input.Status = do.PostStatusPublished
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
// @Success 200 {object} common.BasicResponse  "获取成功"
// @Failure 400 {object} common.BasicResponse"无效的帖子ID"
// @Failure 404 {object} common.BasicResponse"帖子不存在"
// @Router /posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFailed(c, []response.ValidationError{
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
		response.HandleError(c, apperrors.Wrapf(apperrors.ErrPostNotFound, err, "ID: %d", postID))
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
// @Success 200 {object} common.BasicResponse  "获取成功"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /posts [get]
func (h *PostHandler) List(c *gin.Context) {

	var req request.ListPosts

	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	postType := do.ParsePostType(req.PostType)

	listPostsBO := &common.PageQuery[bo.ListPosts]{
		Page:     req.Page,
		PageSize: req.PageSize,
		Data: bo.ListPosts{
			AuthorID:         req.AuthorID,
			TagNames:         req.TagNames,
			SortBy:           req.SortBy,
			PostStatus:       do.PostStatusPublished,
			Keyword:          req.Keyword,
			Type:             postType,
			ModerationStatus: do.ModerationStatusApproved,
		},
	}

	posts, total, err := h.postSvc.List(c, listPostsBO)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, req.Page, req.PageSize)
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
// @Success 200 {object} common.BasicResponse "更新成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 404 {object} common.BasicResponse"帖子不存在"
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
// @Success 200 {object} common.BasicResponse  "删除成功"
// @Failure 400 {object} common.BasicResponse"无效的帖子ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 404 {object} common.BasicResponse"帖子不存在"
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
