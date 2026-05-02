// bo/user_post_query.go
package bo

type UserPostQuery struct {
	// 分页参数
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"` // created_at, hot, view_count 等
	Order    string `json:"order"`   // asc, desc

	// 过滤条件
	UserID           uint   `json:"user_id"`
	Keyword          string `json:"keyword"`
	Status           string `json:"status"`            // 用户感知状态（草稿、已发布等）
	ModerationStatus string `json:"moderation_status"` // 风控状态
	TagID            uint   `json:"tag_id"`
	PostType         string `json:"post_type"`

	// 其他控制
	NeedAuthor bool `json:"-"` // 是否预加载作者（按需加载）
	NeedTags   bool `json:"-"` // 是否预加载标签
}

// Normalize 设置默认值、校验
func (q *UserPostQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	if q.Order == "" {
		q.Order = "desc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.ModerationStatus == "" {
		q.ModerationStatus = "approved" // 默认只查审核通过的
	}
	if q.Status == "" {
		// 根据业务决定，例如只查 published
		q.Status = "published"
	}
}

func (q *UserPostQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

func (q *UserPostQuery) Limit() int {
	return q.PageSize
}

// OrderClause 生成 ORDER BY 子句（需防注入：只允许预设的 SortBy 值）
func (q *UserPostQuery) OrderClause() string {
	sortMap := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"view_count": "view_count",
		"like_count": "like_count",
		"hot":        "pin_top DESC, like_count DESC, view_count DESC",
	}
	sortExpr, ok := sortMap[q.SortBy]
	if !ok {
		sortExpr = "created_at"
	}
	// 普通字段附加方向（hot 自带排序方向）
	if q.SortBy != "hot" {
		if q.Order == "asc" {
			sortExpr += " ASC"
		} else {
			sortExpr += " DESC"
		}
	}
	return sortExpr
}
