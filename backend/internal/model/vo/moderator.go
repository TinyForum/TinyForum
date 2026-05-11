package vo

import (
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
)

type ModeratorApplication struct {
	common.BaseModel
	UserID     uint                 `             json:"user_id"` //
	BoardID    uint                 `gorm:"not null;index"               json:"board_id"`
	Reason     string               `gorm:"size:500"                     json:"reason"`
	Status     do.ApplicationStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ReviewerID *uint                `gorm:"index"                        json:"reviewer_id"` // 审批人（管理员）
	ReviewNote string               `gorm:"size:500"                     json:"review_note"` // 审批备注

	// 申请时希望获得的初始权限（审批人可在通过时调整）
	ReqDeletePost      bool `json:"req_delete_post"`
	ReqPinPost         bool `json:"req_pin_post"`
	ReqEditAnyPost     bool `json:"req_edit_any_post"`
	ReqManageModerator bool `json:"req_manage_moderator"`
	ReqBanUser         bool `json:"req_ban_user"`

	User     do.User  `gorm:"foreignKey:UserID"    json:"user,omitempty"`
	Board    do.Board `gorm:"foreignKey:BoardID"   json:"board,omitempty"`
	Reviewer *do.User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// ModeratorApplicationVO 版主申请脱敏视图
type ModeratorApplicationVO struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"` // 时间格式按需
	UpdatedAt string `json:"updated_at"`

	UserID  uint `json:"user_id"`
	BoardID uint `json:"board_id"`

	// 脱敏后的申请人信息（仅展示必要公开字段）
	User struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		// 不包含 email, password, role 等
	} `json:"user"`

	// 脱敏后的版块信息
	Board struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"board"`

	Reason     string `json:"reason"`
	Status     string `json:"status"` // 申请状态
	ReviewNote string `json:"review_note"`

	// 审批人脱敏
	Reviewer struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	} `json:"reviewer"`

	// 申请权限（按需保留）
	ReqDeletePost      bool `json:"req_delete_post"`
	ReqPinPost         bool `json:"req_pin_post"`
	ReqEditAnyPost     bool `json:"req_edit_any_post"`
	ReqManageModerator bool `json:"req_manage_moderator"`
	ReqBanUser         bool `json:"req_ban_user"`
}
