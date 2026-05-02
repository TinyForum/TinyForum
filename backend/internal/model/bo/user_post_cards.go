package bo

import "tiny-forum/internal/model/do"

type UserPostByUserIdQuery struct {
	Page             int                 // 页码
	PageSize         int                 // 每页条数
	SortBy           string              // 排序字段
	Order            string              // 排序方向
	UserID           uint                // 用户ID
	Keyword          string              // 帖子标题或内容关键词
	Status           string              // 用户感知帖子状态
	ModerationStatus do.ModerationStatus // 风控状态
	Tag              string              // 标签
	BoardNmae        string              // 板块名称

}

// Normalize 设置默认值、校验
func (q *UserPostByUserIdQuery) Normalize() {
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

func (q *UserPostByUserIdQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

func (q *UserPostByUserIdQuery) Limit() int {
	return q.PageSize
}

// OrderClause 生成 ORDER BY 子句（需防注入：只允许预设的 SortBy 值）
func (q *UserPostByUserIdQuery) OrderClause() string {
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
