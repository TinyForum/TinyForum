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
	mw *middleware.MiddlewareSet,
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
	handlers.Auth.RegisterRoutes(api, mw)
	// MARK: Auth routes
	// authGroup := api.Group("/auth")
	// {

	// 	authGroup.POST("/register", handlers.Auth.Register)                               // 用户注册
	// 	authGroup.POST("/login", handlers.Auth.Login)                                     // 用户登录
	// 	authGroup.POST("/logout", handlers.Auth.Logout)                                   // 用户登出
	// 	authGroup.DELETE("/delete-account", mw.AuthMW(), handlers.Auth.DeleteAccount)     // 用户删除账号
	// 	authGroup.GET("/deletion-status", mw.AuthMW(), handlers.Auth.DeletionStatus)      // 用户查询账号删除状态
	// 	authGroup.POST("/cancel-deletion", mw.AuthMW(), handlers.Auth.CancelDeletion)     // 用户取消账号删除
	// 	authGroup.DELETE("/confirm-deletion", mw.AuthMW(), handlers.Auth.ConfirmDeletion) // 用户确认账号删除
	// 	// authGroup.GET("/me", mw.AuthMW(), handlers.User.Me)                               // 获取当前用户信息
	// 	authGroup.POST("/forgot-password", handlers.Auth.ForgotPassword)                  // 用户忘记密码
	// 	authGroup.POST("/reset-password", handlers.Auth.ResetPassword)                    // 用户重置密码
	// 	authGroup.GET("/validate-reset-token", handlers.Auth.ValidateResetToken)          // 用户验证重置密码 token
	// }

	// MARK: Tag routes

	handlers.Tag.RegisterRoutes(api, mw)
	// tagGroup := api.Group("/tags")
	// {
	// 	tagGroup.GET("", handlers.Tag.List)
	// 	tagGroup.POST("", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Tag.Create)
	// 	tagGroup.PUT("/:id", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Tag.Update)
	// 	tagGroup.DELETE("/:id", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Tag.Delete)
	// }

	handlers.Post.RegisterRoutes(api, mw)
	// MARK: Post routes
	// postGroup := api.Group("/posts")
	// {
	// 	postGroup.GET("", mw.OptionalAuthMW(), handlers.Post.List) // 用户获取帖子列表
	// 	postGroup.GET("/:id", mw.OptionalAuthMW(), handlers.Post.GetByID)
	// 	postGroup.POST("",
	// 		mw.AuthMW(),
	// 		mw.RateLimitMW("create_post"),
	// 		mw.ContentCheckMW([]string{"title", "content"}),
	// 		handlers.Post.Create,
	// 	) // 用户发布帖子
	// 	postGroup.PUT("/:id", mw.AuthMW(), handlers.Post.Update)
	// 	postGroup.DELETE("/:id", mw.AuthMW(), handlers.Post.Delete)
	// 	postGroup.POST("/:id/like", mw.AuthMW(), handlers.Post.Like)
	// 	postGroup.DELETE("/:id/like", mw.AuthMW(), handlers.Post.Unlike)
	// }

	// MARK: Comment routes
	handlers.Comment.RegisterRoutes(api, mw)
	// commentGroup := api.Group("/comments")
	// {
	// 	commentGroup.GET("/post/:post_id", handlers.Comment.List)
	// 	commentGroup.POST("",
	// 		mw.AuthMW(),
	// 		mw.RateLimitMW("create_comment"),
	// 		mw.ContentCheckMW([]string{"content"}),
	// 		handlers.Comment.Create,
	// 	)
	// 	commentGroup.DELETE("/:id", mw.AuthMW(), handlers.Comment.Delete)
	// }

	// MARK: User routes
	handlers.User.RegisterRoutes(api, mw)
	// userGroup := api.Group("/users")
	// {

	// }

	// MARK: Notification routes
	handlers.Notification.RegisterRoutes(api, mw)

	// MARK: Board routes
	handlers.Board.RegisterRoutes(api, mw, repos.Board)
	// boardGroup := api.Group("/boards")
	// {

	// 	boardGroup.GET("", handlers.Board.List)
	// 	boardGroup.GET("/tree", handlers.Board.GetTree)

	// 	boardGroup.DELETE("/applications/:application_id", mw.AuthMW(), handlers.Board.CancelApplication)

	// 	boardGroup.POST("", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Board.Create)

	// 	modGroup := boardGroup.Group("/:id/moderators", mw.AuthMW())
	// 	{
	// 		modGroup.GET("", mw.ModeratorRequiredMW(repos.Board), handlers.Board.GetModerators)
	// 		modGroup.POST("", mw.CanManageModeratorMW(repos.Board), handlers.Board.AddModerator)
	// 		modGroup.DELETE("/:user_id", mw.CanManageModeratorMW(repos.Board), handlers.Board.RemoveModerator)
	// 		modGroup.PUT("/:user_id/permissions", mw.AdminRequiredMW(), handlers.Board.UpdateModeratorPermissions)
	// 	}

	// 	banGroup := boardGroup.Group("/:id/bans", mw.AuthMW())
	// 	{
	// 		banGroup.POST("", mw.CanBanUserMW(repos.Board), handlers.Board.BanUser)
	// 		banGroup.DELETE("/:user_id", mw.CanBanUserMW(repos.Board), handlers.Board.UnbanUser)
	// 	}

	// 	postManageGroup := boardGroup.Group("/:id/posts", mw.AuthMW())
	// 	{
	// 		postManageGroup.DELETE("/:post_id", mw.CanDeletePostMW(repos.Board), handlers.Board.DeletePost)
	// 		postManageGroup.PUT("/:post_id/pin", mw.CanPinPostMW(repos.Board), handlers.Board.PinPost)
	// 	}
	// }

	// MARK: Timeline routes
	handlers.Timeline.RegisterRoutes(api, mw)

	// MARK: Topic routes
	handlers.Topic.RegisterRoutes(api, mw)

	// MARK: Answer routes
	answerGroup := api.Group("/answers")
	{
		answerGroup.GET("/:id", mw.OptionalAuthMW(), handlers.Answer.GetAnswer)
		answerGroup.DELETE("/:id", mw.AuthMW(), handlers.Answer.DeleteAnswer)
		answerGroup.GET("/:id/status", mw.OptionalAuthMW(), handlers.Answer.GetVoteStatus)
		answerGroup.POST("/:id/vote", mw.OptionalAuthMW(), handlers.Answer.VoteAnswer)
		answerGroup.DELETE("/:id/vote", mw.AuthMW(), handlers.Answer.RemoveVote)
		answerGroup.POST("/:id/accept", mw.AuthMW(), handlers.Answer.AcceptAnswer)
		answerGroup.POST("/:id/unaccept", mw.AuthMW(), handlers.Answer.UnacceptAnswer)
	}

	// MARK: Question routes
	questionGroup := api.Group("/questions")
	{
		questionGroup.GET("/simple", handlers.Question.GetQuestionSimple)
		questionGroup.GET("/list", mw.OptionalAuthMW(), handlers.Question.GetQuestionsList)
		questionGroup.POST("/create", mw.AuthMW(), handlers.Question.CreateQuestion)
		questionGroup.GET("/detail/:id", mw.OptionalAuthMW(), handlers.Question.GetQuestionDetail)
		questionGroup.GET("/:id/answers", mw.OptionalAuthMW(), handlers.Answer.GetQuestionAnswers)
		questionGroup.POST("/:id/answers", mw.AuthMW(), handlers.Answer.CreateAnswer)
	}

	// MARK: Announcement routes
	announcementGroup := api.Group("/announcements")
	{
		announcementGroup.GET("", handlers.Announcement.List)
		announcementGroup.GET("/pinned", handlers.Announcement.GetPinned)
		announcementGroup.GET("/:id", handlers.Announcement.GetByID)
	}

	announcementAdminGroup := api.Group("/admin/announcements", mw.AuthMW(), mw.AdminRequiredMW())
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
	statsGroup := api.Group("/statistics")
	{
		statsGroup.GET("", handlers.Stats.GetStatsTotal)
	}

	// MARK: Admin routes
	adminGroup := api.Group("/admin", mw.AuthMW(), mw.AdminRequiredMW())
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
		handlers.Stats.RegisterRoutes(adminGroup)
		handlers.Risk.RegisterRoutes(adminGroup)
	}
	uploadGroup := api.Group("/upload", mw.AuthMW())
	{
		handlers.Upload.RegisterRoutes(uploadGroup)

	}
}
