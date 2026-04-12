package wire

import (
	"fmt"

	"tiny-forum/config"
	"tiny-forum/internal/handler"
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
	"tiny-forum/internal/service"
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
		// 新功能表
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
	); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

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
	// 原有
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	tagRepo := repository.NewTagRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// 新功能
	boardRepo := repository.NewBoardRepository(db)
	timelineRepo := repository.NewTimelineRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	questionRepo := repository.NewQuestionRepository(db)

	// ========== Services ==========
	// 基础服务
	notifSvc := service.NewNotificationService(notifRepo)
	userSvc := service.NewUserService(userRepo, jwtMgr, notifSvc)
	tagSvc := service.NewTagService(tagRepo)

	// 新功能服务（注意顺序：被依赖的先创建）
	boardSvc := service.NewBoardService(boardRepo, userRepo, postRepo, notifSvc)
	timelineSvc := service.NewTimelineService(timelineRepo, userRepo, postRepo, commentRepo)
	topicSvc := service.NewTopicService(topicRepo, postRepo, userRepo, notifSvc)
	questionSvc := service.NewQuestionService(questionRepo, postRepo, commentRepo, userRepo, notifSvc)

	// 依赖其他服务的服务
	postSvc := service.NewPostService(postRepo, tagRepo, userRepo, notifSvc)
	commentSvc := service.NewCommentService(commentRepo, postRepo, userRepo, notifSvc)

	// ========== Handlers ==========
	// 原有
	authHandler := handler.NewAuthHandler(userSvc)
	userHandler := handler.NewUserHandler(userSvc)
	tagHandler := handler.NewTagHandler(tagSvc)
	notifHandler := handler.NewNotificationHandler(notifSvc)

	// 更新构造函数（需要传入 questionSvc）
	postHandler := handler.NewPostHandler(postSvc, questionSvc)
	commentHandler := handler.NewCommentHandler(commentSvc, questionSvc)

	// 新功能
	boardHandler := handler.NewBoardHandler(boardSvc)
	timelineHandler := handler.NewTimelineHandler(timelineSvc)
	topicHandler := handler.NewTopicHandler(topicSvc)
	questionHandler := handler.NewQuestionHandler(questionSvc, commentSvc, postSvc)

	// ========== Gin Engine ==========
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ========== Routes ==========
	api := engine.Group("/api/v1")

	// ----- Auth routes -----
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middleware.Auth(jwtMgr), authHandler.Me)
	}

	// ----- Tag routes -----
	tagGroup := api.Group("/tags")
	{
		tagGroup.GET("", tagHandler.List)
		tagGroup.POST("", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Create)
		tagGroup.PUT("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Update)
		tagGroup.DELETE("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Delete)
	}

	// ----- Post routes (包含问答) -----
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

		// 问答相关
		postGroup.GET("/questions", middleware.OptionalAuth(jwtMgr), questionHandler.GetQuestions)
		postGroup.POST("/question", middleware.Auth(jwtMgr), questionHandler.CreateQuestion)
		postGroup.GET("/question/:id", middleware.OptionalAuth(jwtMgr), questionHandler.GetQuestionDetail)
		postGroup.POST("/questions/:post_id/answer/:comment_id/accept", middleware.Auth(jwtMgr), questionHandler.AcceptAnswer)
		postGroup.POST("/question/:id/answer", middleware.Auth(jwtMgr), questionHandler.CreateAnswer)
		postGroup.POST("/questions/answer/:comment_id/vote", middleware.OptionalAuth(jwtMgr), questionHandler.VoteAnswer)

		postGroup.GET("/questions/:post_id/answers", middleware.Auth(jwtMgr), questionHandler.GetQuestionAnswers)

	}

	// ----- Comment routes (包含答案投票) -----
	commentGroup := api.Group("/comments")
	{
		commentGroup.GET("/post/:post_id", commentHandler.List)
		commentGroup.POST("", middleware.Auth(jwtMgr), commentHandler.Create)
		commentGroup.DELETE("/:id", middleware.Auth(jwtMgr), commentHandler.Delete)

		// 答案投票
		commentGroup.POST("/:id/vote", middleware.Auth(jwtMgr), commentHandler.VoteAnswer)
		commentGroup.PUT("/:id/answer", middleware.Auth(jwtMgr), commentHandler.MarkAsAnswer)
		// commentGroup.GET("/post/:post_id/answers", commentHandler.GetAnswers)
	}

	// ----- User routes -----
	userGroup := api.Group("/users")
	{
		userGroup.GET("/leaderboard", userHandler.Leaderboard)
		userGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), userHandler.GetProfile)
		userGroup.PUT("/profile", middleware.Auth(jwtMgr), userHandler.UpdateProfile)
		userGroup.POST("/:id/follow", middleware.Auth(jwtMgr), userHandler.Follow)
		userGroup.DELETE("/:id/follow", middleware.Auth(jwtMgr), userHandler.Unfollow)
	}

	// ----- Notification routes -----
	notifGroup := api.Group("/notifications", middleware.Auth(jwtMgr))
	{
		notifGroup.GET("", notifHandler.List)
		notifGroup.GET("/unread-count", notifHandler.UnreadCount)
		notifGroup.POST("/read-all", notifHandler.MarkAllRead)
	}

	// ----- Board routes (板块) -----
	boardGroup := api.Group("/boards")
	{
		// 公开接口
		boardGroup.GET("", boardHandler.List)
		boardGroup.GET("/tree", boardHandler.GetTree)
		boardGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), boardHandler.GetByID)
		boardGroup.GET("/:id/posts", middleware.OptionalAuth(jwtMgr), boardHandler.GetPosts)
		boardGroup.GET("/slug/:slug/", boardHandler.GetBySlug)

		// 管理员接口
		boardGroup.POST("", middleware.Auth(jwtMgr), middleware.AdminRequired(), boardHandler.Create)
		boardGroup.PUT("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), boardHandler.Update)
		boardGroup.DELETE("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), boardHandler.Delete)

		// 版主管理（使用版主中间件）
		modGroup := boardGroup.Group("/:id/moderators", middleware.Auth(jwtMgr))
		{
			// 查看版主列表（普通版主可看）
			modGroup.GET("", middleware.ModeratorRequired(jwtMgr, boardRepo), boardHandler.GetModerators)

			// 添加/移除版主（需要管理版主权限）
			modGroup.POST("", middleware.CanManageModerator(jwtMgr, boardRepo), boardHandler.AddModerator)
			modGroup.DELETE("/:user_id", middleware.CanManageModerator(jwtMgr, boardRepo), boardHandler.RemoveModerator)
		}

		// 禁言管理（需要禁言权限）
		banGroup := boardGroup.Group("/:id/bans", middleware.Auth(jwtMgr))
		{
			banGroup.POST("", middleware.CanBanUser(jwtMgr, boardRepo), boardHandler.BanUser)
			banGroup.DELETE("/:user_id", middleware.CanBanUser(jwtMgr, boardRepo), boardHandler.UnbanUser)
		}

		// 帖子管理（需要对应权限）
		postManageGroup := boardGroup.Group("/:id/posts", middleware.Auth(jwtMgr))
		{
			// 删除帖子
			postManageGroup.DELETE("/:post_id", middleware.CanDeletePost(jwtMgr, boardRepo), boardHandler.DeletePost)
			// 置顶帖子
			postManageGroup.PUT("/:post_id/pin", middleware.CanPinPost(jwtMgr, boardRepo), boardHandler.PinPost)
		}
	}

	// ----- Timeline routes (时间线) -----
	timelineGroup := api.Group("/timeline", middleware.Auth(jwtMgr))
	{
		timelineGroup.GET("", timelineHandler.GetHomeTimeline)
		timelineGroup.GET("/following", timelineHandler.GetFollowingTimeline)
		timelineGroup.POST("/subscribe/:user_id", timelineHandler.Subscribe)
		timelineGroup.DELETE("/subscribe/:user_id", timelineHandler.Unsubscribe)
		timelineGroup.GET("/subscriptions", timelineHandler.GetSubscriptions)
	}

	// ----- Topic routes (专题) -----
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
	questionGroup := api.Group("/questions")
	{
		questionGroup.GET("/post/:post_id", middleware.OptionalAuth(jwtMgr), questionHandler.GetQuestionAnswers)
		questionGroup.POST("/answer/:comment_id/accept", middleware.Auth(jwtMgr), questionHandler.AcceptAnswer)
		questionGroup.POST("/answer/:comment_id/vote", middleware.Auth(jwtMgr), questionHandler.VoteAnswer)
	}
	// ----- Admin routes -----
	adminGroup := api.Group("/admin", middleware.Auth(jwtMgr), middleware.AdminRequired())
	{
		adminGroup.GET("/users", userHandler.AdminList)
		adminGroup.PUT("/users/:id/active", userHandler.AdminSetActive)
		adminGroup.PUT("/users/:id/blocked", userHandler.AdminSetBlocked)
		adminGroup.GET("/posts", postHandler.AdminList)
		adminGroup.PUT("/posts/:id/pin", postHandler.AdminTogglePin)
		// adminGroup.PUT("/posts/:id/pin-board", postHandler.AdminTogglePinInBoard)
		adminGroup.PUT("/users/:id/role", userHandler.AdminSetRole)

		// 板块管理
		adminGroup.GET("/boards", boardHandler.List)
		// adminGroup.PUT("/boards/:id/sort", boardHandler.UpdateSortOrder)
	}

	return &App{Engine: engine, DB: db, Cfg: cfg}, nil
}
