package comment

import (
	"strconv"
	"tiny-forum/internal/model/common"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Like 点赞评论
// @Summary 点赞评论
// @Tags 评论管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} common.BasicResponse  "点赞成功"
// @Failure 400 {object} common.BasicResponse "无效的评论ID"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /comments/{id}/like [post]
func (h *CommentHandler) Like(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	userID := c.GetUint("user_id")
	if err := h.commentSvc.Like(userID, uint(commentID)); err != nil {
		response.HandleError(c, err)
		return
	}
	responseData := common.ResponseMessage{
		Message: "点赞成功",
	}
	response.Success(c, responseData)
}

// Unlike 取消点赞评论
// @Summary 取消点赞评论
// @Tags 评论管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "评论ID"
// @Success 200 {object} common.BasicResponse  "已取消点赞"
// @Failure 400 {object} common.BasicResponse "无效的评论ID"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /comments/{id}/like [delete]
func (h *CommentHandler) Unlike(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	userID := c.GetUint("user_id")
	if err := h.commentSvc.Unlike(userID, uint(postID)); err != nil {
		response.HandleError(c, err)
		return
	}
	responseData := common.ResponseMessage{
		Message: "已取消点赞",
	}
	response.Success(c, responseData)
}
