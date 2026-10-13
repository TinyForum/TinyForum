package user

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetUserComments 获取当前用户发表的评论
// @Summary 获取我的评论
// @Description 分页获取当前登录用户发表的评论列表
// @Tags 用户中心
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /users/me/comments [get]
func (h *UserHandler) GetUserComments(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req request.GetUserCommentsRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}

	result, err := h.userSvc.GetUserComments(c.Request.Context(), req, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, result)
}
