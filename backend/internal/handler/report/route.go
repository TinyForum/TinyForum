package report

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册举报相关路由
func (h *ReportHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	reportsGroup := api.Group("/reports")
	{
		reportsGroup.POST("", mw.Auth(), h.CreateReport) // 提交举报
	}

	// 管理员举报处理
	adminReports := api.Group("/admin/reports", mw.Auth(), mw.AdminRequired())
	{
		adminReports.GET("/:id", h.GetReport)          // GET /api/v1/admin/reports/:id 获取举报详情
		adminReports.POST("/:id/handle", h.HandleReport) // POST /api/v1/admin/reports/:id/handle 处理举报
	}
}
