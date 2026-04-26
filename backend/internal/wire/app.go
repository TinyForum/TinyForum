package wire

import (
	"tiny-forum/config"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
	Cfg    *config.Config
}

// InitApp 是原来的入口函数，现在调用各个模块进行组装
func InitApp(cfg *config.Config) (*App, error) {
	// 1. 数据库
	db, err := InitDB(cfg)
	if err != nil {
		return nil, err
	}

	// 2. JWT 管理器
	jwtMgr := jwtpkg.NewJWTManager(cfg.Private.JWT.Secret, cfg.Private.JWT.Expire)

	// 3. 基础设施（Redis、限流、敏感词）
	infra, err := InitInfra(cfg)
	if err != nil {
		return nil, err
	}

	// 4. 数据仓库层
	repos := NewRepositories(db)

	// 5. 服务层
	services := NewServices(cfg, jwtMgr, repos, infra)

	// 辅助工具
	helpers := NewHelpers()

	// 6. 控制器层
	handlers := NewHandlers(services, helpers.TimeHelpers)

	// 7. 中间件工厂
	mw := NewMiddlewareSet(jwtMgr, db, services)

	// 8. Gin 引擎
	gin.SetMode(cfg.Basic.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 9. 注册路由（传入 repos 是因为部分中间件需要 boardRepo 动态判断）
	RegisterRoutes(engine, handlers, mw, repos, cfg)

	return &App{
		Engine: engine,
		DB:     db,
		Cfg:    cfg,
	}, nil
}
