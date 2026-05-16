// package do

// import (
// 	"fmt"
// 	"tiny-forum/internal/model/common"
// )

// // VoteType 投票类型（字符串枚举）
// type VoteType string

// const (
// 	VoteUp   VoteType = "up"   // 赞同
// 	VoteDown VoteType = "down" // 反对
// )

// // IsValid 验证投票类型是否合法
// func (v VoteType) IsValid() bool {
// 	switch v {
// 	case VoteUp, VoteDown:
// 		return true
// 	}
// 	return false
// }

// // ParseVoteType 严格解析字符串，用于 Service 层校验
// func ParseVoteType(s string) (VoteType, error) {
// 	vt := VoteType(s)
// 	if vt.IsValid() {
// 		return vt, nil
// 	}
// 	return "", fmt.Errorf("invalid vote type: %s", s)
// }

// // Vote 投票记录（针对评论的点赞/点踩）
// type Vote struct {
// 	common.BaseModel          // 内含 ID, CreatedAt, UpdatedAt, DeletedAt
// 	UserID           uint     `gorm:"not null;uniqueIndex:idx_user_comment;comment:用户ID" json:"user_id"`
// 	CommentID        uint     `gorm:"not null;uniqueIndex:idx_user_comment;comment:评论ID" json:"comment_id"`
// 	Type             VoteType `gorm:"type:varchar(10);not null;comment:投票类型(up/down)" json:"type"`
// 	// 如果未来需要扩展投票目标（如帖子），可添加 TargetType 字段，但当前简洁为主
// }

// // TableName 指定表名（可选，GORM 默认会使用 votes）
// func (Vote) TableName() string {
// 	return "votes"
// }

package do
