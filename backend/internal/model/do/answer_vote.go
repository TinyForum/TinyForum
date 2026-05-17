package do

import (
	"fmt"
	"tiny-forum/internal/model/common"
)

type AnswerVote struct {
	common.BaseModel
	UserID    uint            `gorm:"uniqueIndex:idx_user_answer;not null" json:"user_id"`
	CommentID uint            `gorm:"uniqueIndex:idx_user_answer;not null" json:"comment_id"`
	VoteType  *AnswerVoteType `gorm:"type:varchar(10)" json:"vote_type"` // up/down
}

// AnnouncementType 公告类型（可存储）
type AnswerVoteType string

const (
	AnswerVoteTypeUp   AnswerVoteType = "up"   // 支持
	AnswerVoteTypeDown AnswerVoteType = "down" // 反对
)

var validVoteTypes = map[AnswerVoteType]bool{
	AnswerVoteTypeUp:   true,
	AnswerVoteTypeDown: true,
}

func (vt AnswerVoteType) IsValid() bool {
	return validVoteTypes[vt]
}

// 从字符串安全转换
func ParseAnswerVoteType(s string) (AnswerVoteType, error) {
	vt := AnswerVoteType(s)
	if vt.IsValid() {
		return vt, nil
	}
	return "", fmt.Errorf("invalid vote type: %s", s)
}
