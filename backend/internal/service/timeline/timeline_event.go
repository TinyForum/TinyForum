package timeline

import (
	"encoding/json"
	"tiny-forum/internal/model"
)

type CreateEventInput struct {
	UserID     uint
	ActorID    uint
	Action     model.ActionType
	TargetID   uint
	TargetType string
	Payload    interface{}
	Score      int
}

func (s *TimelineService) CreateEvent(input CreateEventInput) error {
	payloadJSON, _ := json.Marshal(input.Payload)

	event := &model.TimelineEvent{
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
