package stats

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *StatsHandler) RegisterRoutes(stats *gin.RouterGroup, mw middleware.MiddlewareSet) {
	g := stats.Group("/statistics")
	{
		g.GET("", mw.Auth(), mw.AdminRequired(), h.GetStatsTotal)       // 获取所有统计指标
		g.GET("/day", mw.Auth(), mw.AdminRequired(), h.GetStatsDay)     // 获取日数据
		g.GET("/total", mw.Auth(), mw.AdminRequired(), h.GetStatsTotal) // 获取所有统计指标
		g.GET("/trend", mw.Auth(), mw.AdminRequired(), h.GetStatsTrend) // 获取趋势指标
		g.GET("/range", mw.Auth(), mw.AdminRequired(), h.GetStatsRange) // 获取指定时间范围内的数据
	}
}
