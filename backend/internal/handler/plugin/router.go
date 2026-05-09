package plugin

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册插件相关路由
func (h *Handler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// 插件列表（可公开，也可认证，按需）
	api.GET("/plugins", h.ListPlugins) // GET /api/v1/plugins

	// 需要认证的插件操作
	pluginGroup := api.Group("/plugins")
	pluginGroup.Use(mw.Auth())
	{
		// 上传安装插件（需认证）
		pluginGroup.POST("/upload", h.UploadPlugin) // POST /api/v1/plugins/upload

		// 获取当前用户已安装的插件列表
		api.GET("/users/me/plugins", h.ListMyPlugins) // GET /api/v1/users/me/plugins
	}
}