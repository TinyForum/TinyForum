package board

import (
	"time"
)

type ApplicationStatusDetail struct {
	HasApplication bool            `json:"has_application"`
	ApplicationID  uint            `json:"application_id,omitempty"`
	Status         string          `json:"status,omitempty"`
	Reason         string          `json:"reason,omitempty"`
	CreatedAt      time.Time       `json:"created_at,omitempty"`
	ReviewNote     string          `json:"review_note,omitempty"`
	ReviewerID     uint            `json:"reviewer_id,omitempty"`
	ReviewedAt     *time.Time      `json:"reviewed_at,omitempty"`
	CanCancel      bool            `json:"can_cancel"`
	CanResubmit    bool            `json:"can_resubmit"`
	RequestedPerms map[string]bool `json:"requested_perms,omitempty"`
	CanApply       bool            `json:"can_apply"`
}

// type ApplyModeratorInput struct {
// 	UserID               uint                     `json:"user_id"`               // 申请人ID
// 	BoardID              uint                     `json:"board_id"`              // 版块ID
// 	Reason               string                   `json:"reason"`                // 申请理由
// 	RequestedPermissions []do.ModeratorPermission `json:"requested_permissions"` // 请求的权限列表（外部传入）
// }
