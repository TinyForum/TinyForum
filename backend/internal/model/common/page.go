package common

type PageParam struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type PageResult[T any] struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	List     []T   `json:"list"`
	HasMore  bool  `json:"has_more"`
}
