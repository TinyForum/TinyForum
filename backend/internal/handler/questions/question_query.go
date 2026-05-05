package question

import (
	"strconv"

	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetQuestionsList 获取问答帖列表
// @Summary 获取问答帖列表
// @Description 分页获取所有问答类型的帖子
// @Tags 问题管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param unanswered query bool false "只看未解决"
// @Success 200 {object} vo.BasicResponse "获取成功"
// @Failure 500 {object} vo.BasicResponse"服务器内部错误"
// @Router /questions/list [get]
func (h *QuestionHandler) GetQuestionsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	unanswered := c.Query("unanswered") == "true"

	posts, total, err := h.questionSvc.GetQuestionsList(page, pageSize, unanswered)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// GetQuestionDetail 获取问答帖详情
// @Summary 获取问答帖详情
// @Description 获取指定问答帖的详细信息，包括回答列表
// @Tags 问题管理
// @Produce json
// @Param id path int true "帖子ID"
// @Param page query int false "回答页码" default(1)
// @Param page_size query int false "每页回答数" default(20)
// @Param sort query string false "排序方式" default(vote) Enums(vote,newest,oldest)
// @Success 200 {object} vo.BasicResponse  "获取成功"
// @Failure 400 {object} vo.BasicResponse"无效的帖子ID或非问答帖"
// @Failure 404 {object} vo.BasicResponse"帖子不存在"
// @Router /questions/detail/{id} [get]
func (h *QuestionHandler) GetQuestionDetail(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	logger.Infof("查询问题, questionID: %d", questionID)
	if err != nil {
		response.BadRequest(c, "无效的问题 ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sortBy := c.DefaultQuery("sort", "vote")

	question, err := h.questionSvc.GetQuestionDetail(uint(questionID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	answers, total, err := h.commentSvc.GetAnswersByPostID(uint(question.PostID), page, pageSize, sortBy)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"post":      question,
		"answers":   answers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetQuestionSimple 获取问题精简列表
// @Summary 获取问题精简列表
// @Description 获取所有问题的精简信息，可选传递 board_id 获取指定板块的问题
// @Tags 问题管理
// @Produce json
// @Param board_id query int false "板块ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param filter query string false "过滤条件" Enums(all, unanswered, answered)
// @Param sort query string false "排序方式" Enums(latest, hot, score)
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} vo.BasicResponse
// @Failure 500 {object} vo.BasicResponse
// @Router /questions/simple [get]
func (h *QuestionHandler) GetQuestionSimple(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var boardID *uint
	if boardIDStr := c.Query("board_id"); boardIDStr != "" {
		id, err := strconv.ParseUint(boardIDStr, 10, 64)
		if err == nil {
			bid := uint(id)
			boardID = &bid
		}
	}

	filter := c.Query("filter") // all, unanswered, answered
	sort := c.Query("sort")     // latest, hot, score
	keyword := c.Query("keyword")

	questions, total, err := h.questionSvc.GetQuestionSimpleList(pageSize, offset, boardID, filter, sort, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, questions, total, page, pageSize)
}
