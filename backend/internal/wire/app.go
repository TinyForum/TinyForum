// internal/wire/app.go
//
// 变更说明：
//   1. InitInfra 调用新增 db 参数
//   2. NewMiddlewareSet 调用新增 infra.Enforcer 参数

package wire

import (
	"fmt"
	initdata "tiny-forum/init"
	"tiny-forum/internal/botapi"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/job"
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/service/bot"
	"tiny-forum/internal/storage"
	"tiny-forum/internal/strategy"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Cfg    *config.Config
	BotSvc bot.Service
}

func InitApp(cfg *config.Config) (*App, error) {
	// 1. 数据库
	db, err := InitDB(cfg)
	if err != nil {
		return nil, err
	}

	// 2. JWT 管理器
	jwtMgr := jwtpkg.NewJWTManager(cfg.Private.JWT.Secret, cfg.Private.JWT.Expire)

	// 3. 基础设施（Redis、限流、敏感词、Casbin）
	// db 传入是为了让 Casbin 的 GORM adapter 复用同一个数据库连接
	infra, err := InitInfra(cfg, db)
	if err != nil {
		return nil, err
	}
	registry := strategy.NewHandlerRegistry()

	userStorage := storage.NewLocalStorage("./uploads")
	pluginsStorage := storage.NewLocalStorage("./store")

	// 4. 数据仓库层
	repos := NewRepositories(db, infra.RedisClient)
	// api
	bot, err := initdata.InitDefaultBot(db) // 调用上面的函数
	if err != nil {
		return nil, fmt.Errorf("init default bot: %w", err)
	}

	// 6. 创建 ForumAPI（使用上面得到的 bot 实例）
	forumAPI := botapi.NewForumAPI(
		bot.ID, // *do.Bot
		// int64(bot.BaseModel.ID), // botActorID
		repos.Post,
		repos.Comment,
		repos.User,
		// repos.
		// repos.Message, // 确保 repos.Message 存在且类型正确
		repos.Notification,
		// repos.Stats,
	)

	// 5. 服务层
	services := NewServices(cfg, jwtMgr, repos, infra, userStorage, pluginsStorage, registry, forumAPI)
	services.Bot.StartScheduler()

	// 6. 辅助工具
	helpers := NewHelpers()

	// 7. 控制器层
	handlers := NewHandlers(services, helpers.TimeHelpers, cfg)

	// 8. 中间件层（新增 infra.Enforcer 参数）
	mw := middleware.NewMiddlewareSet(
		jwtMgr,
		db,
		services.Risk,
		services.ContentCheck,
		repos.Token,
		&cfg.RiskControl.RateLimit,
		infra.Enforcer, // Casbin enforcer
	)

	// 9. Gin 引擎
	gin.SetMode(cfg.Basic.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 10. 注册路由
	RegisterRoutes(engine, handlers, mw, repos, cfg)

	// 自动清理
	go job.CleanTempFiles(db, userStorage, repos.Attachment)
	go job.CleanTempFiles(db, pluginsStorage, repos.Attachment)

	// services.Bot.StopScheduler()
	return &App{
		Engine: engine,
		DB:     db,
		Cfg:    cfg,
		BotSvc: services.Bot,
	}, nil
}
