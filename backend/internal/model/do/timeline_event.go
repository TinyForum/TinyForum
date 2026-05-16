package do

import (
	"fmt"
	"tiny-forum/internal/model/common"

	"gorm.io/datatypes"
)

// TimelineEvent 时间线事件（用户动态）
type TimelineEvent struct {
	common.BaseModel
	UserID     uint               `gorm:"not null;index:idx_user_time,priority:1" json:"user_id"` // 事件所有者（被影响用户）
	ActorID    uint               `gorm:"not null;index" json:"actor_id"`                         // 操作者（执行动作的用户）
	Action     ActionType         `gorm:"type:varchar(50);index" json:"action"`                   // 动作类型
	TargetID   uint               `json:"target_id"`                                              // 目标对象ID（如帖子ID、评论ID）
	TargetType TilelineTargetType `gorm:"size:50" json:"target_type"`                             // 目标类型（post, comment, user）
	Payload    datatypes.JSON     `gorm:"type:json" json:"payload"`                               // 附加数据（如标题摘要等）
	Score      int                `gorm:"default:0" json:"score"`                                 // 动态权重（用于热度排序）

	User  User `gorm:"foreignKey:UserID" json:"user,omitempty"`   // 事件所有者
	Actor User `gorm:"foreignKey:ActorID" json:"actor,omitempty"` // 操作者
}

// 目标类型
type TilelineTargetType string

const (
	TargetTypePost    TilelineTargetType = "post"    // 帖子
	TargetTypeComment TilelineTargetType = "comment" // 评论
	TargetTypeUser    TilelineTargetType = "user"    // 用户
)

// enum [post, comment, user]

// IsValid 验证动作类型是否合法
func (a TilelineTargetType) IsValid() bool {
	switch a {
	case TargetTypePost, TargetTypeComment, TargetTypeUser:
		return true
	}
	return false
}

// ParseActionType 严格解析，返回错误
func ParseTilelineTargetType(s string) (TilelineTargetType, error) {
	a := TilelineTargetType(s)
	if a.IsValid() {
		return a, nil
	}
	return "", fmt.Errorf("invalid target type: %s", s)
}

// ========================
// 时间线事件
// ========================

// ActionType 动作类型
type ActionType string

const (
	ActionCreatePost    ActionType = "create_post"    // 创建帖子
	ActionCreateComment ActionType = "create_comment" // 创建评论
	ActionLikePost      ActionType = "like_post"      // 点赞帖子
	ActionLikeComment   ActionType = "like_comment"   // 点赞评论
	ActionFollowUser    ActionType = "follow_user"    // 关注用户
	ActionAcceptAnswer  ActionType = "accept_answer"  // 接受回答
	ActionSignIn        ActionType = "sign_in"        // 签到
)

// enum [create_post, create_comment, like_post, like_comment, follow_user, accept_answer, sign_in]
// IsValid 验证动作类型是否合法
func (a ActionType) IsValid() bool {
	switch a {
	case ActionCreatePost, ActionCreateComment, ActionLikePost,
		ActionLikeComment, ActionFollowUser, ActionAcceptAnswer, ActionSignIn:
		return true
	}
	return false
}

// ParseActionType 严格解析，返回错误
func ParseActionType(s string) (ActionType, error) {
	a := ActionType(s)
	if a.IsValid() {
		return a, nil
	}
	return "", fmt.Errorf("invalid action type: %s", s)
}
