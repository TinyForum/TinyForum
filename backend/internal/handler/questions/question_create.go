package question

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// CreateQuestion 创建问答帖
// @Summary 创建问答帖
// @Description 创建一个新的问答类型帖子并设置悬赏积分
// @Tags 问题管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.CreateQuestionRequest true "问答信息"
// @Success 200 {object} vo.BasicResponse "创建成功"
// @Failure 400 {object} vo.BasicResponse"请求参数错误"
// @Failure 401 {object} vo.BasicResponse"未授权"
// @Failure 403 {object} vo.BasicResponse"积分不足"
// @Router /questions/create [post]
func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("读取请求体失败: %v\n", err)
	} else {
		fmt.Printf("原始请求体: %s\n", string(body))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	var input dto.CreateQuestionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if input.BoardID == 0 {
		response.BadRequest(c, "board_id is required")
		return
	}
	if input.Status == "" {
		input.Status = do.PostStatusPublished
	}

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
