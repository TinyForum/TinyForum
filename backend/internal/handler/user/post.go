package user

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// handler/user_handler.go
func (h *UserHandler) GetUserPosts(c *gin.Context) {
	// 获取当前登录用户 ID（例如从 JWT 中间件）
	userID := c.GetUint("user_id")

	var req request.GetUserPostsRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}

	result, err := h.userSvc.GetUserPosts(c.Request.Context(), req, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, result)
}
