package timeline

import (
	"encoding/json"
	"tiny-forum/internal/model/do"

	"gorm.io/datatypes"
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
	targetType, err := do.ParseTilelineTargetType(input.TargetType)
	if err != nil {
		return err
	}

	event := &do.TimelineEvent{
		UserID:     input.UserID,
		ActorID:    input.ActorID,
		Action:     input.Action,
		TargetID:   input.TargetID,
		TargetType: targetType,
		Payload:    datatypes.JSON(payloadJSON),
		Score:      input.Score,
	}

	return s.timelineRepo.CreateEvent(event)
}
