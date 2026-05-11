package user

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetUserPosts 获取用户发布的文章
// @Summary 获取用户发布的文章
// @Description 获取当前登录用户已安装（上传）的插件列表，通常是通过 author_id 查询
// @Tags plugin
// @Accept json
// @Produce json
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /users/me/plugins [get]
func (h *UserHandler) GetUserPosts(c *gin.Context) {
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
