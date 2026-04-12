package handler

import (
	"strconv"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentSvc  *service.CommentService
	questionSvc *service.QuestionService
}

func NewCommentHandler(commentSvc *service.CommentService, questionSvc *service.QuestionService) *CommentHandler {
	return &CommentHandler{
		commentSvc:  commentSvc,
		questionSvc: questionSvc,
	}
}

// Create 创建评论
// @Summary 创建评论
// @Description 创建一条新的评论
// @Tags 评论管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.CreateCommentInput true "评论信息"
// @Success 200 {object} response.Response{data=model.Comment} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /comments [post]
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

// List 获取评论列表
// @Summary 获取帖子的评论列表
// @Description 分页获取指定帖子的所有评论
// @Tags 评论管理
// @Produce json
// @Param post_id path int true "帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Comment}} "获取成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 500 {object} response.Response "服务器内部错误"
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
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "无效的评论ID或删除失败"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "评论不存在"
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

// ========== 新增的投票相关方法 ==========

// VoteAnswer 对答案进行投票（赞成/反对）
// @Summary 对答案投票
// @Description 对问答帖的答案进行投票（赞成up/反对down），重复投票会取消
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param body body object true "投票类型" example({"vote_type":"up"})
// @Success 200 {object} response.Response{data=object} "投票成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "不能给自己的答案投票"
// @Router /comments/{id}/vote [post]
func (h *CommentHandler) VoteAnswer(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	var input struct {
		VoteType string `json:"vote_type" binding:"required,oneof=up down"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	voteInput := service.VoteAnswerInput{
		CommentID: uint(commentID),
		VoteType:  input.VoteType,
	}

	result, err := h.questionSvc.VoteAnswer(userID, voteInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetAnswerVoteStatus 获取当前用户对答案的投票状态
// @Summary 获取答案投票状态
// @Description 获取当前用户对指定答案的投票状态
// @Tags 问答管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的评论ID"
// @Router /comments/{id}/vote [get]
func (h *CommentHandler) GetAnswerVoteStatus(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	userID := c.GetUint("user_id")
	status, err := h.questionSvc.GetAnswerVoteStatus(userID, uint(commentID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, status)
}

// MarkAsAnswer 标记/取消标记为答案（版主或作者）
// @Summary 标记为答案
// @Description 将评论标记为问题的答案（帖子作者或版主可操作）
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param body body object true "是否标记为答案" example({"is_answer":true})
// @Success 200 {object} response.Response{data=object} "操作成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /comments/{id}/answer [put]
func (h *CommentHandler) MarkAsAnswer(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	var input struct {
		IsAnswer bool `json:"is_answer" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin"

	if err := h.commentSvc.MarkAsAnswer(uint(commentID), userID, isAdmin, input.IsAnswer); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}

// GetAnswers 获取帖子的所有答案
// @Summary 获取帖子的答案列表
// @Description 获取指定帖子的所有答案（仅限问答帖）
// @Tags 问答管理
// @Produce json
// @Param post_id path int true "帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort query string false "排序方式" default(vote) Enums(vote, newest, oldest)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Comment}} "获取成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Router /comments/post/{post_id}/answers [get]
// func (h *CommentHandler) GetAnswers(c *gin.Context) {
// 	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
// 	if err != nil {
// 		response.BadRequest(c, "无效的帖子ID")
// 		return
// 	}

// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
// 	sortBy := c.DefaultQuery("sort", "vote") // vote, newest, oldest

// 	answers, total, err := h.commentSvc.GetAnswersByPostID(uint(postID), page, pageSize, sortBy)
// 	if err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}
// 	response.SuccessPage(c, answers, total, page, pageSize)
// }

// AcceptAnswer 采纳答案（保留在PostHandler中，这里也提供一个便捷接口）
// @Summary 采纳答案
// @Description 采纳某个回答作为最佳答案（仅帖子作者可操作）
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param post_id query int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "采纳成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /comments/{id}/accept [post]
func (h *CommentHandler) AcceptAnswer(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	postIDStr := c.Query("post_id")
	if postIDStr == "" {
		response.BadRequest(c, "缺少帖子ID")
		return
	}

	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.questionSvc.AcceptAnswer(uint(postID), uint(commentID), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "采纳答案成功"})
}
