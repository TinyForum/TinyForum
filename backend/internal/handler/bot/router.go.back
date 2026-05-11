package bot

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// 需要认证的附件操作
	botGroup := api.Group("/bots")
	botGroup.Use(mw.Auth())
	{
		// bots.POST("", h.Create)
		botGroup.GET("", h.List)          // 获取机器人列表
		botGroup.POST("", h.Create)       // 创建机器人
		botGroup.PUT("/:id", h.Update)    // 更新机器人
		botGroup.DELETE("/:id", h.Delete) // 删除机器人
		//
		botGroup.GET("/:id", h.Get)           // 获取机器人
		botGroup.GET("/user/me", h.ListMyBot) // 获取用户创建的机器人列表

		botGroup.POST("/:id/run", h.RunNow) // 立即运行机器人
	}
}
