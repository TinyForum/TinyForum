package common

type PageQuery[T any] struct {
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Cursor   string   `json:"cursor"` // 游标分页（created_at 时间戳），与 Page 互斥
	Data     T        `json:"data"`
	Keyword  string   `json:"keyword"`
	SortBy   string   `json:"sort_by"`
	TagNames []string `json:"tag_names"`
	Order    string   `json:"order"`
}
