package stats

import (
	statsService "tiny-forum/internal/service/stats"
	"tiny-forum/pkg/timeutil"
)

type StatsHandler struct {
	statsSvc    statsService.StatsService
	timeHelpers *timeutil.TimeHelpers // 解析时间范围，如 "last7days"
	// timeParser  *utils.TimeParser      // 解析单个时间表达式，如 "2025-01-01" 或 "today"
	// rangeParser *utils.TimeRangeParser // 解析时间范围表达式，如 start=last7days&end=today
}

func NewStatsHandler(svc statsService.StatsService, timeHelpers *timeutil.TimeHelpers) *StatsHandler {
	return &StatsHandler{
		statsSvc:    svc,
		timeHelpers: timeHelpers,
	}
}
