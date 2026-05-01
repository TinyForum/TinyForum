package query

import (
	"time"
	"tiny-forum/internal/model/po"
)

type ListAnnouncements struct {
	Page      int                       `form:"page" binding:"min=1"`
	PageSize  int                       `form:"page_size" binding:"min=1,max=100"`
	BoardID   *uint                     `form:"board_id"`
	Type      *po.AnnouncementType   `form:"type"`
	Status    *po.AnnouncementStatus `form:"status"`
	IsPinned  *bool                     `form:"is_pinned"`
	IsGlobal  *bool                     `form:"is_global"`
	Keyword   string                    `form:"keyword"`
	StartTime *time.Time                `form:"start_time"`
	EndTime   *time.Time                `form:"end_time"`
}
