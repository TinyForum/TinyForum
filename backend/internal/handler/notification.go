package handler

import (
	"strconv"

	"bbs-forum/internal/service"
	"bbs-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifSvc *service.NotificationService
}

func NewNotificationHandler(notifSvc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}

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

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.notifSvc.MarkAllRead(userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "已全部标记为已读"})
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	count, err := h.notifSvc.UnreadCount(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"count": count})
}
