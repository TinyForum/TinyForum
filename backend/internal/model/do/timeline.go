package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type ActionType string

const (
	ActionCreatePost    ActionType = "create_post"
	ActionCreateComment ActionType = "create_comment"
	ActionLikePost      ActionType = "like_post"
	ActionLikeComment   ActionType = "like_comment"
	ActionFollowUser    ActionType = "follow_user"
	ActionAcceptAnswer  ActionType = "accept_answer"
	ActionSignIn        ActionType = "sign_in"
)

type TimelineEvent struct {
	common.BaseModel
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	ActorID    uint       `gorm:"not null;index" json:"actor_id"`
	Action     ActionType `gorm:"type:varchar(50);index" json:"action"`
	TargetID   uint       `json:"target_id"`
	TargetType string     `gorm:"size:50" json:"target_type"`
	Payload    string     `gorm:"type:json" json:"payload"`
	Score      int        `gorm:"default:0" json:"score"`

	User  User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Actor User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

type UserTimeline struct {
	common.BaseModel
	UserID       uint      `gorm:"not null;uniqueIndex:idx_user_timeline" json:"user_id"`
	TimelineType string    `gorm:"type:varchar(20);default:'home'" json:"timeline_type"`
	LastReadAt   time.Time `json:"last_read_at"`
}

type TimelineSubscription struct {
	common.BaseModel
	SubscriberID uint   `gorm:"not null;index" json:"subscriber_id"`
	TargetUserID uint   `gorm:"not null;index" json:"target_user_id"`
	TargetType   string `gorm:"size:20;default:'user'" json:"target_type"`
	TargetID     uint   `json:"target_id"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`
}
