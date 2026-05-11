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

// // PageData 分页数据结构
// type BasicPageData struct {
// 	List     any   `json:"list"`      // 数据列表
// 	Total    int64 `json:"total"`     // 数据总数
// 	Page     int   `json:"page"`      // 当前页码
// 	PageSize int   `json:"page_size"` // 每页数量
// 	HasMore  bool  `json:"has_more"`  // 是否有更多数据
// }
