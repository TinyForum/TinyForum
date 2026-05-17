package request

import "tiny-forum/internal/model/do"

type CreateEventRequest struct {
	UserID     uint
	ActorID    uint
	Action     do.ActionType
	TargetID   uint
	TargetType string
	Payload    interface{}
	Score      int
}
