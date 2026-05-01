package dto

import (
	"time"
	"tiny-forum/internal/model/po"
)

type CreateAnnouncementRequest struct {
	Title       string                `json:"title" binding:"required,min=1,max=200"`
	Content     string                `json:"content" binding:"required"`
	Summary     string                `json:"summary"`
	Cover       string                `json:"cover"`
	Type        po.AnnouncementType   `json:"type" binding:"required,oneof=normal important emergency event"`
	IsPinned    bool                  `json:"is_pinned"`
	IsGlobal    bool                  `json:"is_global"`
	Status      po.AnnouncementStatus `json:"status"`
	BoardID     *uint                 `json:"board_id"`
	PublishedAt *time.Time            `json:"published_at"`
	ExpiredAt   *time.Time            `json:"expired_at"`
}

type UpdateAnnouncementRequest struct {
	Title       *string              `json:"title"`
	Content     *string              `json:"content"`
	Summary     *string              `json:"summary"`
	Cover       *string              `json:"cover"`
	Type        *po.AnnouncementType `json:"type"`
	IsPinned    *bool                `json:"is_pinned"`
	IsGlobal    *bool                `json:"is_global"`
	BoardID     *uint                `json:"board_id"`
	PublishedAt *time.Time           `json:"published_at"`
	ExpiredAt   *time.Time           `json:"expired_at"`
}
