package do

import "tiny-forum/internal/model/common"

type Question struct {
	common.BaseModel
	CreationID       uint     `gorm:"uniqueIndex;not null" json:"creation_id"`
	Creation         Creation `gorm:"foreignKey:CreationID;references:ID" json:"creation,omitempty"`
	AcceptedAnswerID *uint    `json:"accepted_answer_id"`
	RewardScore      int      `gorm:"default:0" json:"reward_score"`
	AnswerCount      int      `gorm:"default:0" json:"answer_count"`
	ViewCount        int      `gorm:"default:0" json:"view_count"`
	LikeCount        int      `gorm:"default:0" json:"like_count"`
	AcceptedAnswer   Comment  `gorm:"foreignKey:AcceptedAnswerID" json:"accepted_answer,omitempty"`
}
