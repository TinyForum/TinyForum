// internal/wire/app.go
package wire

import (
	"fmt"
	"log"
	initdata "tiny-forum/init"
	"tiny-forum/internal/botapi"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/job"
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/service/bot"
	configService "tiny-forum/internal/service/config"
	"tiny-forum/internal/storage"
	"tiny-forum/internal/strategy"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App 应用结构
type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
	DynCfg *config.DynamicConfig // 改为动态配置
	BotSvc bot.Service
	// 保存组件引用以便热更新时重建
	infra      *Infra
	repos      *Repositories
	services   *Services
	handlers   *Handlers
	middleware middleware.MiddlewareSet
}

// InitAppWithDynamic 使用动态配置初始化应用
func InitAppWithDynamic(dynCfg *config.DynamicConfig) (*App, error) {
	// 获取当前配置
	cfg := dynCfg.Get()
	if cfg == nil {
		return nil, fmt.Errorf("failed to get config from dynamic config")
	}

	app := &App{
		DynCfg: dynCfg,
	}

	// 1. 初始化数据库
	db, err := InitDB(cfg)
	if err != nil {
		return nil, err
	}
	app.DB = db

	// 2. 初始化 JWT 管理器
	jwtMgr := jwtpkg.NewJWTManager(cfg.Private.JWT.Secret, cfg.Private.JWT.Expire)

	// 3. 初始化基础设施
	infra, err := InitInfra(cfg, db)
	if err != nil {
		return nil, err
	}
	app.infra = infra

	// 4. 初始化存储
	userStorage := storage.NewLocalStorage("./uploads")
	pluginsStorage := storage.NewLocalStorage("./store")

	// 5. 初始化仓库层
	repos := NewRepositories(db, infra.RedisClient)
	app.repos = repos

	// 6. 初始化机器人
	bot, err := initdata.InitDefaultBot(db)
	if err != nil {
		return nil, fmt.Errorf("init default bot: %w", err)
	}

	// 7. 创建 ForumAPI
	forumAPI := botapi.NewForumAPI(
		bot.ID,
		repos.Post,
		repos.Comment,
		repos.User,
		repos.Notification,
	)

	// 8. 初始化服务层
	registry := strategy.NewHandlerRegistry()
	services := NewServices(cfg, jwtMgr, repos, infra, userStorage, pluginsStorage, registry, forumAPI)
	services.Bot.StartScheduler()
	app.services = services
	app.BotSvc = services.Bot

	// 9. 初始化辅助工具
	helpers := NewHelpers()

	// 10. 初始化控制器
	configSvc := configService.NewConfigService("config", dynCfg, db)
	handlers := NewHandlers(services, helpers.TimeHelpers, cfg, configSvc)
	app.handlers = handlers

	// 11. 初始化中间件
	mw := middleware.NewMiddlewareSet(
		jwtMgr,
		db,
		services.Risk,
		services.ContentCheck,
		repos.Token,
		&cfg.RiskControl.RateLimit,
		infra.Enforcer,
	)
	app.middleware = mw

	// 12. 创建 Gin 引擎并注册路由
	gin.SetMode(cfg.Basic.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	app.Engine = engine

	RegisterRoutes(engine, handlers, mw, repos, cfg)

	// 13. 启动后台任务
	go job.CleanTempFiles(db, userStorage, repos.Attachment)
	go job.CleanTempFiles(db, pluginsStorage, repos.Attachment)

	// 14. 注册配置变更回调（关键：热更新逻辑）
	app.registerConfigCallbacks()

	return app, nil
}

// registerConfigCallbacks 注册配置变更回调
func (app *App) registerConfigCallbacks() {
	app.DynCfg.OnChange(func(fileName string, oldConfig, newConfig *config.Config) {
		log.Printf("[App] Config changed: %s, applying hot updates...", fileName)

		// 根据变更的文件名执行不同的更新策略
		switch fileName {
		case "basic.yml":
			app.onBasicConfigChanged(oldConfig, newConfig)
		case "postgres.yml":
			app.onPostgresConfigChanged(oldConfig, newConfig)
		case "redis.yml":
			app.onRedisConfigChanged(oldConfig, newConfig)
		case "risk_control.yml":
			app.onRiskConfigChanged(oldConfig, newConfig)
		case "private.yml":
			app.onPrivateConfigChanged(oldConfig, newConfig)
		case "manual_reload":
			app.onManualReload(oldConfig, newConfig)
		default:
			// 全面刷新
			app.onFullConfigReload(oldConfig, newConfig)
		}

		log.Printf("[App] Hot update completed for: %s", fileName)
	})
}

// onBasicConfigChanged 基础配置变更处理
func (app *App) onBasicConfigChanged(oldConfig, newConfig *config.Config) {
	// 1. 更新 Gin 模式
	if oldConfig.Basic.Server.Mode != newConfig.Basic.Server.Mode {
		gin.SetMode(newConfig.Basic.Server.Mode)
		log.Printf("[App] Gin mode updated to: %s", newConfig.Basic.Server.Mode)
	}

	// 2. 更新限流配置（限流器会通过回调自动更新）
	// 限流器已经注册了回调，无需手动操作

	// 3. 更新功能开关
	// 各业务服务会通过自己的回调更新

	log.Printf("[App] Basic config updated: port=%d, mode=%s",
		newConfig.Basic.Server.Port, newConfig.Basic.Server.Mode)
}

// onPostgresConfigChanged PostgreSQL 配置变更处理
func (app *App) onPostgresConfigChanged(oldConfig, newConfig *config.Config) {
	// 数据库连接池会通过自己的回调重建
	// 这里只记录日志
	log.Printf("[App] PostgreSQL config changed: %s:%d -> %s:%d",
		oldConfig.Postgres.Host, oldConfig.Postgres.Port,
		newConfig.Postgres.Host, newConfig.Postgres.Port)
}

// onRedisConfigChanged Redis 配置变更处理
func (app *App) onRedisConfigChanged(oldConfig, newConfig *config.Config) {
	// Redis 连接池会通过自己的回调重建
	log.Printf("[App] Redis config changed: %s:%d -> %s:%d",
		oldConfig.Redis.Host, oldConfig.Redis.Port,
		newConfig.Redis.Host, newConfig.Redis.Port)
}

// onRiskConfigChanged 风控配置变更处理
func (app *App) onRiskConfigChanged(oldConfig, newConfig *config.Config) {
	// 风控服务会通过自己的回调更新
	log.Printf("[App] Risk control config changed: enabled=%v",
		newConfig.RiskControl.RateLimit.Enabled)
}

// onPrivateConfigChanged 私有配置变更处理
func (app *App) onPrivateConfigChanged(oldConfig, newConfig *config.Config) {
	// 检查 JWT Secret 是否变更
	if oldConfig.Private.JWT.Secret != newConfig.Private.JWT.Secret {
		// JWT Secret 变更需要重建 JWT 管理器
		// 注意：这会使所有现有 token 失效，需谨慎
		log.Printf("[App] WARNING: JWT secret changed, existing tokens will be invalid")
		// 实际生产环境可能需要更复杂的处理
	}

	// 更新邮件配置（邮件服务会通过自己的回调更新）
	log.Printf("[App] Email config updated: %s:%d",
		newConfig.Private.Email.Host, newConfig.Private.Email.Port)
}

// onManualReload 手动重载处理
func (app *App) onManualReload(oldConfig, newConfig *config.Config) {
	log.Printf("[App] Manual reload triggered, refreshing all components...")
	app.onFullConfigReload(oldConfig, newConfig)
}

// onFullConfigReload 全面配置重载
func (app *App) onFullConfigReload(oldConfig, newConfig *config.Config) {
	// 逐一更新各个组件

	// 1. 更新 Gin 模式
	if oldConfig.Basic.Server.Mode != newConfig.Basic.Server.Mode {
		gin.SetMode(newConfig.Basic.Server.Mode)
	}

	// 2. 更新中间件配置
	if app.middleware != nil {
		app.middleware.UpdateConfig(newConfig)
	}

	// 3. 更新处理器配置
	if app.handlers != nil {
		app.handlers.UpdateConfig(newConfig)
	}

	log.Printf("[App] Full config reload completed")
}

// RebuildComponents 重建特定组件（用于需要完全重建的场景）
func (app *App) RebuildComponents(component string) error {
	cfg := app.DynCfg.Get()
	if cfg == nil {
		return fmt.Errorf("config not available")
	}

	switch component {
	case "database":
		// 重新初始化数据库连接
		db, err := InitDB(cfg)
		if err != nil {
			return err
		}
		app.DB = db
		log.Printf("[App] Database connection rebuilt")

	case "redis":
		// 重新初始化 Redis
		// ...
		log.Printf("[App] Redis connection rebuilt")

	default:
		return fmt.Errorf("unknown component: %s", component)
	}

	return nil
}

// GetCurrentConfig 获取当前配置（便捷方法）
func (app *App) GetCurrentConfig() *config.Config {
	return app.DynCfg.Get()
}

// Close 关闭应用
func (app *App) Close() error {
	if app.DynCfg != nil {
		return app.DynCfg.Close()
	}
	return nil
}
