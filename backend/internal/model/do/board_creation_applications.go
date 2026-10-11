package do

import (
	"time"

	"tiny-forum/internal/model/common"
	apperrors "tiny-forum/pkg/errors"
)

// ── BoardCreationApplication 板块创建申请表 ─────────────────────────────────

// BoardCreationApplication 用户申请创建新板块的记录。
// 唯一性约束：同一板块标识（slug）同时只能存在一条 pending 状态的申请（由 service 层保证）。
type BoardCreationApplication struct {
	common.BaseModel // 包含 ID, CreatedAt, UpdatedAt, DeletedAt

	// 申请基本信息
	UserID uint `gorm:"not null;index" json:"user_id"` // 申请用户ID

	// 拟创建板块信息
	Name        string `gorm:"not null;size:50" json:"name"`       // 拟建板块名称
	Slug        string `gorm:"not null;size:50;index" json:"slug"` // 拟建板块标识
	Description string `gorm:"size:500" json:"description"`        // 拟建板块描述

	// 申请内容
	Reason string `gorm:"size:500" json:"reason"` // 申请理由

	// 审批相关
	Status         ApplicationStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"` // 审批状态
	ReviewerID     *uint             `gorm:"index" json:"reviewer_id"`                               // 审批人ID
	ReviewNote     string            `gorm:"size:500" json:"review_note"`                            // 审批备注
	ReviewedAt     *time.Time        `json:"reviewed_at"`                                            // 审批时间
	CreatedBoardID *uint             `gorm:"index" json:"created_board_id"`                          // 审批通过后创建的板块ID

	// 关联实体（仅用于查询时预加载）
	User         User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reviewer     *User  `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	CreatedBoard *Board `gorm:"foreignKey:CreatedBoardID" json:"created_board,omitempty"`
}

// IsPending 是否处于待审核状态
func (b *BoardCreationApplication) IsPending() bool {
	return b.Status == ApplicationPending
}

// Approve 批准申请
func (b *BoardCreationApplication) Approve(reviewerID uint, note string) error {
	if !b.IsPending() {
		return apperrors.ErrValidation
	}
	b.Status = ApplicationApproved
	b.ReviewerID = &reviewerID
	b.ReviewNote = note
	now := time.Now()
	b.ReviewedAt = &now
	return nil
}

// Reject 拒绝申请
func (b *BoardCreationApplication) Reject(reviewerID uint, note string) error {
	if !b.IsPending() {
		return apperrors.ErrValidation
	}
	b.Status = ApplicationRejected
	b.ReviewerID = &reviewerID
	b.ReviewNote = note
	now := time.Now()
	b.ReviewedAt = &now
	return nil
}
