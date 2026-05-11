package answer

import (
	"strconv"
	commentService "tiny-forum/internal/service/comment"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateAnswer 提交回答
// @Summary 提交回答
// @Description 对指定问答帖提交一个回答
// @Tags 回答管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Param body body CreateAnswerRequest true "回答内容"
// @Success 200 {object} common.BasicResponse "提交成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误或非问答帖"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 404 {object} common.BasicResponse"帖子不存在"
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
	input := commentService.CreateCommentInput{
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
