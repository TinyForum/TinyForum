package notification

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// base URL: /api/v1
func (h *NotificationHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	notifGroup := api.Group("/notifications", mw.AuthMW())
	{
		notifGroup.GET("", h.List) // 获取列表
		notifGroup.GET("/count/unread", h.UnreadCount) // 获取未读数量
		notifGroup.POST("/:id/read", h.MarkRead)
	}
	batchGroup := notifGroup.Group("/batch")
		{
			batchGroup.PATCH("/read", h.BatchMarkRead)
		}

}
