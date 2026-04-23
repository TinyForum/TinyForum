package risk

// TODO: Refactory
import (
	riskservice "tiny-forum/internal/service/risk"

	"github.com/gin-gonic/gin"
)

type RiskHandler struct {
	checkSvc *riskservice.ContentCheckService
	riskSvc  *riskservice.RiskService
}

func NewRiskHandler(checkSvc *riskservice.ContentCheckService, riskSvc *riskservice.RiskService) *RiskHandler {
	return &RiskHandler{checkSvc: checkSvc, riskSvc: riskSvc}
}

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
