package request

import (
	"time"
	"tiny-forum/internal/model/do"
)

type ListAnnouncements struct {
	Page      int                    `form:"page" binding:"min=1"`
	PageSize  int                    `form:"page_size" binding:"min=1,max=100"`
	BoardID   *uint                  `form:"board_id"`
	Type      *do.AnnouncementType   `form:"type"`
	Status    *do.AnnouncementStatus `form:"status"`
	IsPinned  *bool                  `form:"is_pinned"`
	IsGlobal  *bool                  `form:"is_global"`
	Keyword   string                 `form:"keyword"`
	StartTime *time.Time             `form:"start_time"`
	EndTime   *time.Time             `form:"end_time"`
}
type CreateAnnouncement struct {
	Title       string                `json:"title" binding:"required,min=1,max=200"`                         // 标题
	Content     string                `json:"content" binding:"required"`                                     // 内容
	Summary     string                `json:"summary"`                                                        // 摘要
	Cover       string                `json:"cover"`                                                          //
	Type        do.AnnouncementType   `json:"type" binding:"required,oneof=normal important emergency event"` // 类型
	IsPinned    bool                  `json:"is_pinned"`                                                      // 是否置顶
	IsGlobal    bool                  `json:"is_global"`                                                      // 是否全局
	Status      do.AnnouncementStatus `json:"status"`                                                         // 状态
	BoardID     *uint                 `json:"board_id"`                                                       // 版块ID
	PublishedAt *time.Time            `json:"published_at"`                                                   // 发布时间
	ExpiredAt   *time.Time            `json:"expired_at"`                                                     // 过期时间
}

type UpdateAnnouncement struct {
	Title       *string              `json:"title"`
	Content     *string              `json:"content"`
	Summary     *string              `json:"summary"`
	Cover       *string              `json:"cover"`
	Type        *do.AnnouncementType `json:"type"`
	IsPinned    *bool                `json:"is_pinned"`
	IsGlobal    *bool                `json:"is_global"`
	BoardID     *uint                `json:"board_id"`
	PublishedAt *time.Time           `json:"published_at"`
	ExpiredAt   *time.Time           `json:"expired_at"`
}
