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
		notifGroup.GET("", h.List)
		notifGroup.GET("/unread-count", h.UnreadCount)
		notifGroup.POST("/read-all", h.MarkAllRead)
	}

}
