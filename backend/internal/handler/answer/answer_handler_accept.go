package answer

import (
	"errors"
	"strconv"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

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
