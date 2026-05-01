package admin

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *AdminHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	adminGroup := api.Group("/admin", mw.Auth(), mw.AdminRequired())
	{
		announcementsGroup := adminGroup.Group("/announcements")
		{
			announcementsGroup.GET("", h.ListAnnouncements)
			// announcementsGroup.POST("", h.Create)
			// announcementsGroup.PUT("/:id", h.Update)
			// announcementsGroup.DELETE("/:id", h.Delete)
			// announcementsGroup.POST("/:id/publish", h.Publish)
			// announcementsGroup.POST("/:id/archive", h.Archive)
			// announcementsGroup.PUT("/:id/pin", h.Pin)
		}

	}
}
