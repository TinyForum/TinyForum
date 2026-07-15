package request

type CreateTopicReqeust struct {
	Title       string `json:"title" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	CoverUrl    string `json:"cover_url" binding:"max=500"`
	IsPublic    bool   `json:"is_public"`
}

type AddPostToTopicRequest struct {
	TopicID   uint `json:"topic_id" binding:"required"`
	PostID    uint `json:"post_id" binding:"required"`
	SortOrder int  `json:"sort_order"`
}
