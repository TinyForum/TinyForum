package reviewer

import (
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/model/do"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册审核员相关路由
// 鉴权：登录 + 角色优先级不低于 reviewer（reviewer/moderator/admin/super_admin）
func (h *ReviewerHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	g := api.Group("/reviewer", mw.Auth(), middleware.RequireRole(do.RoleReviewer))
	{
		g.GET("/audit/tasks", h.ListPendingTasks)            // GET /api/v1/reviewer/audit/tasks 待审核任务
		g.GET("/audit/history", h.ListReviewedTasks)         // GET /api/v1/reviewer/audit/history 审核历史
		g.PUT("/audit/tasks/:id/approve", h.ApproveTask)     // PUT /api/v1/reviewer/audit/tasks/:id/approve 审核通过
		g.PUT("/audit/tasks/:id/reject", h.RejectTask)       // PUT /api/v1/reviewer/audit/tasks/:id/reject 审核拒绝
	}
}
