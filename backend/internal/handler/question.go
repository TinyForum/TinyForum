package handler

import (
	"errors"
	"net/http"
	"strconv"

	apperrors "tiny-forum/internal/errors"
	"tiny-forum/internal/model"
	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionSvc *service.QuestionService
	commentSvc  *service.CommentService
	postSvc     *service.PostService
}

func NewQuestionHandler(
	questionSvc *service.QuestionService,
	commentSvc *service.CommentService,
	postSvc *service.PostService,
) *QuestionHandler {
	return &QuestionHandler{
		questionSvc: questionSvc,
		commentSvc:  commentSvc,
		postSvc:     postSvc,
	}
}

// GetQuestions 获取问答帖列表
// @Summary 获取问答帖列表
// @Description 分页获取所有问答类型的帖子
// @Tags 问答管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param unanswered query bool false "只看未解决"
// @Success 200 {object} response.Response{data=response.PageData} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /posts/questions [get]
func (h *QuestionHandler) GetQuestions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	unanswered := c.Query("unanswered") == "true"

	posts, total, err := h.questionSvc.GetQuestions(page, pageSize, unanswered)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// CreateQuestion 创建问答帖
// @Summary 创建问答帖
// @Description 创建一个新的问答类型帖子并设置悬赏积分
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body model.CreateQuestionInput true "问答信息"
// @Success 200 {object} response.Response{data=model.QuestionResponse} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "积分不足"
// @Router /posts/question [post]
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var input model.CreateQuestionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 从JWT token中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	question, err := h.questionSvc.CreateQuestion(userID.(uint), input)
	if err != nil {
		switch err.Error() {
		case "积分不足":
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		default:
			response.BadRequest(c, err.Error())
		}
		return
	}

	response.Success(c, question)
}

// GetQuestionDetail 获取问答帖详情
// @Summary 获取问答帖详情
// @Description 获取指定问答帖的详细信息，包括回答列表
// @Tags 问答管理
// @Produce json
// @Param id path int true "帖子ID"
// @Param page query int false "回答页码" default(1)
// @Param page_size query int false "每页回答数" default(20)
// @Param sort query string false "排序方式" default(vote) Enums(vote,newest,oldest)
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的帖子ID或非问答帖"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /posts/question/{id} [get]
func (h *QuestionHandler) GetQuestionDetail(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	// 游客也可访问，viewerID 为 0 时跳过 liked 查询
	viewerID, _ := c.Get("user_id")
	var viewerUint uint
	if v, ok := viewerID.(uint); ok {
		viewerUint = v
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort", "vote")

	post, liked, err := h.postSvc.GetByID(uint(postID), viewerUint)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	if !post.IsQuestion {
		response.BadRequest(c, "该帖子不是问答类型")
		return
	}

	answers, total, err := h.commentSvc.GetAnswersByPostID(uint(postID), page, pageSize, sortBy)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"post":      post,
		"liked":     liked,
		"answers":   answers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateAnswer 提交回答
// @Summary 提交回答
// @Description 对指定问答帖提交一个回答
// @Tags 问答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Param body body CreateAnswerRequest true "回答内容"
// @Success 200 {object} response.Response{data=model.Comment} "提交成功"
// @Failure 400 {object} response.Response "请求参数错误或非问答帖"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /posts/question/{id}/answer [post]
func (h *QuestionHandler) CreateAnswer(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	authorID := c.GetUint("user_id")
	input := service.CreateCommentInput{
		PostID:  uint(postID),
		Content: req.Content,
	}

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
// @Router /posts/questions/{post_id}/answer/{comment_id}/accept [post]
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
// @Router /posts/questions/answer/{comment_id}/vote [post]
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
}

// GetQuestionAnswers 获取问题的回答列表
// @Summary 获取问题的回答列表
// @Description 分页获取指定问题的所有回答，已采纳的回答排在最前
// @Tags 问答管理
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path int true "问题帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 404 {object} response.Response "问题不存在"
// @Router /posts/questions/{post_id}/answers [get]
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

// ---- Request types ----

// CreateAnswerRequest 提交回答的请求参数
type CreateAnswerRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"` // 回答内容
}

// VoteAnswerRequest 投票请求参数
type VoteAnswerRequest struct {
	VoteType string `json:"vote_type" binding:"required,oneof=up down" example:"up"` // up-赞同，down-反对
}
