package risk

// TODO: Refactory
import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func (h *RiskHandler) RegisterRoutes(admin *gin.RouterGroup) {
	g := admin.Group("/risk")
	{
		// 审核队列(暂时没用，可以删除)
		// g.GET("/audit/tasks", h.ListAuditTasks)
		// g.POST("/audit/tasks/:id/approve", h.ApproveTask)
		// g.POST("/audit/tasks/:id/reject", h.RejectTask)

		// 操作日志 未使用，可以删除
		g.GET("/audit/logs", h.ListAuditLogs)
	}
}
