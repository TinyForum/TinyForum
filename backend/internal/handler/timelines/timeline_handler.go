package timeline

import (
	timelineService "tiny-forum/internal/service/timeline"
)

type TimelineHandler struct {
	timelineSvc *timelineService.TimelineService
}

func NewTimelineHandler(timelineSvc *timelineService.TimelineService) *TimelineHandler {
	return &TimelineHandler{timelineSvc: timelineSvc}
}
