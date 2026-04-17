package stats

import (
	statsService "tiny-forum/internal/service/stats"
)

type StatsHandler struct {
	statsSvc *statsService.StatsService
}

func NewStatsHandler(svc *statsService.StatsService) *StatsHandler {
	return &StatsHandler{statsSvc: svc}
}
