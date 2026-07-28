package comment

import (
	"strconv"
	"tiny-forum/internal/model/common"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

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
