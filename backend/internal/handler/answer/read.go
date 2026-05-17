package answer

import (
	"strconv"
	"tiny-forum/internal/model/do"
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
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 404 {object} common.BasicResponse"回答不存在"
// @Router /answers/{id} [get]
func (h *AnswerHandler) GetAnswer(c *gin.Context) {
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
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
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 400 {object} common.BasicResponse "无效的帖子ID"
// @Failure 404 {object} common.BasicResponse "问题不存在"
// @Router /answers/{post_id}/answers [get]
func (h *AnswerHandler) GetQuestionAnswers(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	question, answers, total, err := h.questionSvc.GetQuestionWithAnswers(uint(postID), page, pageSize)
	if err != nil {
		response.HandleError(c, err)
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
// @Success      200  {object}  common.BasicResponse  "获取成功"
// @Failure      400  {object}  common.BasicResponse "无效的回答ID"
// @Failure      401  {object}  common.BasicResponse "未授权"
// @Failure      500  {object}  common.BasicResponse "服务器内部错误"
// @Router       /answers/{id}/status [get]

// GetVoteStatus 获取用户对指定答案的投票状态及统计信息
func (h *AnswerHandler) GetVoteStatus(c *gin.Context) {
	// 1. 解析并校验 answerID
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || answerID == 0 {
		response.HandleError(c, err)
		return
	}

	// 2. 获取当前用户ID（未登录则为 0）
	userID := c.GetUint("user_id")

	// 3. 获取用户投票状态（仅登录用户）
	var userVote *do.AnswerVoteType
	if userID != 0 {
		userVote, err = h.commentSvc.GetUserVoteStatus(uint(answerID), userID)
		if err != nil {
			response.HandleError(c, err)
			return
		}
	}

	// 4. 获取投票统计（up/down 数量）
	upCount, downCount, err := h.commentSvc.GetVoteStatistics(uint(answerID))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 5. 返回结果
	response.Success(c, VoteStatusResponse{
		UserVote:  userVote,
		UpCount:   upCount,
		DownCount: downCount,
		Total:     upCount + downCount,
	})
}
