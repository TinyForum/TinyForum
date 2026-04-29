package announcement

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *AnnouncementHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	announcementGroup := api.Group("/announcements")
	{
		announcementGroup.GET("", h.List)
		announcementGroup.GET("/pinned", h.GetPinned)
		announcementGroup.GET("/:id", h.GetByID)
	}
}
