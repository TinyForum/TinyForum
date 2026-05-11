package bot

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	g := api.Group("/bots")
	g.Use(mw.Auth())
	{
		// ── CRUD ──────────────────────────────────────────────────────
		g.GET("", h.List)          // 获取所有机器人列表
		g.POST("", h.Create)       // 创建机器人（Lua 脚本 or 零代码）
		g.PUT("/:id", h.Update)    // 更新机器人
		g.DELETE("/:id", h.Delete) // 删除机器人
		g.GET("/:id", h.Get)       // 获取机器人详情

		// ── 用户维度 ──────────────────────────────────────────────────
		// 注意：/user/me 必须在 /:id 之前注册，否则 gin 会把 "user" 当 id 匹配
		g.GET("/user/me", h.ListMyBot) // 我创建的机器人

		// ── 执行 ──────────────────────────────────────────────────────
		g.POST("/:id/run", h.RunNow) // 手动触发执行

		// ── 零代码支持 ────────────────────────────────────────────────
		nocodeg := g.Group("/nocode")
		{
			nocodeg.GET("/metadata", h.GetNocodeMetadata) // 前端加载节点定义
			nocodeg.POST("/validate", h.ValidateFlow)     // 提交前校验 Flow
		}
	}
}
