package answer

import (
	"strconv"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// DeleteAnswer 删除回答
// @Summary 删除回答
// @Description 删除指定的回答（作者本人、问题作者、版主或管理员可操作）
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回答ID"
// @Success 200 {object} vo.BasicResponse "删除成功"
// @Failure 400 {object} vo.BasicResponse"请求参数错误"
// @Failure 401 {object} vo.BasicResponse"未授权"
// @Failure 403 {object} vo.BasicResponse"无权限"
// @Failure 404 {object} vo.BasicResponse"回答不存在"
// @Router /answers/{id} [delete]
func (h *AnswerHandler) DeleteAnswer(c *gin.Context) {
	// 1. 获取并验证回答ID
	answerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的回答ID")
		return
	}

	// 2. 获取当前用户信息
	userID := c.GetUint("user_id")
	role, exists := c.Get("user_role")
	if !exists {
		response.Unauthorized(c, "未获取到用户信息")
		return
	}

	// 3. 获取用户角色（用于权限判断）
	userRole, ok := role.(string)
	if !ok {
		userRole = "user"
	}
	isAdmin := userRole == "admin" || userRole == "moderator"

	// 4. 调用服务层删除回答
	if err := h.commentSvc.DeleteAnswer(uint(answerID), userID, isAdmin); err != nil {
		// 根据错误类型返回不同的响应
		if err.Error() == "回答不存在" {
			response.NotFound(c, err.Error())
			return
		}
		if err.Error() == "没有权限删除此回答" {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	// 5. 返回成功响应
	response.Success(c, gin.H{
		"message":   "删除成功",
		"answer_id": answerID,
	})
}

// MARK: Vote
// RemoveVote 取消投票
// @Summary 取消回答的投票
// @Description 取消用户对指定回答的投票
// @Tags 回答管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回答ID"
// @Success 200 {object} vo.BasicResponse "取消投票成功"
// @Failure 400 {object} vo.BasicResponse"无效的回答ID或尚未投票"
// @Failure 401 {object} vo.BasicResponse"未授权"
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
