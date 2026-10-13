package vo

import "time"

// ModeratorPostVO 版主面板的帖子视图对象
type ModeratorPostVO struct {
	ID           uint      `json:"id"`                     // 帖子ID
	Title        string    `json:"title"`                  // 标题
	Content      string    `json:"content"`                // 内容摘要（截断）
	AuthorID     uint      `json:"author_id"`              // 作者ID
	AuthorName   string    `json:"author_name"`            // 作者名称
	AuthorAvatar string    `json:"author_avatar,omitempty"` // 作者头像
	BoardID      uint      `json:"board_id"`               // 所属板块ID
	IsPinned     bool      `json:"is_pinned"`              // 是否板块内置顶
	ViewCount    int       `json:"view_count"`             // 浏览量
	LikeCount    int       `json:"like_count"`             // 点赞数
	CommentCount int64     `json:"comment_count"`          // 评论数
	Status       string    `json:"status"`                 // 内容状态：pending/published/deleted
	CreatedAt    time.Time `json:"created_at"`             // 创建时间
}

// ModeratorReportVO 版主面板的举报视图对象
type ModeratorReportVO struct {
	ID            uint      `json:"id"`                       // 举报ID
	TargetType    string    `json:"target_type"`              // 举报目标类型：post/comment
	TargetID      uint      `json:"target_id"`                // 举报目标ID
	TargetTitle   string    `json:"target_title"`             // 目标标题
	TargetContent string    `json:"target_content"`           // 目标内容摘要（截断）
	ReporterID    uint      `json:"reporter_id"`              // 举报人ID
	ReporterName  string    `json:"reporter_name"`            // 举报人名称
	Reason        string    `json:"reason"`                   // 举报原因
	Status        string    `json:"status"`                   // 举报状态：pending/resolved/rejected
	CreatedAt     time.Time `json:"created_at"`               // 举报时间
}

// BannedUserVO 版主面板的封禁用户视图对象
type BannedUserVO struct {
	ID        uint       `json:"id"`               // 封禁记录ID
	UserID    uint       `json:"user_id"`          // 被禁用户ID
	Username  string     `json:"username"`         // 被禁用户名称
	AvatarURL string     `json:"avatar_url,omitempty"` // 被禁用户头像
	BoardID   uint       `json:"board_id"`         // 封禁所属板块ID
	Reason    string     `json:"reason"`           // 封禁原因
	ExpiresAt *time.Time `json:"expires_at,omitempty"` // 解封时间（nil 表示永久）
	IsActive  bool       `json:"is_active"`        // 是否仍在封禁期内
	CreatedAt time.Time  `json:"created_at"`       // 封禁时间
}