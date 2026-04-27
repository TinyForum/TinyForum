package dto

import (
	"time"
	"tiny-forum/internal/model"
)

type CreateAnnouncementRequest struct {
	Title       string                   `json:"title" binding:"required,min=1,max=200"`
	Content     string                   `json:"content" binding:"required"`
	Summary     string                   `json:"summary"`
	Cover       string                   `json:"cover"`
	Type        model.AnnouncementType   `json:"type" binding:"required,oneof=normal important emergency event"`
	IsPinned    bool                     `json:"is_pinned"`
	IsGlobal    bool                     `json:"is_global"`
	Status      model.AnnouncementStatus `json:"status"`
	BoardID     *uint                    `json:"board_id"`
	PublishedAt *time.Time               `json:"published_at"`
	ExpiredAt   *time.Time               `json:"expired_at"`
}

type UpdateAnnouncementRequest struct {
	Title       *string                 `json:"title"`
	Content     *string                 `json:"content"`
	Summary     *string                 `json:"summary"`
	Cover       *string                 `json:"cover"`
	Type        *model.AnnouncementType `json:"type"`
	IsPinned    *bool                   `json:"is_pinned"`
	IsGlobal    *bool                   `json:"is_global"`
	BoardID     *uint                   `json:"board_id"`
	PublishedAt *time.Time              `json:"published_at"`
	ExpiredAt   *time.Time              `json:"expired_at"`
}

type ListAnnouncementRequest struct {
	Page      int                       `form:"page" binding:"min=1"`
	PageSize  int                       `form:"page_size" binding:"min=1,max=100"`
	BoardID   *uint                     `form:"board_id"`
	Type      *model.AnnouncementType   `form:"type"`
	Status    *model.AnnouncementStatus `form:"status"`
	IsPinned  *bool                     `form:"is_pinned"`
	IsGlobal  *bool                     `form:"is_global"`
	Keyword   string                    `form:"keyword"`
	StartTime *time.Time                `form:"start_time"`
	EndTime   *time.Time                `form:"end_time"`
}

type ListAnnouncementResponse struct {
	Total         int64                `json:"total"`
	Page          int                  `json:"page"`
	PageSize      int                  `json:"page_size"`
	Announcements []model.Announcement `json:"announcements"`
}
