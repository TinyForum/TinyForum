package wire

import (
	"fmt"
	"time"

	"tiny-forum/config"
	announcementHandler "tiny-forum/internal/handler/announcement"
	answerHandler "tiny-forum/internal/handler/answer"
	authHandler "tiny-forum/internal/handler/auth"
	boardHandler "tiny-forum/internal/handler/board"
	commentHandler "tiny-forum/internal/handler/comment"
	notificationHandler "tiny-forum/internal/handler/notification"
	postHandler "tiny-forum/internal/handler/post"
	questionHandler "tiny-forum/internal/handler/questions"
	statsHandler "tiny-forum/internal/handler/stats"
	tagHandler "tiny-forum/internal/handler/tags"
	timelineHandler "tiny-forum/internal/handler/timelines"
	topicHandler "tiny-forum/internal/handler/topic"
	userHandler "tiny-forum/internal/handler/user"

	announcementService "tiny-forum/internal/service/announcement"
	boardService "tiny-forum/internal/service/board"
	commentService "tiny-forum/internal/service/comment"
	notificationService "tiny-forum/internal/service/notification"
	postService "tiny-forum/internal/service/post"
	questionService "tiny-forum/internal/service/question"
	statsService "tiny-forum/internal/service/stats"
	tagService "tiny-forum/internal/service/tag"
	timelineService "tiny-forum/internal/service/timeline"
	topicService "tiny-forum/internal/service/topic"
	userService "tiny-forum/internal/service/user"

	"tiny-forum/internal/middleware"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
	jwtpkg "tiny-forum/pkg/jwt"
	"tiny-forum/pkg/logger"

	_ "tiny-forum/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Cfg    *config.Config
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone,
	)

	logLevel := gormlogger.Silent
	if cfg.Server.Mode == "debug" {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate all models
	if err := db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Tag{},
		&model.Post{},
		&model.Comment{},
		&model.Like{},
		&model.Notification{},
		&model.SignIn{},
		&model.Report{},
		&model.Board{},
		&model.Moderator{},
		&model.BoardBan{},
		&model.ModeratorLog{},
		&model.Question{},
		&model.AnswerVote{},
		&model.TimelineEvent{},
		&model.UserTimeline{},
		&model.TimelineSubscription{},
		&model.Topic{},
		&model.TopicPost{},
		&model.TopicFollow{},
		&model.Announcement{},
		&model.ModeratorApplication{},
		&model.Moderator{},
		&model.Vote{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	// 核心配置：避免打爆 PostgreSQL
	sqlDB.SetMaxOpenConns(80)                 // 最大打开连接数（PG 默认 max_connections=100）
	sqlDB.SetMaxIdleConns(20)                 // 空闲连接池大小
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // 空闲连接超时

	logger.Info("Database connected and migrated successfully")
	return db, nil
}

func InitApp(cfg *config.Config) (*App, error) {
	// Init DB
	db, err := InitDB(cfg)
	if err != nil {
		return nil, err
	}

	// JWT manager
	jwtMgr := jwtpkg.NewManager(cfg.JWT.Secret, cfg.JWT.Expire)

	// ========== Repositories ==========
	tokenRepo := repository.NewTokenRepository(db)
	userRepo := repository.NewUserRepository(db, tokenRepo)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	tagRepo := repository.NewTagRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	timelineRepo := repository.NewTimelineRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	voteRepo := repository.NewVoteRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// ========== Services ==========
	// 基础服务
	notifSvc := notificationService.NewNotificationService(notifRepo)
	userSvc := userService.NewUserService(userRepo, jwtMgr, notifSvc)
	tagSvc := tagService.NewTagService(tagRepo)
	boardSvc := boardService.NewBoardService(boardRepo, userRepo, postRepo, notifSvc)
	timelineSvc := timelineService.NewTimelineService(timelineRepo, userRepo, postRepo, commentRepo)
	topicSvc := topicService.NewTopicService(topicRepo, postRepo, userRepo, notifSvc)
	questionSvc := questionService.NewQuestionService(questionRepo, postRepo, commentRepo, userRepo, notifSvc, tagRepo)
	postSvc := postService.NewPostService(postRepo, tagRepo, userRepo, boardRepo, notifSvc)
	commentSvc := commentService.NewCommentService(commentRepo, postRepo, userRepo, notifSvc, voteRepo)
	announcementSvc := announcementService.NewAnnouncementService(announcementRepo)
	statsSvc := statsService.NewStatsService(statsRepo, postRepo, tagRepo, boardRepo, userRepo, commentRepo)

	// ========== Handlers ==========
	authHandler := authHandler.NewAuthHandler(userSvc)
	userHandler := userHandler.NewUserHandler(userSvc, notifSvc)
	tagHandler := tagHandler.NewTagHandler(tagSvc)
	notifHandler := notificationHandler.NewNotificationHandler(notifSvc)
	postHandler := postHandler.NewPostHandler(postSvc)
	commentHandler := commentHandler.NewCommentHandler(commentSvc, questionSvc)
	boardHandler := boardHandler.NewBoardHandler(boardSvc)
	timelineHandler := timelineHandler.NewTimelineHandler(timelineSvc)
	topicHandler := topicHandler.NewTopicHandler(topicSvc)
	questionHandler := questionHandler.NewQuestionHandler(questionSvc, commentSvc, postSvc)
	answerHandler := answerHandler.NewAnswerHandler(questionSvc, commentSvc, postSvc)
	announcementHandler := announcementHandler.NewAnnouncementHandler(announcementSvc)
	statsHandler := statsHandler.NewStatsHandler(statsSvc)

	// ========== Gin Engine ==========
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},

		AllowCredentials: true,
	}))

	// Swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ========== Routes ==========
	api := engine.Group("/api/v1")

	// ----- MARK: Auth routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout) // 新增
		authGroup.GET("/me", middleware.Auth(jwtMgr), authHandler.Me)
	}

	// ----- MARK: Tag routes
	tagGroup := api.Group("/tags")
	{
		tagGroup.GET("", tagHandler.List)
		tagGroup.POST("", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Create)
		tagGroup.PUT("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Update)
		tagGroup.DELETE("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Delete)
	}

	// ----- MARK: Post routes
	postGroup := api.Group("/posts")
	{
		// 普通帖子
		postGroup.GET("", middleware.OptionalAuth(jwtMgr), postHandler.List)
		postGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), postHandler.GetByID)
		postGroup.POST("", middleware.Auth(jwtMgr), postHandler.Create)
		postGroup.PUT("/:id", middleware.Auth(jwtMgr), postHandler.Update)
		postGroup.DELETE("/:id", middleware.Auth(jwtMgr), postHandler.Delete)
		postGroup.POST("/:id/like", middleware.Auth(jwtMgr), postHandler.Like)
		postGroup.DELETE("/:id/like", middleware.Auth(jwtMgr), postHandler.Unlike)

	}

	// ----- MARK: Comment routes
	commentGroup := api.Group("/comments")
	{
		commentGroup.GET("/post/:post_id", commentHandler.List)
		commentGroup.POST("", middleware.Auth(jwtMgr), commentHandler.Create)
		commentGroup.DELETE("/:id", middleware.Auth(jwtMgr), commentHandler.Delete)

		// commentGroup.GET("/post/:post_id/answers", commentHandler.GetAnswers)
	}

	// ----- MARK: User routes
	userGroup := api.Group("/users")
	{
		userGroup.GET("/leaderboard", userHandler.Leaderboard)                                     // 获取排行榜
		userGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), userHandler.GetProfile)             // 获取用户信息
		userGroup.PUT("/profile", middleware.Auth(jwtMgr), userHandler.UpdateProfile)              // 更新用户信息
		userGroup.PATCH("/password", middleware.Auth(jwtMgr), userHandler.ChangePassword)          // 更新用户密码
		userGroup.POST("/:id/follow", middleware.Auth(jwtMgr), userHandler.Follow)                 // 关注用户
		userGroup.DELETE("/:id/follow", middleware.Auth(jwtMgr), userHandler.Unfollow)             // 取消关注用户
		userGroup.GET("/:id/followers", middleware.OptionalAuth(jwtMgr), userHandler.GetFollowers) // 获取用户的粉丝
		userGroup.GET("/:id/following", middleware.OptionalAuth(jwtMgr), userHandler.GetFollowing) // 获取用户的关注列表
		userGroup.GET("/:id/Score", middleware.OptionalAuth(jwtMgr), userHandler.GetScore)         // 获取用户的积分
		userGroup.GET("/me/role", middleware.OptionalAuth(jwtMgr), userHandler.GetCurrentUserRole) // 获取当前用户的角色
	}

	// ----- MARK: Notification routes
	notifGroup := api.Group("/notifications", middleware.Auth(jwtMgr))
	{
		notifGroup.GET("", notifHandler.List)
		notifGroup.GET("/unread-count", notifHandler.UnreadCount)
		notifGroup.POST("/read-all", notifHandler.MarkAllRead)
	}

	// ----- MARK: Board routes
	boardGroup := api.Group("/boards")
	{
		// ── 公开接口 ──────────────────────────────────────────────────────────
		boardGroup.GET("", boardHandler.List)                                                             // 获取板块列表
		boardGroup.GET("/tree", boardHandler.GetTree)                                                     // 获取板块树
		boardGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), boardHandler.GetByID)                     // 获取板块信息
		boardGroup.GET("/slug/:slug", boardHandler.GetBoardBySlug)                                        // 获取板块信息
		boardGroup.GET("/slug/:slug/posts", middleware.OptionalAuth(jwtMgr), boardHandler.GetPostsBySlug) // 获取板块帖子

		// ── 用户：申请 / 撤销版主申请 ────────────────────────────────────────
		// POST /boards/:id/moderators/apply   提交申请
		boardGroup.POST("/:id/moderators/apply-moderator",
			middleware.Auth(jwtMgr),
			boardHandler.ApplyModerator) // 提交申请
		// 查看申请状态
		boardGroup.GET("/moderators/apply",
			middleware.Auth(jwtMgr),
			boardHandler.GetUserApplications) // 查看申请状态
		// 用户获取自己管理的板块
		boardGroup.GET("/moderators/managed",
			middleware.Auth(jwtMgr),
			boardHandler.GetUserModeratorBoards) // 查看申请状态

		// 撤销申请（操作自己的申请）
		boardGroup.DELETE("/applications/:application_id",
			middleware.Auth(jwtMgr),
			boardHandler.CancelApplication) // 撤销申请

		// ── 管理员：板块 CRUD ──────────────────────────────────────────────
		boardGroup.POST("", middleware.Auth(jwtMgr), middleware.AdminRequired(), boardHandler.Create)
		boardGroup.PUT("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), boardHandler.Update)
		boardGroup.DELETE("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), boardHandler.Delete)

		// ── 版主管理（需要 manage_moderator 权限 或 管理员身份）───────────────
		modGroup := boardGroup.Group("/:id/moderators", middleware.Auth(jwtMgr))
		{
			// 查看版主列表（普通版主即可）
			modGroup.GET("",
				middleware.ModeratorRequired(jwtMgr, boardRepo),
				boardHandler.GetModerators)

			// 直接任命版主
			modGroup.POST("",
				middleware.CanManageModerator(jwtMgr, boardRepo),
				boardHandler.AddModerator)

			// 移除版主
			modGroup.DELETE("/:user_id",
				middleware.CanManageModerator(jwtMgr, boardRepo),
				boardHandler.RemoveModerator)

			// 升级 / 降级版主权限（仅管理员）
			modGroup.PUT("/:user_id/permissions",
				middleware.AdminRequired(),
				boardHandler.UpdateModeratorPermissions)
		}

		// ── 禁言管理（需要 ban_user 权限）────────────────────────────────────
		banGroup := boardGroup.Group("/:id/bans", middleware.Auth(jwtMgr))
		{
			banGroup.POST("", middleware.CanBanUser(jwtMgr, boardRepo), boardHandler.BanUser)
			banGroup.DELETE("/:user_id", middleware.CanBanUser(jwtMgr, boardRepo), boardHandler.UnbanUser)
		}

		// ── 帖子管理（版主）──────────────────────────────────────────────────
		postManageGroup := boardGroup.Group("/:id/posts", middleware.Auth(jwtMgr))
		{
			postManageGroup.DELETE("/:post_id",
				middleware.CanDeletePost(jwtMgr, boardRepo),
				boardHandler.DeletePost)
			postManageGroup.PUT("/:post_id/pin",
				middleware.CanPinPost(jwtMgr, boardRepo),
				boardHandler.PinPost)
		}
	}

	// ── 管理员：版主申请审批（挂在 /admin 下，已有 AdminRequired 中间件）──────
	// 建议放在 adminGroup 之内，与其他 admin 路由统一鉴权：
	adminBoardGroup := api.Group("/admin/boards", middleware.Auth(jwtMgr), middleware.AdminRequired())
	{
		// GET  /admin/boards/applications?board_id=&status=pending&page=1&page_size=20
		adminBoardGroup.GET("/applications", boardHandler.ListApplications)
		// POST /admin/boards/applications/:application_id/review
		adminBoardGroup.POST("/applications/:application_id/review", boardHandler.ReviewApplication)
	}

	// ----- MARK: Timeline routes
	timelineGroup := api.Group("/timeline", middleware.Auth(jwtMgr))
	{
		timelineGroup.GET("", timelineHandler.GetHomeTimeline)
		timelineGroup.GET("/following", timelineHandler.GetFollowingTimeline)
		timelineGroup.POST("/subscribe/:user_id", timelineHandler.Subscribe)
		timelineGroup.DELETE("/subscribe/:user_id", timelineHandler.Unsubscribe)
		timelineGroup.GET("/subscriptions", timelineHandler.GetSubscriptions)
	}

	// ----- MARK: Topic routes
	topicGroup := api.Group("/topics")
	{
		topicGroup.GET("", topicHandler.List)
		topicGroup.GET("/:id", topicHandler.GetByID)
		topicGroup.GET("/:id/posts", topicHandler.GetTopicPosts)
		topicGroup.GET("/:id/followers", topicHandler.GetFollowers)

		// 需要认证
		topicGroup.POST("", middleware.Auth(jwtMgr), topicHandler.Create)
		topicGroup.PUT("/:id", middleware.Auth(jwtMgr), topicHandler.Update)
		topicGroup.DELETE("/:id", middleware.Auth(jwtMgr), topicHandler.Delete)
		topicGroup.POST("/:id/posts", middleware.Auth(jwtMgr), topicHandler.AddPost)
		topicGroup.DELETE("/:id/posts/:post_id", middleware.Auth(jwtMgr), topicHandler.RemovePost)
		topicGroup.POST("/:id/follow", middleware.Auth(jwtMgr), topicHandler.Follow)
		topicGroup.DELETE("/:id/follow", middleware.Auth(jwtMgr), topicHandler.Unfollow)
	}

	// ----- MARK: Answer routes
	answerGroup := api.Group("/answers")
	{
		// 单个答案操作
		answerGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), answerHandler.GetAnswer) // 获取答案
		// answerGroup.PUT("/:id", middleware.Auth(jwtMgr), answerHandler.UpdateAnswer)      // 更新答案
		answerGroup.DELETE("/:id", middleware.Auth(jwtMgr), answerHandler.DeleteAnswer) // 删除答案

		// 答案交互
		answerGroup.GET("/:id/status", middleware.OptionalAuth(jwtMgr), answerHandler.GetVoteStatus) // 获取答案投票状态
		answerGroup.POST("/:id/vote", middleware.OptionalAuth(jwtMgr), answerHandler.VoteAnswer)     // 答案投票
		answerGroup.DELETE("/:id/vote", middleware.Auth(jwtMgr), answerHandler.RemoveVote)           // 取消投票
		answerGroup.POST("/:id/accept", middleware.Auth(jwtMgr), answerHandler.AcceptAnswer)         // 接受答案
		answerGroup.POST("/:id/unaccept", middleware.Auth(jwtMgr), answerHandler.UnacceptAnswer)     // 取消接受答案
	}

	// ----- MARK: Question routes
	questionGroup := api.Group("/questions")
	{
		questionGroup.GET("/simple", questionHandler.GetQuestionSimple)
		questionGroup.GET("/list", middleware.OptionalAuth(jwtMgr), questionHandler.GetQuestionsList)
		questionGroup.POST("/create", middleware.Auth(jwtMgr), questionHandler.CreateQuestion)
		questionGroup.GET("/detail/:id", middleware.OptionalAuth(jwtMgr), questionHandler.GetQuestionDetail)

		// 问题的答案
		questionGroup.GET("/:id/answers", middleware.OptionalAuth(jwtMgr), answerHandler.GetQuestionAnswers)
		questionGroup.POST("/:id/answers", middleware.Auth(jwtMgr), answerHandler.CreateAnswer)
	}

	// ========== MARK: Announcement routes ==========
	announcementGroup := api.Group("/announcements")
	{
		// 公开接口（所有用户可访问）
		announcementGroup.GET("", announcementHandler.List)
		announcementGroup.GET("/pinned", announcementHandler.GetPinned)
		announcementGroup.GET("/:id", announcementHandler.GetByID)
	}

	// 公告管理接口（需要管理员权限）
	announcementAdminGroup := api.Group("/admin/announcements", middleware.Auth(jwtMgr), middleware.AdminRequired())
	{
		announcementAdminGroup.GET("", announcementHandler.AdminList) // 管理员获取公告
		announcementAdminGroup.POST("", announcementHandler.Create)
		announcementAdminGroup.PUT("/:id", announcementHandler.Update)
		announcementAdminGroup.DELETE("/:id", announcementHandler.Delete)
		announcementAdminGroup.POST("/:id/publish", announcementHandler.Publish)
		announcementAdminGroup.POST("/:id/archive", announcementHandler.Archive)
		announcementAdminGroup.PUT("/:id/pin", announcementHandler.Pin)
	}
	statsGrop := api.Group("/statistics")
	{
		statsGrop.GET("", statsHandler.GetStatsTotal)
	}

	// ----- MARK: Admin routes -----
	adminGroup := api.Group("/admin", middleware.Auth(jwtMgr), middleware.AdminRequired())
	{
		adminGroup.GET("/users", userHandler.AdminList)
		adminGroup.PUT("/users/:id/active", userHandler.AdminSetActive)                  // 激活用户
		adminGroup.PUT("/users/:id/blocked", userHandler.AdminSetBlocked)                // 封禁用户
		adminGroup.DELETE("/users/:id/", userHandler.AdminDeleteUser)                    // 删除用户
		adminGroup.POST("/users/:id/reset-password", userHandler.AdminResetUserPassword) // 重置密码

		adminGroup.GET("/posts", postHandler.AdminList)
		adminGroup.PUT("/posts/:id/pin", postHandler.AdminTogglePin)
		adminGroup.PUT("/users/:id/role", userHandler.AdminSetRole)
		adminGroup.GET("/boards", boardHandler.List)
		// 平台统计
		adminGroup.GET("/statistics/day", statsHandler.GetStatsDay)     // 获取日数据
		adminGroup.GET("/statistics/total", statsHandler.GetStatsTotal) // 获取所有统计指标
		adminGroup.GET("/statistics/trend", statsHandler.GetStatsTrend) // 获取趋势指标
		// 积分
		adminGroup.GET("/users/score", userHandler.AdminGetUserScore) // 获取用户积分
		adminGroup.PUT("/users/:id/score", userHandler.AdminSetScore) // 设置用户积分

		// adminGroup.PUT("/boards/:id/sort", boardHandler.UpdateSortOrder)
	}

	return &App{Engine: engine, DB: db, Cfg: cfg}, nil
}
