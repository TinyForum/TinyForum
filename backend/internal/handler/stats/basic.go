package stats

import (
	statsService "tiny-forum/internal/service/stats"
	"tiny-forum/pkg/utils"
)

type StatsHandler struct {
	statsSvc    statsService.StatsService
	timeHelpers *utils.TimeHelpers // 解析时间范围，如 "last7days"
	// timeParser  *utils.TimeParser      // 解析单个时间表达式，如 "2025-01-01" 或 "today"
	// rangeParser *utils.TimeRangeParser // 解析时间范围表达式，如 start=last7days&end=today
}

func NewStatsHandler(svc statsService.StatsService, timeHelpers *utils.TimeHelpers) *StatsHandler {
	return &StatsHandler{
		statsSvc:    svc,
		timeHelpers: timeHelpers,
	}
}
