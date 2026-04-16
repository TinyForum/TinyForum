package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tiny-forum/internal/service"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type AnswerHandler struct {
	questionSvc *service.QuestionService
	commentSvc  *service.CommentService
	postSvc     *service.PostService
}

func NewAnswerHandler(
	questionSvc *service.QuestionService,
	commentSvc *service.CommentService,
	postSvc *service.PostService,
) *AnswerHandler {
	return &AnswerHandler{
		questionSvc: questionSvc,
		commentSvc:  commentSvc,
		postSvc:     postSvc,
	}
}

// CreateAnswer 提交回答
// @Summary 提交回答
// @Description 对指定问答帖提交一个回答
// @Tags 回答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Param body body CreateAnswerRequest true "回答内容"
// @Success 200 {object} response.Response{data=model.Comment} "提交成功"
// @Failure 400 {object} response.Response "请求参数错误或非问答帖"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /answers/{id}/answer [post]
func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	// 1. 获取问题ID（URL 参数）
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的问题 ID")
		return
	}

	// 2. 通过 questionID 查询 Question 记录，获取 postID
	question, err := h.questionSvc.GetQuestionByID(uint(questionID))
	if err != nil {
		response.BadRequest(c, "问题不存在")
		return
	}

	// 3. 获取 postID
	postID := question.PostID

	// 4. 绑定请求参数
	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 5. 获取当前用户 ID
	authorID := c.GetUint("user_id")

	// 6. 构建输入
	input := service.CreateCommentInput{
		PostID:  postID,
		Content: req.Content,
	}

	// 7. 创建回答
	comment, err := h.commentSvc.CreateAnswer(authorID, input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, comment)
}

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
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		case errors.Is(err, apperrors.ErrAcceptForbidden):
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "采纳答案成功"})
}

// // VoteAnswer 投票回答
// // @Summary 投票回答
// // @Description 对问题的回答进行投票（赞同或反对）
// // @Tags 回答管理
// // @Accept json
// // @Produce json
// // @Security ApiKeyAuth
// // @Param id path int true "回答评论ID"
// // @Param body body VoteAnswerRequest true "投票信息"
// // @Success 200 {object} response.Response{data=object} "投票成功"
// // @Failure 400 {object} response.Response "无效的评论ID或投票类型"
// // @Failure 401 {object} response.Response "未授权"
// // @Failure 404 {object} response.Response "回答不存在"
// // @Router /answers/answer/{id}/vote [post]
// func (h *AnswerHandler) VoteAnswer(c *gin.Context) {
// 	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		response.BadRequest(c, "无效的评论ID")
// 		return
// 	}

// 	var input struct {
// 		VoteType string `json:"vote_type" binding:"required,oneof=up down"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		response.BadRequest(c, err.Error())
// 		return
// 	}

// 	userID := c.GetUint("user_id")
// 	voteInput := service.VoteAnswerInput{
// 		CommentID: uint(commentID),
// 		VoteType:  input.VoteType,
// 	}

// 	voteResult, err := h.questionSvc.VoteAnswer(userID, voteInput)
// 	if err != nil {
// 		response.BadRequest(c, err.Error())
// 		return
// 	}
// 	response.Success(c, gin.H{"message": "投票成功", "result": voteResult})
// }

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

// GetAnswer 获取回答详情
// @Summary 获取回答详情
// @Description 根据ID获取回答的详细信息
// @Tags 问答管理
// @Accept json
// @Produce json
// @Param id path int true "回答ID"
// @Success 200 {object} response.Response{data=model.Comment} "获取成功"
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

// DeleteAnswer 删除回答
// @Summary 删除回答
// @Description 删除指定的回答（作者本人、问题作者、版主或管理员可操作）
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "回答ID"
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "回答不存在"
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

type VoteStatusResponse struct {
	UserVote  int `json:"user_vote"`
	UpCount   int `json:"up_count"`
	DownCount int `json:"down_count"`
	// Total     int `json:"total"`
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

// =============
// ---- Request types ----

// CreateAnswerRequest 提交回答的请求参数
type CreateAnswerRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"` // 回答内容
}

// VoteAnswerRequest 投票请求参数
type VoteAnswerRequest struct {
	VoteType string `json:"vote_type" binding:"required,oneof=up down" example:"up"` // up-赞同，down-反对
}
