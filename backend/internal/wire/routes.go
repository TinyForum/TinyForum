package wire

import (
	"tiny-forum/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(
	engine *gin.Engine,
	handlers *Handlers,
	mw *MiddlewareSet,
	repos *Repositories,
	cfg *config.Config,

) {
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

	// MARK: Auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", handlers.Auth.Register)
		authGroup.POST("/login", handlers.Auth.Login)
		authGroup.POST("/logout", handlers.Auth.Logout)
		authGroup.DELETE("/delete-account", mw.AuthMW(), handlers.Auth.DeleteAccount)
		authGroup.GET("/deletion-status", mw.AuthMW(), handlers.Auth.DeletionStatus)
		authGroup.POST("/cancel-deletion", mw.AuthMW(), handlers.Auth.CancelDeletion)
		authGroup.DELETE("/confirm-deletion", mw.AuthMW(), handlers.Auth.ConfirmDeletion)
		authGroup.GET("/me", mw.AuthMW(), handlers.User.Me)
		authGroup.POST("/forgot-password", handlers.Auth.ForgotPassword)
		authGroup.POST("/reset-password", handlers.Auth.ResetPassword)
		authGroup.GET("/validate-reset-token", handlers.Auth.ValidateResetToken)
	}

	// MARK: Tag routes
	tagGroup := api.Group("/tags")
	{
		tagGroup.GET("", handlers.Tag.List)
		tagGroup.POST("", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Tag.Create)
		tagGroup.PUT("/:id", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Tag.Update)
		tagGroup.DELETE("/:id", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Tag.Delete)
	}

	// MARK: Post routes
	postGroup := api.Group("/posts")
	{
		postGroup.GET("", mw.OptionalAuthMW(), handlers.Post.List) // 用户获取帖子列表
		postGroup.GET("/:id", mw.OptionalAuthMW(), handlers.Post.GetByID)
		postGroup.POST("",
			mw.AuthMW(),
			mw.RateLimitMW("create_post"),
			mw.ContentCheckMW([]string{"title", "content"}),
			handlers.Post.Create,
		)
		postGroup.PUT("/:id", mw.AuthMW(), handlers.Post.Update)
		postGroup.DELETE("/:id", mw.AuthMW(), handlers.Post.Delete)
		postGroup.POST("/:id/like", mw.AuthMW(), handlers.Post.Like)
		postGroup.DELETE("/:id/like", mw.AuthMW(), handlers.Post.Unlike)
	}

	// MARK: Comment routes
	commentGroup := api.Group("/comments")
	{
		commentGroup.GET("/post/:post_id", handlers.Comment.List)
		commentGroup.POST("",
			mw.AuthMW(),
			mw.RateLimitMW("create_comment"),
			mw.ContentCheckMW([]string{"content"}),
			handlers.Comment.Create,
		)
		commentGroup.DELETE("/:id", mw.AuthMW(), handlers.Comment.Delete)
	}

	// MARK: User routes
	userGroup := api.Group("/users")
	{
		userGroup.GET("/leaderboard/simple", handlers.User.LeaderboardSimple)
		userGroup.GET("/leaderboard/detail", handlers.User.LeaderboardDetail)
		userGroup.GET("/:id", mw.OptionalAuthMW(), handlers.User.GetProfile)
		userGroup.PUT("/profile", mw.AuthMW(), handlers.User.UpdateProfile)
		userGroup.PATCH("/password", mw.AuthMW(), handlers.User.ChangePassword)
		userGroup.POST("/:id/follow", mw.AuthMW(), handlers.User.Follow)
		userGroup.DELETE("/:id/follow", mw.AuthMW(), handlers.User.Unfollow)
		userGroup.GET("/:id/followers", mw.OptionalAuthMW(), handlers.User.GetFollowers)
		userGroup.GET("/:id/following", mw.OptionalAuthMW(), handlers.User.GetFollowing)
		userGroup.GET("/:id/Score", mw.OptionalAuthMW(), handlers.User.GetScore)
		userGroup.GET("/me/role", mw.OptionalAuthMW(), handlers.User.GetCurrentUserRole)
	}

	// MARK: Notification routes
	notifGroup := api.Group("/notifications", mw.AuthMW())
	{
		notifGroup.GET("", handlers.Notification.List)
		notifGroup.GET("/unread-count", handlers.Notification.UnreadCount)
		notifGroup.POST("/read-all", handlers.Notification.MarkAllRead)
	}

	// MARK: Board routes
	boardGroup := api.Group("/boards")
	{
		boardGroup.GET("", handlers.Board.List)
		boardGroup.GET("/tree", handlers.Board.GetTree)
		boardGroup.GET("/:id", mw.OptionalAuthMW(), handlers.Board.GetByID)
		boardGroup.GET("/slug/:slug", handlers.Board.GetBoardBySlug)
		boardGroup.GET("/slug/:slug/posts", mw.OptionalAuthMW(), handlers.Board.GetPostsBySlug)

		boardGroup.POST("/:id/moderators/apply-moderator", mw.AuthMW(), handlers.Board.ApplyModerator)
		boardGroup.GET("/moderators/apply-status", mw.AuthMW(), handlers.Board.GetUserApplications)
		boardGroup.GET("/moderators/managed", mw.AuthMW(), handlers.Board.GetUserModeratorBoards)
		boardGroup.DELETE("/applications/:application_id", mw.AuthMW(), handlers.Board.CancelApplication)

		boardGroup.POST("", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Board.Create)
		boardGroup.PUT("/:id", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Board.Update)
		boardGroup.DELETE("/:id", mw.AuthMW(), mw.AdminRequiredMW(), handlers.Board.Delete)

		modGroup := boardGroup.Group("/:id/moderators", mw.AuthMW())
		{
			modGroup.GET("", mw.ModeratorRequiredMW(repos.Board), handlers.Board.GetModerators)
			modGroup.POST("", mw.CanManageModeratorMW(repos.Board), handlers.Board.AddModerator)
			modGroup.DELETE("/:user_id", mw.CanManageModeratorMW(repos.Board), handlers.Board.RemoveModerator)
			modGroup.PUT("/:user_id/permissions", mw.AdminRequiredMW(), handlers.Board.UpdateModeratorPermissions)
		}

		banGroup := boardGroup.Group("/:id/bans", mw.AuthMW())
		{
			banGroup.POST("", mw.CanBanUserMW(repos.Board), handlers.Board.BanUser)
			banGroup.DELETE("/:user_id", mw.CanBanUserMW(repos.Board), handlers.Board.UnbanUser)
		}

		postManageGroup := boardGroup.Group("/:id/posts", mw.AuthMW())
		{
			postManageGroup.DELETE("/:post_id", mw.CanDeletePostMW(repos.Board), handlers.Board.DeletePost)
			postManageGroup.PUT("/:post_id/pin", mw.CanPinPostMW(repos.Board), handlers.Board.PinPost)
		}
	}

	// MARK: Timeline routes
	timelineGroup := api.Group("/timeline", mw.AuthMW())
	{
		timelineGroup.GET("", handlers.Timeline.GetHomeTimeline)
		timelineGroup.GET("/following", handlers.Timeline.GetFollowingTimeline)
		timelineGroup.POST("/subscribe/:user_id", handlers.Timeline.Subscribe)
		timelineGroup.DELETE("/subscribe/:user_id", handlers.Timeline.Unsubscribe)
		timelineGroup.GET("/subscriptions", handlers.Timeline.GetSubscriptions)
	}

	// MARK: Topic routes
	topicGroup := api.Group("/topics")
	{
		topicGroup.GET("", handlers.Topic.List)
		topicGroup.GET("/:id", handlers.Topic.GetByID)
		topicGroup.GET("/:id/posts", handlers.Topic.GetTopicPosts)
		topicGroup.GET("/:id/followers", handlers.Topic.GetFollowers)
		topicGroup.POST("", mw.AuthMW(), handlers.Topic.Create)
		topicGroup.PUT("/:id", mw.AuthMW(), handlers.Topic.Update)
		topicGroup.DELETE("/:id", mw.AuthMW(), handlers.Topic.Delete)
		topicGroup.POST("/:id/posts", mw.AuthMW(), handlers.Topic.AddPost)
		topicGroup.DELETE("/:id/posts/:post_id", mw.AuthMW(), handlers.Topic.RemovePost)
		topicGroup.POST("/:id/follow", mw.AuthMW(), handlers.Topic.Follow)
		topicGroup.DELETE("/:id/follow", mw.AuthMW(), handlers.Topic.Unfollow)
	}

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
}
