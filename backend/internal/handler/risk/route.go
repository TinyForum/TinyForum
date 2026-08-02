package risk

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func (h *RiskHandler) RegisterRoutes(admin *gin.RouterGroup, mw middleware.MiddlewareSet) {
	g := admin.Group("/risk")
	{
		g.GET("/audit/logs", h.ListAuditLogs)
	}
}
