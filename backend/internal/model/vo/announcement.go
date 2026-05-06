package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

type ListAnnouncements struct {
	Total         int64             `json:"total"`
	Page          int               `json:"page"`
	PageSize      int               `json:"page_size"`
	Announcements []do.Announcement `json:"announcements"`
}

// AnnouncementVO 公告脱敏视图
type AnnouncementVO struct {
	ID          uint                   `json:"id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"` // 内容可能需要全部展示
	Summary     string                 `json:"summary,omitempty"`
	Cover       string                 `json:"cover,omitempty"`
	Type        *do.AnnouncementType   `json:"type,omitempty"`
	Status      *do.AnnouncementStatus `json:"status,omitempty"`
	IsPinned    bool                   `json:"is_pinned"`
	IsGlobal    bool                   `json:"is_global"`
	BoardID     *uint                  `json:"board_id,omitempty"` // 只保留 ID，不暴露完整板块
	PublishedAt *time.Time             `json:"published_at,omitempty"`
	ExpiredAt   *time.Time             `json:"expired_at,omitempty"`
	ViewCount   int                    `json:"view_count"`
	CreatedBy   uint                   `json:"created_by"` // 创建人 ID

	// 脱敏后的创建人信息（可选，如果需要展示用户名）
	Creator struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar,omitempty"`
	} `json:"creator,omitempty"`

	// 脱敏后的板块信息（可选，如果前端需要展示板块名称）
	Board struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"board,omitempty"`
}
