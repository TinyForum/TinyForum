package vo

import (
	"time"
)

type ReportVO struct {
	ID              uint       `json:"id"`               // 举报记录ID
	TargetID        uint       `json:"target_id"`        // 举报目标ID
	TargetType      string     `json:"target_type"`      // 举报目标类型
	Type            string     `json:"type"`             // ReportType 映射为字符串
	Reason          string     `json:"reason"`           // 举报原因
	Status          string     `json:"status"`           // 举报状态
	HandleNote      string     `json:"handle_note"`      // 处理备注
	HandleAt        *time.Time `json:"handle_at"`        // 处理时间
	ContentSnapshot string     `json:"content_snapshot"` // 按规则脱敏
	IsAnonymous     bool       `json:"is_anonymous"`     // 是否匿名
	Priority        int8       `json:"priority"`         // 优先级
	// 脱敏后的举报人（可能为匿名结构）
	Reporter *UserPublicVO `json:"reporter,omitempty"` // 若 IsAnonymous=true 则不返回或返回空
	Handler  *UserPublicVO `json:"handler,omitempty"`  // 若有处理人则返回
}
