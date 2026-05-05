package plugin

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *PluginHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// 需要认证的上传接口
	plugin := api.Group("/plugins")
	plugin.Use(mw.Auth())
	{
		plugin.GET("", h.List) // GET /api/v1/plugins/ - 获取插件列表
		// plugin.POST("", h.Upload)                 // POST /api/v1/upload - 上传文件
		// plugin.GET("/user/files", h.GetUserFiles) // GET /api/v1/upload/user/files - 获取用户文件列表
		// plugin.GET("/:file_id", h.GetFile)        // GET /api/v1/upload/:file_id - 获取文件信息
		// plugin.DELETE("/:file_id", h.DeleteFile)  // DELETE /api/v1/upload/:file_id - 删除文件
	}

}
