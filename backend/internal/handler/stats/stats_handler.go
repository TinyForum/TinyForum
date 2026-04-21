package stats

import (
	statsService "tiny-forum/internal/service/stats"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsSvc *statsService.StatsService
}

func NewStatsHandler(svc *statsService.StatsService) *StatsHandler {
	return &StatsHandler{statsSvc: svc}
}

func (h *StatsHandler) RegisterRoutes(stats *gin.RouterGroup) {
	g := stats.Group("/statistics")
	{
		g.GET("/day", h.GetStatsDay)     // 获取日数据
		g.GET("/total", h.GetStatsTotal) // 获取所有统计指标
		g.GET("/trend", h.GetStatsTrend) // 获取趋势指标
		g.GET("/range", h.GetStatsRange) // 获取指定时间范围内的数据
	}
}
