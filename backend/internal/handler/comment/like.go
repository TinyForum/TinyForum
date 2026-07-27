package comment

import (
	"strconv"
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
	response.Success(c, gin.H{"message": "点赞成功"})
}
