package answer

import (
	"fmt"
	"strconv"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// VoteAnswer 处理对回答的投票操作
// @Summary      投票回答（支持赞同/反对）
// @Description  用户可以对指定回答进行“赞同”（up）或“反对”（down）投票。如果用户再次点击相同的投票类型，则会取消之前的投票。需要用户已登录认证。
// @Tags         回答管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      int     true  "回答ID"
// @Param        request  body      object  true  "投票类型"  example({"vote_type": "up"})
// @Success      200      {object}  object  "返回操作结果、当前赞同票数及当前用户的投票状态"
// @Failure      400      {object}  object  "请求参数错误（如无效ID、缺失或非法的投票类型）"
// @Router       /answers/{id}/vote [post]
func (h *AnswerHandler) VoteAnswer(c *gin.Context) {
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的回答ID")
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

	// 转换 vote_type: "up" -> 1, "down" -> -1
	voteValue := 1
	if input.VoteType == "down" {
		voteValue = -1
	}

	comment, err := h.commentSvc.VoteAnswer(uint(answerID), userID, voteValue)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取用户当前投票状态
	userVote, _ := h.commentSvc.GetUserVoteStatus(uint(answerID), userID)

	response.Success(c, gin.H{
		"message":    "操作成功",
		"vote_count": comment.VoteCount, // 赞同票数
		"user_vote":  userVote,          // 0:未投票, 1:赞同, -1:反对
	})
}

// RemoveVote 取消投票
// @Summary 取消回答的投票
// @Description 取消用户对指定回答的投票
// @Tags 回答管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回答ID"
// @Success 200 {object} response.Response{data=object{message=string,vote_count=int,user_vote=int}} "取消投票成功"
// @Failure 400 {object} response.Response "无效的回答ID或尚未投票"
// @Failure 401 {object} response.Response "未授权"
// @Router /answers/{id}/vote [delete]
func (h *AnswerHandler) RemoveVote(c *gin.Context) {
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的回答ID")
		return
	}

	userID := c.GetUint("user_id")

	// 调用取消投票的服务方法
	comment, err := h.commentSvc.RemoveVote(uint(answerID), userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 获取用户当前投票状态（应该是 0）
	userVote, _ := h.commentSvc.GetUserVoteStatus(uint(answerID), userID)

	response.Success(c, gin.H{
		"message":    "取消投票成功",
		"vote_count": comment.VoteCount,
		"user_vote":  userVote, // 0:未投票
	})
}

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
