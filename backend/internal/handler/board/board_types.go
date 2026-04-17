package board

// ReviewApplicationRequest 审批版主申请请求
type ReviewApplicationRequest struct {
	Approve            bool   `json:"approve" binding:"required"`
	ReviewNote         string `json:"review_note" binding:"max=500"`
	CanDeletePost      bool   `json:"can_delete_post"`
	CanPinPost         bool   `json:"can_pin_post"`
	CanEditAnyPost     bool   `json:"can_edit_any_post"`
	CanManageModerator bool   `json:"can_manage_moderator"`
	CanBanUser         bool   `json:"can_ban_user"`
}
