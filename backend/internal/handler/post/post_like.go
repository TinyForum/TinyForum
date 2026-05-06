package post

import (
	"strconv"

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
func (h *PostHandler) Like(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	userID := c.GetUint("user_id")
	if err := h.postSvc.Like(userID, uint(postID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "点赞成功"})
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
func (h *PostHandler) Unlike(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}
	userID := c.GetUint("user_id")
	if err := h.postSvc.Unlike(userID, uint(postID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "已取消点赞"})
}
