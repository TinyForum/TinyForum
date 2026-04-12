package model

type Question struct {
	BaseModel
	PostID           uint  `gorm:"uniqueIndex;not null" json:"post_id"`
	AcceptedAnswerID *uint `json:"accepted_answer_id"`
	RewardScore      int   `gorm:"default:0" json:"reward_score"`
	AnswerCount      int   `gorm:"default:0" json:"answer_count"`

	Post           Post    `gorm:"foreignKey:PostID" json:"post,omitempty"`
	AcceptedAnswer Comment `gorm:"foreignKey:AcceptedAnswerID" json:"accepted_answer,omitempty"`
}

type AnswerVote struct {
	BaseModel
	UserID    uint   `gorm:"uniqueIndex:idx_user_answer;not null" json:"user_id"`
	CommentID uint   `gorm:"uniqueIndex:idx_user_answer;not null" json:"comment_id"`
	VoteType  string `gorm:"type:varchar(10)" json:"vote_type"` // up/down
}
