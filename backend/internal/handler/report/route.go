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
}
