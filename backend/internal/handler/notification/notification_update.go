package notification

import (
	"strconv"
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

// MarkRead 标记单个通知为已读
// @Summary 标记单个通知为已读
// @Description 将指定ID的通知标记为已读
// @Tags 通知管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path uint true "通知ID"
// @Success 200 {object} response.Response{data=object} "标记成功"
// @Failure 400 {object} response.Response "无效的通知ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权操作"
// @Failure 404 {object} response.Response "通知不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
    userID := c.GetUint("user_id")
    notifIDStr := c.Param("id")
    
    // 参数验证在 Handler 层
    notifID, err := strconv.ParseUint(notifIDStr, 10, 64)
    if err != nil {
        response.BadRequest(c, "无效的通知ID")
        return
    }
    
    // 调用 Service 层（传递正确的类型）
    if err := h.notifSvc.MarkRead(uint(notifID), userID); err != nil {
        switch err.Error() {
        case "通知不存在":
            response.NotFound(c, err.Error())
        case "无权操作此通知":
            response.Forbidden(c, err.Error())
        default:
            response.InternalError(c, err.Error())
        }
        return
    }
    response.Success(c, gin.H{"message": "已标记为已读"})
}