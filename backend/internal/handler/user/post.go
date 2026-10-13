package user

import (
	"strconv"

	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetUserPosts 获取用户发布的文章
// @Summary 获取用户发布的文章
// @Description 获取当前登录用户已安装（上传）的插件列表，通常是通过 author_id 查询
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /users/me/posts [get]
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

// GetUserPostById 获取当前用户的单个帖子
// @Summary 获取我的单个帖子
// @Description 获取当前登录用户发布的单个帖子详情
// @Tags 用户中心
// @Produce json
// @Security ApiKeyAuth
// @Param post_id path int true "帖子ID"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Failure 403 {object} common.BasicResponse
// @Failure 404 {object} common.BasicResponse
// @Router /users/me/posts/{post_id} [get]
func (h *UserHandler) GetUserPostById(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}
	userID := c.GetUint("user_id")
	result, err := h.userSvc.GetUserPost(c.Request.Context(), uint(postID), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}
