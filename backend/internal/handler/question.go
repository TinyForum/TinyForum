// package handler

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"strconv"

// 	"tiny-forum/internal/model"
// 	"tiny-forum/internal/service"
// 	"tiny-forum/pkg/logger"
// 	"tiny-forum/pkg/response"

// 	"github.com/gin-gonic/gin"
// )

// type QuestionHandler struct {
// 	questionSvc *service.QuestionService
// 	commentSvc  *service.CommentService
// 	postSvc     *service.PostService
// }

// func NewQuestionHandler(
// 	questionSvc *service.QuestionService,
// 	commentSvc *service.CommentService,
// 	postSvc *service.PostService,
// ) *QuestionHandler {
// 	return &QuestionHandler{
// 		questionSvc: questionSvc,
// 		commentSvc:  commentSvc,
// 		postSvc:     postSvc,
// 	}
// }

// // GetQuestions 获取问答帖列表
// // @Summary 获取问答帖列表
// // @Description 分页获取所有问答类型的帖子
// // @Tags 问题管理
// // @Produce json
// // @Param page query int false "页码" default(1)
// // @Param page_size query int false "每页数量" default(20)
// // @Param unanswered query bool false "只看未解决"
// // @Success 200 {object} response.Response{data=response.PageData} "获取成功"
// // @Failure 500 {object} response.Response "服务器内部错误"
// // @Router /questions/list [get]
// func (h *QuestionHandler) GetQuestionsList(c *gin.Context) {
// 	// 获取问题列表
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
// 	unanswered := c.Query("unanswered") == "true"

// 	posts, total, err := h.questionSvc.GetQuestionsList(page, pageSize, unanswered)
// 	if err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}
// 	response.SuccessPage(c, posts, total, page, pageSize)
// }

// // CreateQuestion 创建问答帖
// // @Summary 创建问答帖
// // @Description 创建一个新的问答类型帖子并设置悬赏积分
// // @Tags 问题管理
// // @Accept json
// // @Produce json
// // @Security ApiKeyAuth
// // @Param body body model.CreateQuestionInput true "问答信息"
// // @Success 200 {object} response.Response{data=model.QuestionResponse} "创建成功"
// // @Failure 400 {object} response.Response "请求参数错误"
// // @Failure 401 {object} response.Response "未授权"
// // @Failure 403 {object} response.Response "积分不足"
// // @Router /questions/create [post]
// func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		fmt.Printf("读取请求体失败: %v\n", err)
// 	} else {
// 		fmt.Printf("原始请求体: %s\n", string(body))
// 		// 关键：重新设置请求体，否则后续绑定会失败
// 		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
// 	}

// 	var input model.CreateQuestionInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		response.BadRequest(c, "参数错误: "+err.Error())
// 		return
// 	}

// 	// 从JWT token中获取用户ID
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		response.Unauthorized(c, "请先登录")
// 		return
// 	}

// 	question, err := h.questionSvc.CreateQuestion(userID.(uint), input)
// 	if err != nil {
// 		switch err.Error() {
// 		case "积分不足":
// 			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
// 		default:
// 			response.BadRequest(c, err.Error())
// 		}
// 		return
// 	}

// 	response.Success(c, question)
// }

// // GetQuestionDetail 获取问答帖详情
// // @Summary 获取问答帖详情
// // @Description 获取指定问答帖的详细信息，包括回答列表
// // @Tags 问题管理
// // @Produce json
// // @Param id path int true "帖子ID"
// // @Param page query int false "回答页码" default(1)
// // @Param page_size query int false "每页回答数" default(20)
// // @Param sort query string false "排序方式" default(vote) Enums(vote,newest,oldest)
// // @Success 200 {object} response.Response{data=object} "获取成功"
// // @Failure 400 {object} response.Response "无效的帖子ID或非问答帖"
// // @Failure 404 {object} response.Response "帖子不存在"
// // @Router /questions/detail/{id} [get]
// func (h *QuestionHandler) GetQuestionDetail(c *gin.Context) {
// 	// 先查询 question 表 是否有这个 id
// 	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	logger.Infof("查询问题, questionID: %d", questionID)
// 	if err != nil {
// 		response.BadRequest(c, "无效的问题 ID")
// 		return
// 	}

// 	// 游客也可访问，viewerID 为 0 时跳过 liked 查询
// 	// viewerID, _ := c.Get("user_id")
// 	// var viewerUint uint
// 	// if v, ok := viewerID.(uint); ok {
// 	// 	viewerUint = v
// 	// }

// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
// 	sortBy := c.DefaultQuery("sort", "vote")

// 	question, err := h.questionSvc.GetQuestionDetail(uint(questionID))
// 	if err != nil {
// 		response.NotFound(c, err.Error())
// 		return
// 	}
// 	// fmt.Printf("读取 post: %v\n", post.Type)
// 	// if post.Type != "question" {
// 	// 	response.BadRequest(c, "该帖子不是问答类型")
// 	// 	return
// 	// }

// 	answers, total, err := h.commentSvc.GetAnswersByPostID(uint(question.PostID), page, pageSize, sortBy)
// 	if err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}

// 	response.Success(c, gin.H{
// 		"post": question,
// 		// "liked":     liked,
// 		"answers":   answers,
// 		"total":     total,
// 		"page":      page,
// 		"page_size": pageSize,
// 	})
// }

// // GetQuestionSimple 获取问题精简列表
// // @Summary 获取问题精简列表
// // @Description 获取所有问题的精简信息，可选传递 board_id 获取指定板块的问题
// // @Tags 问题管理
// // @Produce json
// // @Param board_id query int false "板块ID"
// // @Param page query int false "页码" default(1)
// // @Param page_size query int false "每页数量" default(20)
// // @Success 200 {object} response.Response{data=response.PageData{list=[]model.QuestionListResponse}}
// // @Failure 500 {object} response.Response
// // @Router /questions/simple [get]
// func (h *QuestionHandler) GetQuestionSimple(c *gin.Context) {
// 	// 获取分页参数
// 	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 || pageSize > 100 {
// 		pageSize = 20
// 	}

// 	offset := (page - 1) * pageSize

// 	// 获取板块过滤参数
// 	var boardID *uint
// 	if boardIDStr := c.Query("board_id"); boardIDStr != "" {
// 		id, err := strconv.ParseUint(boardIDStr, 10, 64)
// 		if err == nil {
// 			bid := uint(id)
// 			boardID = &bid
// 		}
// 	}

// 	// 获取筛选和排序参数
// 	filter := c.Query("filter")   // all, unanswered, answered
// 	sort := c.Query("sort")       // latest, hot, score
// 	keyword := c.Query("keyword") // 关键词搜索

// 	// 调用 Service 层
// 	questions, total, err := h.questionSvc.GetQuestionSimpleList(pageSize, offset, boardID, filter, sort, keyword)
// 	if err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}

//		response.SuccessPage(c, questions, total, page, pageSize)
//	}
package handler
