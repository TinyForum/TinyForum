package timeline

import (
	"encoding/json"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"

	"gorm.io/datatypes"
)

func (s *timelineService) CreateEvent(input request.CreateEventRequest) error {
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
