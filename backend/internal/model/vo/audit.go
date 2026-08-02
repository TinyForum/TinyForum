package vo

import "time"

// AuditLogVO 审核日志脱敏视图（排除 operator_ip 等敏感字段）
type AuditLogVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Operator  string    `json:"operator"`
	Detail    string    `json:"detail,omitempty"`
}

type AuditLogsVO struct {
	Log []AuditLogVO `json:"log"`
}
