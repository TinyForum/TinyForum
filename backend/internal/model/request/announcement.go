package request

import (
	"time"
	"tiny-forum/internal/model/do"
)

type ListAnnouncements struct {
    Page      int                      `form:"page" binding:"min=1"`
    PageSize  int                      `form:"page_size" binding:"min=1,max=100"`
    BoardID   *uint                    `form:"board_id"`
    Type      *do.AnnouncementType     `form:"type" binding:"omitempty,oneof=0 1 2 3"`
    Status    *do.AnnouncementStatus   `form:"status" binding:"omitempty,oneof=0 1 2"`
    IsPinned  *bool                    `form:"is_pinned"`
    IsGlobal  *bool                    `form:"is_global"`
    Keyword   string                   `form:"keyword"`
    StartTime *time.Time               `form:"start_time"`
    EndTime   *time.Time               `form:"end_time"`
}
type CreateAnnouncement struct {
    Title       string                `json:"title" binding:"required,min=1,max=200"`
    Content     string                `json:"content" binding:"required"`
    Summary     string                `json:"summary"`
    Cover       string                `json:"cover"`
    IsPinned    bool                  `json:"is_pinned"`
    IsGlobal    bool                  `json:"is_global"`
    Type        *do.AnnouncementType   `json:"type" binding:"required,oneof=0 1 2 3"`
    Status      *do.AnnouncementStatus `json:"status"`  
    BoardID     *uint                 `json:"board_id"`
    PublishedAt *time.Time            `json:"published_at"`
    ExpiredAt   *time.Time            `json:"expired_at"`
}

type UpdateAnnouncement struct {
    Title       *string                `json:"title"`
    Content     *string                `json:"content"`
    Summary     *string                `json:"summary"`
    Cover       *string                `json:"cover"`
    Type        *do.AnnouncementType   `json:"type" binding:"omitempty,oneof=0 1 2 3"`
    Status      *do.AnnouncementStatus `json:"status" binding:"omitempty,oneof=0 1 2"`
    IsPinned    *bool                  `json:"is_pinned"`
    IsGlobal    *bool                  `json:"is_global"`
    BoardID     *uint                  `json:"board_id"`
    PublishedAt *time.Time             `json:"published_at"`
    ExpiredAt   *time.Time             `json:"expired_at"`
}