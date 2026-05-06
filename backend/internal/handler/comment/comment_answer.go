package comment

import (
	"strconv"

	questionService "tiny-forum/internal/service/question"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// VoteAnswer 对答案进行投票（赞成/反对）
// @Summary 对答案投票
// @Description 对问答帖的答案进行投票（赞成up/反对down），重复投票会取消
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param body body object true "投票类型" example({"vote_type":"up"})
// @Success 200 {object} common.BasicResponse  "投票成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"不能给自己的答案投票"
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
	voteInput := questionService.VoteAnswerInput{
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
// @Success 200 {object} common.BasicResponse  "获取成功"
// @Failure 400 {object} common.BasicResponse"无效的评论ID"
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

// AcceptAnswer 采纳答案（便捷接口）
// @Summary 采纳答案
// @Description 采纳某个回答作为最佳答案（仅帖子作者可操作）
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Param post_id query int true "帖子ID"
// @Success 200 {object} common.BasicResponse  "采纳成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
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
