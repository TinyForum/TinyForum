package vo

import "time"

// BoardVO 版块脱敏视图（对外暴露）
type BoardVO struct {
	ID          uint      `json:"id"`                    // 版块ID
	CreatedAt   time.Time `json:"created_at"`            // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`            // 更新时间
	Name        string    `json:"name"`                  // 版块名称
	Slug        string    `json:"slug"`                  // 版块别名
	Description string    `json:"description,omitempty"` // 版块描述
	Icon        string    `json:"icon,omitempty"`        // 版块图标
	Cover       string    `json:"cover,omitempty"`       // 版块封面
	ParentID    *uint     `json:"parent_id,omitempty"`   // 仅保留父版块ID，不嵌套完整对象
	SortOrder   int       `json:"sort_order"`            // 排序
	ViewRole    string    `json:"view_role"`             // UserRole 映射为字符串
	PostRole    string    `json:"post_role"`             // UserRole 映射为字符串
	ReplyRole   string    `json:"reply_role"`            // UserRole 映射为字符串
	PostCount   int       `json:"post_count"`            // 帖子数量
	ThreadCount int       `json:"thread_count"`          // 主题数量
	TodayCount  int       `json:"today_count"`           // 今日发帖数量
}

// HotBoardRow 热门板块查询结果行
type HotBoardRowVO struct {
	ID           int64
	Name         string
	Icon         string
	ArticleCount int64
	CommentCount int64
	ActiveUser   int64
}
