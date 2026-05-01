package vo

import "tiny-forum/internal/model/do"

type ListAnnouncements struct {
	Total         int64             `json:"total"`
	Page          int               `json:"page"`
	PageSize      int               `json:"page_size"`
	Announcements []do.Announcement `json:"announcements"`
}
