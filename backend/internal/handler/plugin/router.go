package plugin

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册插件相关路由
func (h *Handler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// api.GET("/plugins", h.ListPlugins) // GET /api/v1/plugins

	// 需要认证的插件操作
	pluginGroup := api.Group("/plugins")
	pluginGroup.Use(mw.Auth())
	{
		// 上传插件
		pluginGroup.POST("", h.UploadPlugin) // POST /api/v1/plugins 上传插件
		// 获取所有插件
		pluginGroup.GET("", h.ListPlugins) // GET /api/v1/plugins 获取所有插件
		// 删除插件
		pluginGroup.DELETE("/:id", h.DeletePlugin) // DELETE /api/v1/plugins/:id 删除插件

		// 获取当前用户的插件列表
		pluginGroup.GET("/me", h.ListMyPlugins)          // GET /api/v1/plugins/me 获取当前用户的插件列表
		pluginGroup.PATCH("/:id/toggle", h.TogglePlugin) // PATCH /api/v1/plugins/toggle 切换插件启用状态

	}
}
