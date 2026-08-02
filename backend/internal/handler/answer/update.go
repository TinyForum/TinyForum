package answer

import (
	"errors"
	"strconv"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
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
// @Param id path int true "回答ID"
// @Param question_id path int true "问题帖子ID"
// @Success 200 {object} common.BasicResponse "采纳成功"
// @Failure 400 {object} common.BasicResponse"无效的ID或操作失败"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（非问题作者）"
// @Failure 404 {object} common.BasicResponse"问题或回答不存在"
// @Router /answers/{id}/accept/{question_id} [post]
func (h *AnswerHandler) AcceptAnswer(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("question_id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	userID := c.GetUint("user_id")

	if err := h.questionSvc.AcceptAnswer(c.Request.Context(), uint(questionID), uint(answerID), userID); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrPostNotFound):
			response.HandleError(c, err)
		case errors.Is(err, apperrors.ErrAcceptForbidden):
			response.HandleError(c, err)
		default:
			response.HandleError(c, err)
		}
		return
	}
	responseData := vo.AcceptAnswerVO{
		Message: vo.ActionAcceptedAnswer,
	}

	response.Success(c, responseData)
}

// UnacceptAnswer 取消接受答案
// @Summary 取消接受答案
// @Description 取消将回答标记为问题的正确答案
// @Tags 回答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回答ID"
// @Success 200 {object} common.BasicResponse "取消成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 404 {object} common.BasicResponse"回答或问题不存在"
// @Router /answers/{id}/unaccept [post]
func (h *AnswerHandler) UnacceptAnswer(c *gin.Context) {
	// 1. 获取回答ID
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
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
			response.HandleError(c, err)
		case "问题不存在":
			response.HandleError(c, err)
		case "该回答未被接受为答案":
			response.HandleError(c, err)
		case "没有权限操作":
			response.HandleError(c, err)
		default:
			response.HandleError(c, err)
		}
		return
	}

	responseData := common.ResponseMessage{
		Message: "已取消接受答案",
	}
	response.Success(c, responseData)

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
	// 1. 解析并校验 answerID
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || answerID == 0 {
		response.HandleError(c, err)
		return
	}

	// 2. 绑定并校验请求体
	var input struct {
		VoteType string `json:"vote_type" binding:"required,oneof=up down"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.HandleError(c, err)
		return
	}

	// 3. 获取当前用户ID（必须已登录）
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.HandleError(c, err)
		return
	}

	// 4. 转换 vote_type 为 do.AnswerVoteType 枚举
	var voteType do.AnswerVoteType
	switch input.VoteType {
	case "up":
		voteType = do.AnswerVoteTypeUp
	case "down":
		voteType = do.AnswerVoteTypeDown
	default:
		response.HandleError(c, err)
		return
	}

	// 5. 调用 Service 层投票
	comment, err := h.commentSvc.VoteAnswer(uint(answerID), userID, voteType)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 6. 获取用户最新的投票状态（可能因取消投票变为 nil）
	userVote, _ := h.commentSvc.GetUserVoteStatus(uint(answerID), userID)

	// 7. 返回响应
	responseData := vo.VoteResponseVO{
		Message:   "操作成功",
		VoteCount: comment.DownVotes + comment.UpVotes,
		UserVote:  userVote,
	}
	response.Success(c, responseData)

}
