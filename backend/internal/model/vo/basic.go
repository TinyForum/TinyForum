package vo

// Response 统一响应结构
type BasicResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	Timestamp int64  `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
}

// PageData 分页数据结构
type BasicPageData struct {
	List     any   `json:"list"`      // 数据列表
	Total    int64 `json:"total"`     // 数据总数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页数量
	HasMore  bool  `json:"has_more"`  // 是否有更多数据
}
