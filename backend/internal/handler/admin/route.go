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
			announcementsGroup.POST("", h.CreateAnnouncement)
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
			usersGroup.PUT("/:id/active", h.SetActiveUser)
			usersGroup.PUT("/:id/blocked", h.SetBlockedUser)
			usersGroup.DELETE("/:id/", h.DeleteUser)
			usersGroup.PUT("/:id/role", h.SetRoleUser)

			// adminGroup.PUT("/users/:id/active", handlers.User.AdminSetActive)
			// adminGroup.PUT("/users/:id/blocked", handlers.User.AdminSetBlocked)
		}
		postsGroup := adminGroup.Group("/posts")
		{
			postsGroup.GET("", h.ListPosts)

		}

		// 	adminGroup.POST("/users/:id/reset-password", handlers.User.AdminResetUserPassword)
		// 	adminGroup.GET("/users/score", handlers.User.AdminGetUserScore)
		// 	adminGroup.PUT("/users/:id/score", handlers.User.AdminSetScore)
		// 	adminGroup.GET("/boards/applications", handlers.Board.ListApplications)
		// 	adminGroup.POST("/boards/applications/:application_id/review", handlers.Board.ReviewApplication)
		// 	adminGroup.GET("/boards", handlers.Board.List)
		// 	adminGroup.GET("/posts/pending", handlers.Post.AdminGetModerationRequire)
		// 	adminGroup.PUT("/audit/tasks/:id/approve", handlers.Post.AdminApprovePost)
		// 	adminGroup.PUT("/audit/tasks/:id/reject", handlers.Post.AdminRejectPost)
		// 	adminGroup.PUT("/posts/:id/pin", handlers.Post.AdminTogglePin)

	}
}
