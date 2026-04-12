package handler

import (
	"strconv"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionSvc *service.QuestionService
}

func NewQuestionHandler(questionSvc *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionSvc: questionSvc}
}

// AcceptAnswer 采纳答案
// @Summary 采纳答案
// @Description 采纳某个回答作为问题的正确答案（仅问题作者可操作）
// @Tags 问答管理
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path int true "问题帖子ID"
// @Param comment_id path int true "回答评论ID"
// @Success 200 {object} response.Response{data=object} "采纳成功"
// @Failure 400 {object} response.Response "无效的ID或操作失败"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（非问题作者）"
// @Failure 404 {object} response.Response "问题或回答不存在"
// @Router /questions/{post_id}/answer/{comment_id}/accept [post]
func (h *QuestionHandler) AcceptAnswer(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.questionSvc.AcceptAnswer(uint(postID), uint(commentID), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "采纳答案成功"})
}

// VoteAnswer 投票回答
// @Summary 投票回答
// @Description 对问题的回答进行投票（赞同或反对）
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param comment_id path int true "回答评论ID"
// @Param body body VoteAnswerRequest true "投票信息"
// @Success 200 {object} response.Response{data=object} "投票成功"
// @Failure 400 {object} response.Response "无效的评论ID或投票类型"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "回答不存在"
// @Router /questions/answer/{comment_id}/vote [post]
func (h *QuestionHandler) VoteAnswer(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
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

	voteResult, err := h.questionSvc.VoteAnswer(userID, voteInput)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "投票成功", "result": voteResult})

	response.Success(c, gin.H{"message": "投票成功"})
}

// GetQuestionAnswers 获取问题的回答列表
// @Summary 获取问题的回答列表
// @Description 分页获取指定问题的所有回答，已采纳的回答会排在前面
// @Tags 问答管理
// @Produce json
// @Param post_id path int true "问题帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 404 {object} response.Response "问题不存在"
// @Router /questions/{post_id}/answers [get]
func (h *QuestionHandler) GetQuestionAnswers(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	question, answers, total, err := h.questionSvc.GetQuestionWithAnswers(uint(postID), page, pageSize)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"question":  question,
		"answers":   answers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// VoteAnswerRequest 投票请求参数
type VoteAnswerRequest struct {
	VoteType string `json:"vote_type" binding:"required,oneof=up down" example:"up" Enums:"up,down"` // 投票类型：up-赞同，down-反对
}
