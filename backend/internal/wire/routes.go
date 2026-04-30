package wire

import (
	"fmt"
	"tiny-forum/config"
	"tiny-forum/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(
	engine *gin.Engine,
	handlers *Handlers,
	mw middleware.MiddlewareSet,
	repos *Repositories,
	cfg *config.Config,

) {
	fmt.Printf("DEBUG: AllowOrigins = %v, len = %d\n", cfg.Basic.AllowOrigins, len(cfg.Basic.AllowOrigins))
	// CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Basic.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) { c.JSON(200, "pong") })
	}
	handlers.Auth.RegisterRoutes(api, mw)               // 验证路由（权限相关操作，密码）
	handlers.Tag.RegisterRoutes(api, mw)                // 标签路由
	handlers.Post.RegisterRoutes(api, mw)               // 帖子路由
	handlers.Comment.RegisterRoutes(api, mw)            // 评论路由
	handlers.User.RegisterRoutes(api, mw)               // 用户路由（用户信息、排名）【密码在验证路由】
	handlers.Notification.RegisterRoutes(api, mw)       // 通知路由
	handlers.Board.RegisterRoutes(api, mw, repos.Board) // 板块路由
	handlers.Timeline.RegisterRoutes(api, mw)           // 时间线路由
	handlers.Topic.RegisterRoutes(api, mw)              // 主题路由
	handlers.Answer.RegisterRoutes(api, mw)             // 答案路由
	handlers.Question.RegisterRoutes(api, mw)           // 问题路由
	handlers.Announcement.RegisterRoutes(api, mw)       // 公告路由
	announcementAdminGroup := api.Group("/admin/announcements", mw.Auth(), mw.AdminRequired())
	{
		announcementAdminGroup.GET("", handlers.Announcement.AdminList)
		announcementAdminGroup.POST("", handlers.Announcement.Create)
		announcementAdminGroup.PUT("/:id", handlers.Announcement.Update)
		announcementAdminGroup.DELETE("/:id", handlers.Announcement.Delete)
		announcementAdminGroup.POST("/:id/publish", handlers.Announcement.Publish)
		announcementAdminGroup.POST("/:id/archive", handlers.Announcement.Archive)
		announcementAdminGroup.PUT("/:id/pin", handlers.Announcement.Pin)
	}

	// MARK: Statistics routes
	handlers.Stats.RegisterRoutes(api, mw)

	// FIXME: 通过上下文判断，而不是路径
	// MARK: Admin routes

	adminGroup := api.Group("/admin", mw.Auth(), mw.AdminRequired())
	{
		adminGroup.GET("/users", handlers.User.AdminList)
		adminGroup.PUT("/users/:id/active", handlers.User.AdminSetActive)
		adminGroup.PUT("/users/:id/blocked", handlers.User.AdminSetBlocked)
		adminGroup.DELETE("/users/:id/", handlers.User.AdminDeleteUser)
		adminGroup.PUT("/users/:id/role", handlers.User.AdminSetRole)
		adminGroup.POST("/users/:id/reset-password", handlers.User.AdminResetUserPassword)
		adminGroup.GET("/users/score", handlers.User.AdminGetUserScore)
		adminGroup.PUT("/users/:id/score", handlers.User.AdminSetScore)
		adminGroup.GET("/boards/applications", handlers.Board.ListApplications)
		adminGroup.POST("/boards/applications/:application_id/review", handlers.Board.ReviewApplication)
		adminGroup.GET("/boards", handlers.Board.List)
		adminGroup.GET("/posts", handlers.Post.AdminList)
		adminGroup.GET("/posts/pending", handlers.Post.AdminGetModerationRequire)
		adminGroup.PUT("/audit/tasks/:id/approve", handlers.Post.AdminApprovePost)
		adminGroup.PUT("/audit/tasks/:id/reject", handlers.Post.AdminRejectPost)
		adminGroup.PUT("/posts/:id/pin", handlers.Post.AdminTogglePin)

		// MARK: 挂载子路由
		handlers.Risk.RegisterRoutes(adminGroup)
	}

	handlers.Upload.RegisterRoutes(api, mw)

}
