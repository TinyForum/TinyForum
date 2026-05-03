// internal/wire/routes.go
//
// 变更说明：
//   - adminGroup 改用 mw.CasbinAuth() 替代 mw.AdminRequired()
//   - 公开路由组新增 mw.OptionalAuth() + mw.CasbinAuth() 以支持 guest 策略
//   - 私有路由组统一用 mw.Auth() + mw.CasbinAuth()
//
// 中间件叠加顺序约定：
//   1. mw.Auth() / mw.OptionalAuth()  —— 解析 JWT，注入 user_role
//   2. mw.CasbinAuth()                —— 路由级 RBAC（读取 user_role 做决策）
//   3. mw.RateLimit(...)              —— 限流（读取 user_id）
//   4. mw.ContentCheck(...)           —— 内容安全（读取 request body）
//   5. mw.ModeratorRequired(...)      —— 版主细粒度权限（查数据库）
//
// 版主细粒度权限（CanDeletePost / CanBanUser 等）仍挂在具体路由上，
// Casbin 只做"该角色能否访问这条路由"的粗粒度决策。

package wire

import (
	"fmt"
	"tiny-forum/config"
	"tiny-forum/internal/infra/ratelimit"
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
		// 健康检查：完全公开，无需鉴权
		api.GET("/health", func(c *gin.Context) { c.JSON(200, "pong") })
	}

	// ── 公开路由（guest 可读，已登录用户可写）────────────────────────────────
	// OptionalAuth：有 token 则注入角色，没有 token 则注入 guest
	// CasbinAuth：根据角色决策，guest 只能 GET，user 可以 POST/PUT/DELETE
	handlers.Auth.RegisterRoutes(api, mw)
	handlers.Tag.RegisterRoutes(api, mw)
	handlers.Post.RegisterRoutes(api, mw)
	handlers.Comment.RegisterRoutes(api, mw)
	handlers.User.RegisterRoutes(api, mw)
	handlers.Notification.RegisterRoutes(api, mw)
	handlers.Board.RegisterRoutes(api, mw, repos.Board)
	handlers.Timeline.RegisterRoutes(api, mw)
	handlers.Topic.RegisterRoutes(api, mw)
	handlers.Answer.RegisterRoutes(api, mw)
	handlers.Question.RegisterRoutes(api, mw)
	handlers.Announcement.RegisterRoutes(api, mw)
	handlers.Stats.RegisterRoutes(api, mw)
	handlers.Admin.RegisterRoutes(api, mw)
	handlers.Upload.RegisterRoutes(api, mw)

	// ── Admin 路由组（示例：Casbin 替代 AdminRequired）───────────────────────
	//
	// 原写法：
	//   adminGroup := api.Group("/admin", mw.Auth(), mw.AdminRequired())
	//
	// 新写法：Auth() 注入角色，CasbinAuth() 查策略表决策
	// 策略表中 admin 拥有 /api/v1/admin/* 的所有方法权限，效果等价
	// adminGroup := api.Group("/admin", mw.Auth(), mw.CasbinAuth())
	// {
	// 	adminGroup.PUT("/users/:id/role", handlers.User.AdminSetRole)
	// 	adminGroup.POST("/users/:id/reset-password", handlers.User.AdminResetUserPassword)
	// 	adminGroup.GET("/users/score", handlers.User.AdminGetUserScore)
	// 	adminGroup.PUT("/users/:id/score", handlers.User.AdminSetScore)
	// 	adminGroup.GET("/boards/applications", handlers.Board.ListApplications)
	// 	adminGroup.POST("/boards/applications/:application_id/review", handlers.Board.ReviewApplication)
	// 	adminGroup.GET("/boards", handlers.Board.List)
	// 	adminGroup.GET("/posts", handlers.Post.AdminList)
	// 	adminGroup.GET("/posts/pending", handlers.Post.AdminGetModerationRequire)
	// 	adminGroup.PUT("/audit/tasks/:id/approve", handlers.Post.AdminApprovePost)
	// 	adminGroup.PUT("/audit/tasks/:id/reject", handlers.Post.AdminRejectPost)
	// 	adminGroup.PUT("/posts/:id/pin", handlers.Post.AdminTogglePin)

	// 	handlers.Risk.RegisterRoutes(adminGroup)
	// }

	// ── 版主路由示例（Casbin 粗粒度 + moderator.go 细粒度）────────────────────
	//
	// 第一层：CasbinAuth 确认"版主角色能访问此路由"
	// 第二层：CanBanUser / CanDeletePost 确认"该版主在此板块有对应权限"
	//
	// 路由注册在各 handler 的 RegisterRoutes 内部，此处仅作示意：
	//
	//   boardMod := api.Group("/boards/:id",
	//       mw.Auth(),
	//       mw.CasbinAuth(),                        // 第一层
	//   )
	//   boardMod.DELETE("/posts/:post_id",
	//       mw.CanDeletePost(repos.Board),           // 第二层
	//       handlers.Board.ModeratorDeletePost,
	//   )
	//   boardMod.POST("/ban",
	//       mw.CanBanUser(repos.Board),              // 第二层
	//       handlers.Board.BanUser,
	//   )

	// ── 限流示例（保持原有用法不变）──────────────────────────────────────────
	//
	// 限流中间件挂在具体路由上，与 Casbin 无冲突：
	//
	//   api.POST("/posts",
	//       mw.Auth(),
	//       mw.CasbinAuth(),
	//       mw.RateLimit(ratelimit.ActionCreatePost),
	//       mw.ContentCheck([]string{"title", "content"}),
	//       handlers.Post.Create,
	//   )
	_ = ratelimit.ActionCreatePost // 消除 unused import 提示（实际在 handler 内使用）
}
