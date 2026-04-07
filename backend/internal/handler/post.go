package handler

import (
	"strconv"

	"bbs-forum/internal/repository"
	"bbs-forum/internal/service"
	"bbs-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postSvc *service.PostService
}

func NewPostHandler(postSvc *service.PostService) *PostHandler {
	return &PostHandler{postSvc: postSvc}
}

func (h *PostHandler) Create(c *gin.Context) {
	authorID := c.GetUint("user_id")
	var input service.CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	post, err := h.postSvc.Create(authorID, input)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, post)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	viewerID, _ := c.Get("user_id")
	var viewerUint uint
	if v, ok := viewerID.(uint); ok {
		viewerUint = v
	}

	post, liked, err := h.postSvc.GetByID(uint(postID), viewerUint)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, gin.H{"post": post, "liked": liked})
}

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

	opts := repository.PostListOptions{
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

func (h *PostHandler) Update(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin"

	var input service.UpdatePostInput
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

func (h *PostHandler) Like(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	userID := c.GetUint("user_id")
	if err := h.postSvc.Like(userID, uint(postID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "点赞成功"})
}

func (h *PostHandler) Unlike(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	userID := c.GetUint("user_id")
	if err := h.postSvc.Unlike(userID, uint(postID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "已取消点赞"})
}

// Admin handlers
func (h *PostHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	posts, total, err := h.postSvc.AdminList(page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

func (h *PostHandler) AdminTogglePin(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	if err := h.postSvc.TogglePin(uint(postID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}
