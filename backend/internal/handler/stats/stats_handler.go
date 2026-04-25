package stats

import (
	statsService "tiny-forum/internal/service/stats"
	"tiny-forum/pkg/utils"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsSvc    *statsService.StatsService
	timeParser  *utils.TimeParser      // 解析单个时间表达式，如 "2025-01-01" 或 "today"
	rangeParser *utils.TimeRangeParser // 解析时间范围表达式，如 start=last7days&end=today
}

func NewStatsHandler(svc *statsService.StatsService, timeParser *utils.TimeParser, rangeParser *utils.TimeRangeParser) *StatsHandler {
	return &StatsHandler{
		statsSvc:    svc,
		timeParser:  timeParser,
		rangeParser: rangeParser,
	}
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
