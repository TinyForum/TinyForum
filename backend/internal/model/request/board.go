package request

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
