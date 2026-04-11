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

// App holds all application dependencies
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

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	tagRepo := repository.NewTagRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// Services (build notification service first as dependency)
	notifSvc := service.NewNotificationService(notifRepo)
	userSvc := service.NewUserService(userRepo, jwtMgr, notifSvc)
	postSvc := service.NewPostService(postRepo, tagRepo, userRepo, notifSvc)
	commentSvc := service.NewCommentService(commentRepo, postRepo, userRepo, notifSvc)
	tagSvc := service.NewTagService(tagRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(userSvc)
	userHandler := handler.NewUserHandler(userSvc)
	postHandler := handler.NewPostHandler(postSvc)
	commentHandler := handler.NewCommentHandler(commentSvc)
	tagHandler := handler.NewTagHandler(tagSvc)
	notifHandler := handler.NewNotificationHandler(notifSvc)

	// Gin engine
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
	// swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	api := engine.Group("/api/v1")

	// Auth routes (public)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middleware.Auth(jwtMgr), authHandler.Me)
	}

	// Tag routes (list is public)
	tagGroup := api.Group("/tags")
	{
		tagGroup.GET("", tagHandler.List)
		tagGroup.POST("", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Create)
		tagGroup.PUT("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Update)
		tagGroup.DELETE("/:id", middleware.Auth(jwtMgr), middleware.AdminRequired(), tagHandler.Delete)
	}

	// Post routes
	postGroup := api.Group("/posts")
	{
		postGroup.GET("", middleware.OptionalAuth(jwtMgr), postHandler.List)
		postGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), postHandler.GetByID)
		postGroup.POST("", middleware.Auth(jwtMgr), postHandler.Create)
		postGroup.PUT("/:id", middleware.Auth(jwtMgr), postHandler.Update)
		postGroup.DELETE("/:id", middleware.Auth(jwtMgr), postHandler.Delete)
		postGroup.POST("/:id/like", middleware.Auth(jwtMgr), postHandler.Like)
		postGroup.DELETE("/:id/like", middleware.Auth(jwtMgr), postHandler.Unlike)
	}

	// Comment routes
	commentGroup := api.Group("/comments")
	{
		commentGroup.GET("/post/:post_id", commentHandler.List)
		commentGroup.POST("", middleware.Auth(jwtMgr), commentHandler.Create)
		commentGroup.DELETE("/:id", middleware.Auth(jwtMgr), commentHandler.Delete)
	}

	// User routes
	userGroup := api.Group("/users")
	{
		userGroup.GET("/leaderboard", userHandler.Leaderboard)
		userGroup.GET("/:id", middleware.OptionalAuth(jwtMgr), userHandler.GetProfile)
		userGroup.PUT("/profile", middleware.Auth(jwtMgr), userHandler.UpdateProfile)
		userGroup.POST("/:id/follow", middleware.Auth(jwtMgr), userHandler.Follow)
		userGroup.DELETE("/:id/follow", middleware.Auth(jwtMgr), userHandler.Unfollow)
	}

	// Notification routes
	notifGroup := api.Group("/notifications", middleware.Auth(jwtMgr))
	{
		notifGroup.GET("", notifHandler.List)
		notifGroup.GET("/unread-count", notifHandler.UnreadCount)
		notifGroup.POST("/read-all", notifHandler.MarkAllRead)
	}

	// Admin routes
	adminGroup := api.Group("/admin", middleware.Auth(jwtMgr), middleware.AdminRequired())
	{
		adminGroup.GET("/users", userHandler.AdminList)
		adminGroup.PUT("/users/:id/active", userHandler.AdminSetActive)
		adminGroup.GET("/posts", postHandler.AdminList)
		adminGroup.PUT("/posts/:id/pin", postHandler.AdminTogglePin)
	}

	return &App{Engine: engine, DB: db, Cfg: cfg}, nil
}
