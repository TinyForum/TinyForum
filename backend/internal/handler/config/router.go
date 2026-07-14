package config

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册插件相关路由
func (h *ConfigHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {

	// 需要认证的插件操作
	adminGroup := api.Group("/admin")
	adminGroup.Use(mw.Auth())
	adminGroup.Use(mw.SystemMaintainerRequired())

	// 配置管理
	configGroup := adminGroup.Group("/config")
	{
		configGroup.GET("/list", h.ListConfigs)        // 列出配置文件
		configGroup.GET("/:file", h.GetConfig)         // 获取单个配置内容
		configGroup.GET("/:file/kv", h.GetConfigKV)    // 获取单个配置内容（KV 格式）
		configGroup.PUT("/:file", h.UpdateConfig)      // 更新单个配置内容
		configGroup.PUT("/:file/kv", h.UpdateConfigKV) // 更新单个配置内容（KV 格式）
		configGroup.POST("/reload", h.ReloadConfig)    // 重新加载配置
		configGroup.GET("/history", h.GetHistory)      // 获取配置历史记录
	}
}
