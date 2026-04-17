package notification

import (
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// MarkAllRead 标记所有通知为已读
// @Summary 标记所有通知为已读
// @Description 将当前用户的所有未读通知标记为已读
// @Tags 通知管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=object} "标记成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications/read-all [post]
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.notifSvc.MarkAllRead(userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "已全部标记为已读"})
}
