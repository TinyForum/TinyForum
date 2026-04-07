package handler

import (
	"strconv"

	"bbs-forum/internal/service"
	"bbs-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentSvc *service.CommentService
}

func NewCommentHandler(commentSvc *service.CommentService) *CommentHandler {
	return &CommentHandler{commentSvc: commentSvc}
}

func (h *CommentHandler) Create(c *gin.Context) {
	authorID := c.GetUint("user_id")
	var input service.CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	comment, err := h.commentSvc.Create(authorID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, comment)
}

func (h *CommentHandler) List(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	comments, total, err := h.commentSvc.List(uint(postID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, comments, total, page, pageSize)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin"

	if err := h.commentSvc.Delete(uint(commentID), userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}
