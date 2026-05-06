package vo

import "time"

// ReportVO 举报记录脱敏视图（对外暴露）
type ReportVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ReporterID uint       `json:"reporter_id,omitempty"` // 若 IsAnonymous=true 则不返回或返回0
	TargetID   uint       `json:"target_id"`
	TargetType string     `json:"target_type"`
	Type       string     `json:"type"` // ReportType 映射为字符串
	Reason     string     `json:"reason"`
	Status     string     `json:"status"` // ReportStatus 映射为字符串
	HandlerID  *uint      `json:"handler_id,omitempty"`
	HandleNote string     `json:"handle_note,omitempty"`
	HandleAt   *time.Time `json:"handle_at,omitempty"`

	IsAnonymous bool `json:"is_anonymous"`
	Priority    int8 `json:"priority"`

	// 脱敏后的举报人（仅在非匿名时返回）
	Reporter struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	} `json:"reporter,omitempty"`

	// 脱敏后的处理人（如果有）
	Handler struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	} `json:"handler,omitempty"`
}
