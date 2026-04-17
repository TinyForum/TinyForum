package notification

import (
	"strconv"

	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// List 获取通知列表
// @Summary 获取通知列表
// @Description 分页获取当前用户的通知列表
// @Tags 通知管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Notification}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	notifs, total, err := h.notifSvc.List(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, notifs, total, page, pageSize)
}

// UnreadCount 获取未读通知数量
// @Summary 获取未读通知数量
// @Description 获取当前用户未读通知的总数
// @Tags 通知管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications/unread-count [get]
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	count, err := h.notifSvc.UnreadCount(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": count})
}
