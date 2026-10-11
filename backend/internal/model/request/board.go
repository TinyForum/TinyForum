package request

import "time"

type ReviewApplicationRequest struct {
	ApplicationID      uint   `json:"application_id" binding:"required"`
	Approve            bool   `json:"approve"`
	ReviewNote         string `json:"review_note" binding:"max=500"`
	CanDeletePost      *bool  `json:"can_delete_post"`
	CanPinPost         *bool  `json:"can_pin_post"`
	CanEditAnyPost     *bool  `json:"can_edit_any_post"`
	CanManageModerator *bool  `json:"can_manage_moderator"`
	CanBanUser         *bool  `json:"can_ban_user"`
}

// ApplyBoardCreationRequest 用户申请创建新板块
type ApplyBoardCreationRequest struct {
	UserID      uint   `json:"-"` // 由 handler 从认证上下文注入，禁止客户端伪造
	Name        string `json:"name"        binding:"required,min=2,max=50"`
	Slug        string `json:"slug"        binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=500"`
	Reason      string `json:"reason"      binding:"required,max=500"`
}

// ReviewBoardCreationRequest 管理员审批板块创建申请
type ReviewBoardCreationRequest struct {
	ApplicationID uint   `json:"application_id" binding:"required"`
	Approve       bool   `json:"approve"`
	ReviewNote    string `json:"review_note" binding:"max=500"`
}

type BoardBanUserRequest struct {
	UserID    uint       `json:"user_id"  binding:"required"`
	BoardID   uint       `json:"board_id" binding:"required"`
	Reason    string     `json:"reason"   binding:"required,max=500"`
	ExpiresAt *time.Time `json:"expires_at"`
}
