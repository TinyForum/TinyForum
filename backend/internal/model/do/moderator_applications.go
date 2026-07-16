package do

import (
	"fmt"
	"time"
	"tiny-forum/internal/model/common"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/datatypes"
)

// ── ModeratorApplication 版主申请表 ─────────────────────────────────────────

// ApplicationStatus 申请状态
type ApplicationStatus string

const (
	ApplicationPending  ApplicationStatus = "pending"  // 待审核
	ApplicationApproved ApplicationStatus = "approved" // 已通过
	ApplicationRejected ApplicationStatus = "rejected" // 已拒绝
	ApplicationCanceled ApplicationStatus = "canceled" // 用户撤销
)

// Permission 版主权限标识（可动态扩展）

// ModeratorApplication 用户申请成为版主的记录。
// 唯一性约束：同一用户在同一板块同时只能存在一条 pending 状态的申请（由 service 层保证）。
type ModeratorApplication struct {
	common.BaseModel // 包含 ID, CreatedAt, UpdatedAt, DeletedAt

	// 申请基本信息
	UserID  uint `gorm:"not null;index" json:"user_id"`  // 申请用户ID
	BoardID uint `gorm:"not null;index" json:"board_id"` // 目标版块ID

	// 申请内容
	Reason string `gorm:"size:500" json:"reason"` // 申请理由
	// 申请时希望获得的权限列表（JSON数组，可任意扩展）
	RequestedPermissions datatypes.JSONSlice[ModeratorPermission] `gorm:"type:json" json:"requested_permissions"`

	// 审批相关
	Status     ApplicationStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	ReviewerID *uint             `gorm:"index" json:"reviewer_id"`    // 审批人ID
	ReviewNote string            `gorm:"size:500" json:"review_note"` // 审批备注
	ReviewedAt *time.Time        `json:"reviewed_at"`                 // 审批时间（批准/拒绝时设置）

	// 关联实体（仅用于查询时预加载）
	User     User  `gorm:"foreignKey:UserID"    json:"user,omitempty"`
	Board    Board `gorm:"foreignKey:BoardID"   json:"board,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// 申请状态有效性验证
func (as ApplicationStatus) IsValid() bool {
	switch as {
	case ApplicationPending, ApplicationApproved, ApplicationRejected, ApplicationCanceled:
		return true
	}
	return false
}

// ParseTargetType 严格解析，返回错误
func ParseApplicationStatus(s string) (ApplicationStatus, error) {
	as := ApplicationStatus(s)
	if as.IsValid() {
		return as, nil
	}
	return "", fmt.Errorf("invalid target status: %s", s)
}

// --

// IsPending 是否处于待审核状态
func (m *ModeratorApplication) IsPending() bool {
	return m.Status == ApplicationPending
}

// Approve 批准申请
func (m *ModeratorApplication) Approve(reviewerID uint, note string) error {
	if !m.IsPending() {
		return apperrors.ErrValidation
	}
	m.Status = ApplicationApproved
	m.ReviewerID = &reviewerID
	m.ReviewNote = note
	now := time.Now()
	m.ReviewedAt = &now
	return nil
}

// Reject 拒绝申请
func (m *ModeratorApplication) Reject(reviewerID uint, note string) error {
	if !m.IsPending() {
		return apperrors.ErrValidation
	}
	m.Status = ApplicationRejected
	m.ReviewerID = &reviewerID
	m.ReviewNote = note
	now := time.Now()
	m.ReviewedAt = &now
	return nil
}

// Cancel 用户自行撤销（仅 pending 状态可撤销）
func (m *ModeratorApplication) Cancel() error {
	if !m.IsPending() {
		return apperrors.ErrValidation
	}
	m.Status = ApplicationCanceled
	return nil
}

// ── 权限辅助方法（提升易用性）───────────────────────────────────────────────

// HasPermission 检查申请单是否请求了某个权限
func (m *ModeratorApplication) HasPermission(perm ModeratorPermission) bool {
	for _, p := range m.RequestedPermissions {
		if p == perm {
			return true
		}
	}
	return false
}

// GetRequestedPermissions 返回请求的权限列表（安全拷贝）
func (m *ModeratorApplication) GetRequestedPermissions() []ModeratorPermission {
	if len(m.RequestedPermissions) == 0 {
		return []ModeratorPermission{}
	}
	out := make([]ModeratorPermission, len(m.RequestedPermissions))
	copy(out, m.RequestedPermissions)
	return out
}
