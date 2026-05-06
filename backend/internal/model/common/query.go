package common

type PageQuery[T any] struct {
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Data     T        `json:"data"`
	Keyword  string   `json:"keyword"`
	SortBy   string   `json:"sort_by"`
	TagNames []string `json:"tag_names"`
}
