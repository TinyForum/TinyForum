package bo

type PageQuery[T any] struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Options  []T `json:"options"`
}
