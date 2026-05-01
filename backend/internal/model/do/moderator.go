package do

import "encoding/json"

// ── ModeratorPermissions 版主细粒度权限 ──────────────────────────────────────

// ModeratorPermissions 存储在 Moderator.Permissions（JSON 列）中。
// 每个字段对应一项操作能力，默认全为 false（最小权限原则）。
type ModeratorPermissions struct {
	CanDeletePost      bool `json:"can_delete_post"`
	CanPinPost         bool `json:"can_pin_post"`
	CanEditAnyPost     bool `json:"can_edit_any_post"`
	CanManageModerator bool `json:"can_manage_moderator"`
	CanBanUser         bool `json:"can_ban_user"`
}

// ── Moderator ────────────────────────────────────────────────────────────────

// Moderator 板块版主记录。
// Permissions 字段以 JSON 格式持久化在数据库中，
// 通过 GetPermissions / SetPermissions 读写，上层代码不直接操作 json.RawMessage。
type Moderator struct {
	BaseModel
	UserID      uint            `gorm:"not null;uniqueIndex:idx_user_board" json:"user_id"`
	BoardID     uint            `gorm:"not null;uniqueIndex:idx_user_board" json:"board_id"`
	Permissions json.RawMessage `gorm:"type:json"        json:"permissions" swaggertype:"object" swaggerio:"ignore"`
	User        User            `gorm:"foreignKey:UserID"  json:"user,omitempty"`
	Board       Board           `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}

// GetPermissions 将 JSON 列反序列化为 ModeratorPermissions。
// 若字段为空（旧数据 / NULL）则返回零值权限（全 false），不报错。
func (m *Moderator) GetPermissions() (ModeratorPermissions, error) {
	var p ModeratorPermissions
	if len(m.Permissions) == 0 {
		return p, nil
	}
	err := json.Unmarshal(m.Permissions, &p)
	return p, err
}

// SetPermissions 将 ModeratorPermissions 序列化并写入 JSON 列。
func (m *Moderator) SetPermissions(p ModeratorPermissions) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	m.Permissions = b
	return nil
}

// HasPermission 快速检查单项权限，避免调用方重复 switch。
func (m *Moderator) HasPermission(permission string) bool {
	p, err := m.GetPermissions()
	if err != nil {
		return false
	}
	switch permission {
	case "delete_post":
		return p.CanDeletePost
	case "pin_post":
		return p.CanPinPost
	case "edit_any_post":
		return p.CanEditAnyPost
	case "manage_moderator":
		return p.CanManageModerator
	case "ban_user":
		return p.CanBanUser
	}
	return false
}

// ── ModeratorApplication 版主申请表 ─────────────────────────────────────────

// ApplicationStatus 申请状态
type ApplicationStatus string

const (
	ApplicationPending  ApplicationStatus = "pending"  // 待审核
	ApplicationApproved ApplicationStatus = "approved" // 已通过
	ApplicationRejected ApplicationStatus = "rejected" // 已拒绝
	ApplicationCanceled ApplicationStatus = "canceled" // 用户撤销
)

// ModeratorApplication 用户申请成为版主的记录。
// 一个用户在同一板块同一时间只允许存在一条 pending 申请（unique 约束在 service 层保证）。
type ModeratorApplication struct {
	BaseModel
	UserID     uint              `gorm:"not null;index"               json:"user_id"`
	BoardID    uint              `gorm:"not null;index"               json:"board_id"`
	Reason     string            `gorm:"size:500"                     json:"reason"`
	Status     ApplicationStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ReviewerID *uint             `gorm:"index"                        json:"reviewer_id"` // 审批人（管理员）
	ReviewNote string            `gorm:"size:500"                     json:"review_note"` // 审批备注

	// 申请时希望获得的初始权限（审批人可在通过时调整）
	ReqDeletePost      bool `json:"req_delete_post"`
	ReqPinPost         bool `json:"req_pin_post"`
	ReqEditAnyPost     bool `json:"req_edit_any_post"`
	ReqManageModerator bool `json:"req_manage_moderator"`
	ReqBanUser         bool `json:"req_ban_user"`

	User     User  `gorm:"foreignKey:UserID"    json:"user,omitempty"`
	Board    Board `gorm:"foreignKey:BoardID"   json:"board,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// ApplyModeratorInput 用户申请版主的参数。

type ApplyModeratorInput struct {
	UserID             uint   `json:"-"` // 从 JWT 上下文注入，不接受客户端传参
	Username           string `json:"-"`
	BoardID            uint   `json:"-"` // 从 URL 路径注入
	Reason             string `json:"reason"               binding:"required,max=500"`
	ReqDeletePost      bool   `json:"req_delete_post"`
	ReqPinPost         bool   `json:"req_pin_post"`
	ReqEditAnyPost     bool   `json:"req_edit_any_post"`
	ReqManageModerator bool   `json:"req_manage_moderator"`
	ReqBanUser         bool   `json:"req_ban_user"`
}

// MARK: HTTP
// AddModeratorRequest swagger body
type AddModeratorRequest struct {
	UserID             uint `json:"user_id"              example:"1"    binding:"required"`
	CanDeletePost      bool `json:"can_delete_post"      example:"true"`
	CanPinPost         bool `json:"can_pin_post"         example:"true"`
	CanEditAnyPost     bool `json:"can_edit_any_post"    example:"false"`
	CanManageModerator bool `json:"can_manage_moderator" example:"false"`
	CanBanUser         bool `json:"can_ban_user"         example:"true"`
}

// BanUserRequest swagger body
type BanUserRequest struct {
	UserID    uint   `json:"user_id"    example:"1"                    binding:"required"`
	Reason    string `json:"reason"     example:"发布违规内容"              binding:"required"`
	ExpiresAt string `json:"expires_at" example:"2024-12-31T23:59:59Z"`
}
