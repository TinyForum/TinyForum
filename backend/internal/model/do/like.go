package do

import (
	"fmt"
	"tiny-forum/internal/model/common"
)

type Like struct {
	common.BaseModel
	UserID     uint       `gorm:"not null;uniqueIndex:idx_like,priority:1;index" json:"user_id"`                // 用户ID
	TargetType TargetType `gorm:"type:varchar(20);not null;uniqueIndex:idx_like,priority:2" json:"target_type"` // 目标类型
	TargetID   uint       `gorm:"not null;uniqueIndex:idx_like,priority:3;index" json:"target_id"`              // 目标ID

	// 关联（不序列化，避免递归）
	User *User `gorm:"foreignKey:UserID" json:"-"` // 用户
}

func (Like) TableName() string { return "likes" }

type TargetType string

const (
	LikeTargetPost    TargetType = "post"    // 帖子
	LikeTargetComment TargetType = "comment" // 评论
)

// TargetType 有效性验证
func (tt TargetType) IsValid() bool {
	switch tt {
	case LikeTargetPost, LikeTargetComment:
		return true
	}
	return false
}

// ParseTargetType 严格解析，返回错误
func ParseTargetType(s string) (TargetType, error) {
	tt := TargetType(s)
	if tt.IsValid() {
		return tt, nil
	}
	return "", fmt.Errorf("invalid target type: %s", s)
}
