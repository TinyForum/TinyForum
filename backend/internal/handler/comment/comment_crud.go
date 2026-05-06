package comment

import (
	"strconv"

	commentService "tiny-forum/internal/service/comment"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Create 创建评论
// @Summary 创建评论
// @Description 创建一条新的评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body comment.CreateCommentInput true "评论信息"
// @Success 200 {object} common.BasicResponse "创建成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Router /comments [post]
func (h *CommentHandler) Create(c *gin.Context) {
	authorID := c.GetUint("user_id")
	var input commentService.CreateCommentInput
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

// List 获取评论列表
// @Summary 获取帖子的评论列表
// @Description 分页获取指定帖子的所有评论
// @Tags 评论管理
// @Produce json
// @Param post_id path int true "帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 400 {object} common.BasicResponse"无效的帖子ID"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /comments/post/{post_id} [get]
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

// Delete 删除评论
// @Summary 删除评论
// @Description 删除指定的评论（用户可以删除自己的评论，管理员可以删除任何评论）
// @Tags 评论管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} common.BasicResponse  "删除成功"
// @Failure 400 {object} common.BasicResponse"无效的评论ID或删除失败"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 404 {object} common.BasicResponse"评论不存在"
// @Router /comments/{id} [delete]
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
