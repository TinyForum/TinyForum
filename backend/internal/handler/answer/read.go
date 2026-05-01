package answer

import (
	"fmt"
	"strconv"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// MARK: Answer

// GetAnswer 获取回答详情
// @Summary 获取回答详情
// @Description 根据ID获取回答的详细信息
// @Tags 问答管理
// @Accept json
// @Produce json
// @Param id path int true "回答ID"
// @Success 200 {object} response.Response{data=po.Comment} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "回答不存在"
// @Router /answers/{id} [get]
func (h *AnswerHandler) GetAnswer(c *gin.Context) {
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的回答ID")
		return
	}

	answer, err := h.commentSvc.GetAnswerByID(uint(answerID))
	if err != nil {
		response.NotFound(c, apperrors.ErrAnswerNotFound.Error())
		return
	}

	response.Success(c, answer)
}

// GetQuestionAnswers 获取问题的回答列表
// @Summary 获取问题的回答列表
// @Description 分页获取指定问题的所有回答，已采纳的回答排在最前
// @Tags 回答管理
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path int true "问题帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 404 {object} response.Response "问题不存在"
// @Router /answers/{post_id}/answers [get]
func (h *AnswerHandler) GetQuestionAnswers(c *gin.Context) {
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

// MARK: Vote

// GetVoteStatus 获取回答的投票状态
// @Summary      获取回答投票状态
// @Description  获取指定回答的投票统计信息（赞同数、反对数、总票数）以及当前用户的投票状态
// @Tags         回答管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "回答ID"
// @Success      200  {object}  response.Response{data=VoteStatusResponse}  "获取成功"
// @Failure      400  {object}  response.Response  "无效的回答ID"
// @Failure      401  {object}  response.Response  "未授权"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /answers/{id}/status [get]
func (h *AnswerHandler) GetVoteStatus(c *gin.Context) {
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的回答ID")
		return
	}

	userID := c.GetUint("user_id")

	userVote, err := h.commentSvc.GetUserVoteStatus(uint(answerID), userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	upCount, downCount, err := h.commentSvc.GetVoteStatistics(uint(answerID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	fmt.Printf("user_vote: %+v\n", userVote)
	fmt.Printf("up_count: %d, down_count: %d, total: %d\n", upCount, downCount, upCount+downCount)

	stats := VoteStatusResponse{
		UserVote:  userVote,
		UpCount:   upCount,
		DownCount: downCount,
		// Total:     upCount + downCount,
	}

	response.Success(c, stats)
}
