package risk

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// 审计日志属于管理功能，必须登录且具备管理员权限才能访问
func (h *RiskHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	g := api.Group("/risk", mw.Auth(), mw.AdminRequired())
	{
		g.GET("/audit/logs", h.ListAuditLogs)
	}
}
