package vo

import "tiny-forum/internal/model/po"

type ListAnnouncements struct {
	Total         int64             `json:"total"`
	Page          int               `json:"page"`
	PageSize      int               `json:"page_size"`
	Announcements []po.Announcement `json:"announcements"`
}
