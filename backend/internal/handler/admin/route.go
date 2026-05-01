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
			announcementsGroup.GET("/:id", h.CreateAnnouncement)
			announcementsGroup.PUT("/:id", h.UpdateAnnouncement)
			announcementsGroup.DELETE("/:id", h.DeleteAnnouncement)
			announcementsGroup.POST("/:id/publish", h.PublishAnnouncement)
			announcementsGroup.POST("/:id/archive", h.ArchiveAnnouncement)
			announcementsGroup.PUT("/:id/pin", h.PinAnnouncement)
			// announcementsGroup.POST("", h.Create)
			// announcementsGroup.PUT("/:id", h.Update)
			// announcementsGroup.DELETE("/:id", h.Delete)
			// announcementsGroup.POST("/:id/publish", h.Publish)
			// announcementsGroup.POST("/:id/archive", h.Archive)
			// announcementsGroup.PUT("/:id/pin", h.Pin)
		}
		usersGroup := adminGroup.Group("/users")
		{
			usersGroup.GET("", h.ListUsers)
					usersGroup.PUT("/users/:id/active", h.SetActiveUser)
						usersGroup.PUT("/users/:id/blocked", h.SetBlockedUser)
						usersGroup.DELETE("/users/:id/", h.DeleteUser)
						usersGroup.PUT("/users/:id/role", h.SetRoleUser)
		
			// adminGroup.PUT("/users/:id/active", handlers.User.AdminSetActive)
			// adminGroup.PUT("/users/:id/blocked", handlers.User.AdminSetBlocked)
		}

	}
}
