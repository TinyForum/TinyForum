package article

import (
	"strconv"

	"tiny-forum/internal/model/common"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Like 点赞帖子
// @Summary 点赞帖子
// @Description 为指定帖子点赞
// @Tags 帖子管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} common.BasicResponse  "点赞成功"
// @Failure 400 {object} common.BasicResponse"无效的帖子ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /posts/{id}/like [post]
func (h *ArticleHandler) Like(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	userID := c.GetUint("user_id")
	if err := h.articleSvc.Like(userID, uint(postID)); err != nil {
		response.HandleError(c, err)
		return
	}
	responseData := common.ResponseMessage{
		Message: "点赞成功",
	}
	response.Success(c, responseData)
}

// Unlike 取消点赞帖子
// @Summary 取消点赞帖子
// @Description 取消对指定帖子的点赞
// @Tags 帖子管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} common.BasicResponse  "取消点赞成功"
// @Failure 400 {object} common.BasicResponse"无效的帖子ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /posts/{id}/like [delete]
func (h *ArticleHandler) Unlike(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	userID := c.GetUint("user_id")
	if err := h.articleSvc.Unlike(userID, uint(postID)); err != nil {
		response.HandleError(c, err)
		return
	}
	responseData := common.ResponseMessage{
		Message: "取消点赞成功",
	}
	response.Success(c, responseData)
}
