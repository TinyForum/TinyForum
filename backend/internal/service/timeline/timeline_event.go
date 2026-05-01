package timeline

import (
	"encoding/json"
	"tiny-forum/internal/model/do"
)

type CreateEventInput struct {
	UserID     uint
	ActorID    uint
	Action     do.ActionType
	TargetID   uint
	TargetType string
	Payload    interface{}
	Score      int
}

func (s *timelineService) CreateEvent(input CreateEventInput) error {
	payloadJSON, _ := json.Marshal(input.Payload)

	event := &do.TimelineEvent{
		UserID:     input.UserID,
		ActorID:    input.ActorID,
		Action:     input.Action,
		TargetID:   input.TargetID,
		TargetType: input.TargetType,
		Payload:    string(payloadJSON),
		Score:      input.Score,
	}

	return s.timelineRepo.CreateEvent(event)
}
