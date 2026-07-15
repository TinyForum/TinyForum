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
			announcementsGroup.GET("", h.ListAnnouncements)                // 获取公告列表
			announcementsGroup.POST("", h.CreateAnnouncement)              // 创建公告
			announcementsGroup.PUT("/:id", h.UpdateAnnouncement)           // 更新公告
			announcementsGroup.DELETE("/:id", h.DeleteAnnouncement)        // 删除公告
			announcementsGroup.POST("/:id/publish", h.PublishAnnouncement) // 发布公告
			announcementsGroup.POST("/:id/archive", h.ArchiveAnnouncement) // 归档公告
			announcementsGroup.PUT("/:id/pin", h.PinAnnouncement)          // 置顶公告
			// announcementsGroup.POST("", h.Create)
			// announcementsGroup.PUT("/:id", h.Update)
			// announcementsGroup.DELETE("/:id", h.Delete)
			// announcementsGroup.POST("/:id/publish", h.Publish)
			// announcementsGroup.POST("/:id/archive", h.Archive)
			// announcementsGroup.PUT("/:id/pin", h.Pin)
		}
		usersGroup := adminGroup.Group("/users")
		{
			usersGroup.GET("", h.ListUsers)                  // 获取用户列表
			usersGroup.PUT("/:id/active", h.SetActiveUser)   // 激活用户
			usersGroup.PUT("/:id/blocked", h.SetBlockedUser) // 设置用户是否被禁用
			usersGroup.DELETE("/:id", h.DeleteUser)          // 删除用户
			usersGroup.PUT("/:id/role", h.SetRoleUser)       // 设置用户角色
			usersGroup.GET("/score", h.ListUsersScore)       // 获取用户积分列表
			usersGroup.GET("/:id/score", h.GetUserScore)     // 获取用户积分

			// adminGroup.PUT("/users/:id/active", handlers.User.AdminSetActive)
			// adminGroup.PUT("/users/:id/blocked", handlers.User.AdminSetBlocked)
		}
		postsGroup := adminGroup.Group("/posts")
		{
			postsGroup.GET("", h.ListPosts)                 // 获取帖子列表
			postsGroup.GET("/pending", h.ListReviewRequire) // 获取待审核帖子列表
		}
		boardGroup := adminGroup.Group("/boards")
		{
			boardGroup.GET("/applications", h.ListApplications)                          // 获取版块申请列表
			boardGroup.POST("/applications/:application_id/review", h.ReviewApplication) // 审核版块申请
		}
		reportsGroup := adminGroup.Group("/reports")
		{
			reportsGroup.GET("", h.ListReports) // 获取举报列表
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
