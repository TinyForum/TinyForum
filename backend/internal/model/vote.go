// internal/model/vote.go
package model

import "time"

type Vote struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `gorm:"index;uniqueIndex:idx_user_comment;comment:用户ID"`
	CommentID uint `gorm:"index;uniqueIndex:idx_user_comment;comment:评论ID"`
	Value     int  `gorm:"type:tinyint;comment:投票值 1:赞同 -1:反对"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Vote) TableName() string {
	return "votes"
}
