package bo

type CreateBoardInput struct {
	Name        string `json:"name"        binding:"required,min=2,max=50"`
	Slug        string `json:"slug"        binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon"        binding:"max=100"`
	Cover       string `json:"cover"       binding:"max=500"`
	ParentID    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
	ViewRole    string `json:"view_role"`
	PostRole    string `json:"post_role"`
	ReplyRole   string `json:"reply_role"`
}
