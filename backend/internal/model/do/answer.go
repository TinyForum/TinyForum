package do

import "tiny-forum/internal/model/common"

// type Answer struct {
// 	common.BaseModel
// 	QuestionID uint   `gorm:"not null;index" json:"question_id"`             // 问题ID
// 	ReplyID    uint   `gorm:"not null;index" json:"reply_id"`                // 回复ID
// 	Reply      *Reply `gorm:"foreignKey:ID;references:ReplyID" json:"reply"` // 回复
// 	IsAccepted bool   `gorm:"default:false;index" json:"is_accepted"`        // 是否被采纳
// 	VoteCount  int    `gorm:"default:0" json:"vote_count"`                   // 投票数
// }

// ==================== 回答（问题下的回答，不可嵌套） ====================

// Answer 回答扩展表（依附于 Reply，TargetType = "question"）
type Answer struct {
	common.BaseModel
	QuestionID uint `gorm:"not null;index" json:"question_id"`      // 所属问题ID
	ReplyID    uint `gorm:"not null;uniqueIndex" json:"reply_id"`   // 关联的回复ID（一对一）
	IsAccepted bool `gorm:"default:false;index" json:"is_accepted"` // 是否被采纳
	// VoteCount  int  `gorm:"default:0" json:"vote_count"`            // 投票数

	// 关联回复内容
	Reply *Reply `gorm:"foreignKey:ReplyID;references:ID" json:"reply,omitempty"`

	HelpfulCount int    `gorm:"default:0" json:"helpful_count"`         // 有帮助次数（区别于点赞）
	EditorNote   string `gorm:"type:text" json:"editor_note,omitempty"` // 编辑备注（管理员可见）
	UpVotes      int    `gorm:"column:up_votes" json:"up_votes"`
	DownVotes    int    `gorm:"column:down_votes" json:"down_votes"`
}
