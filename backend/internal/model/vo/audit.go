package vo

import "tiny-forum/internal/model/do"

type AuditLogsVO struct {
	Log []do.AuditLog `json:"log"`
}
