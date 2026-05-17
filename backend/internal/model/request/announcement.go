package request

import (
	"fmt"
	"time"
	"tiny-forum/internal/model/do"
)

type ListAnnouncementsRequest struct {
	Page      int                    `form:"page" binding:"min=1"`
	PageSize  int                    `form:"page_size" binding:"min=1,max=100"`
	BoardID   *uint                  `form:"board_id"`
	Type      *do.AnnouncementType   `form:"type" binding:"omitempty"`
	Status    *do.AnnouncementStatus `form:"status" binding:"omitempty"`
	IsPinned  *bool                  `form:"is_pinned"`
	IsGlobal  *bool                  `form:"is_global"`
	Keyword   string                 `form:"keyword"`
	StartTime *time.Time             `form:"start_time"`
	EndTime   *time.Time             `form:"end_time"`
}

// 检查格式：
func (l *ListAnnouncementsRequest) Validate() error {
	// 默认分页参数
	if l.Page <= 0 {
		l.Page = 1
	}
	if l.PageSize <= 0 {
		l.PageSize = 20 // 默认每页20条
	}
	if l.PageSize > 100 {
		l.PageSize = 100
	}

	// 可选：校验状态枚举值（如果仍用 string 类型）
	if l.Status != nil && !l.Status.IsValid() {
		return fmt.Errorf("invalid announcement status: %s", *l.Status)
	}
	if l.Type != nil && !l.Type.IsValid() {
		return fmt.Errorf("invalid announcement type: %s", *l.Type)
	}

	// 时间范围校验
	if l.StartTime != nil && l.EndTime != nil && l.StartTime.After(*l.EndTime) {
		return fmt.Errorf("start_time cannot be after end_time")
	}

	return nil
}

type CreateAnnouncement struct {
	Title       string                 `json:"title" binding:"required,min=1,max=200"`
	Content     string                 `json:"content" binding:"required"`
	Summary     string                 `json:"summary"`
	Cover       string                 `json:"cover"`
	IsPinned    bool                   `json:"is_pinned"`
	IsGlobal    bool                   `json:"is_global"`
	Type        *do.AnnouncementType   `json:"type" binding:"required"`
	Status      *do.AnnouncementStatus `json:"status"`
	BoardID     *uint                  `json:"board_id"`
	PublishedAt *time.Time             `json:"published_at"`
	ExpiredAt   *time.Time             `json:"expired_at"`
}

type UpdateAnnouncement struct {
	Title       *string                `json:"title"`
	Content     *string                `json:"content"`
	Summary     *string                `json:"summary"`
	Cover       *string                `json:"cover"`
	Type        *do.AnnouncementType   `json:"type" binding:"omitempty"`
	Status      *do.AnnouncementStatus `json:"status" binding:"omitempty"`
	IsPinned    *bool                  `json:"is_pinned"`
	IsGlobal    *bool                  `json:"is_global"`
	BoardID     *uint                  `json:"board_id"`
	PublishedAt *time.Time             `json:"published_at"`
	ExpiredAt   *time.Time             `json:"expired_at"`
}
