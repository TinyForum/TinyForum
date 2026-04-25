package answer

import (
	"errors"
	"strconv"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// MARK: Accept

// AcceptAnswer 采纳答案
// @Summary 采纳答案
// @Description 采纳某个回答作为问题的正确答案（仅问题作者可操作）
// @Tags 回答管理
// @Produce json
// @Security ApiKeyAuth
// @Param question_id path int true "问题帖子ID"
// @Param answer_id path int true "回答评论ID"
// @Success 200 {object} response.Response{data=object} "采纳成功"
// @Failure 400 {object} response.Response "无效的ID或操作失败"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（非问题作者）"
// @Failure 404 {object} response.Response "问题或回答不存在"
// @Router /answers/{question_id}/accept/{answer_id} [post]
func (h *AnswerHandler) AcceptAnswer(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("question_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	commentID, err := strconv.ParseUint(c.Param("answer_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.questionSvc.AcceptAnswer(uint(postID), uint(commentID), userID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrPostNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, apperrors.ErrAcceptForbidden):
			response.Forbidden(c, err.Error())
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "采纳答案成功"})
}

// UnacceptAnswer 取消接受答案
// @Summary 取消接受答案
// @Description 取消将回答标记为问题的正确答案
// @Tags 回答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回答ID"
// @Success 200 {object} response.Response{data=object} "取消成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "回答或问题不存在"
// @Router /answers/{id}/unaccept [post]
func (h *AnswerHandler) UnacceptAnswer(c *gin.Context) {
	// 1. 获取回答ID
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的回答ID")
		return
	}

	// 2. 获取当前用户
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin" || role == "moderator"

	// 3. 调用 service 层取消接受
	if err := h.commentSvc.UnacceptAnswer(uint(answerID), userID, isAdmin); err != nil {
		switch err.Error() {
		case "回答不存在":
			response.NotFound(c, err.Error())
		case "问题不存在":
			response.NotFound(c, err.Error())
		case "该回答未被接受为答案":
			response.BadRequest(c, err.Error())
		case "没有权限操作":
			response.Forbidden(c, err.Error())
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	// 4. 返回成功响应
	response.Success(c, gin.H{
		"message":   "已取消接受答案",
		"answer_id": answerID,
	})
}

// MARK: Vote

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
