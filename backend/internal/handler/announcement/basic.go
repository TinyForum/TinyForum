package announcement

import (
	"tiny-forum/internal/service/announcement"
)

type AnnouncementHandler struct {
	service announcement.AnnouncementService
}

func NewAnnouncementHandler(svc announcement.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{service: svc}
}
