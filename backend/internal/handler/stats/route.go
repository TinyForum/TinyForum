package stats

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *StatsHandler) RegisterRoutes(stats *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	g := stats.Group("/statistics")
	{
		g.GET("", mw.AuthMW(), mw.AdminRequiredMW(), h.GetStatsTotal)       // 获取所有统计指标
		g.GET("/day", mw.AuthMW(), mw.AdminRequiredMW(), h.GetStatsDay)     // 获取日数据
		g.GET("/total", mw.AuthMW(), mw.AdminRequiredMW(), h.GetStatsTotal) // 获取所有统计指标
		g.GET("/trend", mw.AuthMW(), mw.AdminRequiredMW(), h.GetStatsTrend) // 获取趋势指标
		g.GET("/range", mw.AuthMW(), mw.AdminRequiredMW(), h.GetStatsRange) // 获取指定时间范围内的数据
	}
}
