package do

import (
	"tiny-forum/internal/model/common"
)

// BehaviorType 用户行为类型
// enum [view, like, unlike, comment, share, favorite, unfavorite, follow, unfollow, click, dwell, swipe, not_interested, report]
type BehaviorType string

const (
	BehaviorView          BehaviorType = "view"
	BehaviorLike          BehaviorType = "like"
	BehaviorUnlike        BehaviorType = "unlike"
	BehaviorComment       BehaviorType = "comment"
	BehaviorShare         BehaviorType = "share"
	BehaviorFavorite      BehaviorType = "favorite"
	BehaviorUnfavorite    BehaviorType = "unfavorite"
	BehaviorFollow        BehaviorType = "follow"
	BehaviorUnfollow      BehaviorType = "unfollow"
	BehaviorClick         BehaviorType = "click"
	BehaviorDwell         BehaviorType = "dwell"
	BehaviorSwipe         BehaviorType = "swipe"
	BehaviorNotInterested BehaviorType = "not_interested"
	BehaviorReport        BehaviorType = "report"
)

func (b BehaviorType) IsValid() bool {
	switch b {
	case BehaviorView, BehaviorLike, BehaviorUnlike, BehaviorComment, BehaviorShare,
		BehaviorFavorite, BehaviorUnfavorite, BehaviorFollow, BehaviorUnfollow,
		BehaviorClick, BehaviorDwell, BehaviorSwipe, BehaviorNotInterested, BehaviorReport:
		return true
	}
	return false
}

// UserBehaviorEvent 用户行为事件：记录每一次用户与内容的交互
type UserBehaviorEvent struct {
	common.BaseModel
	UserID      uint         `gorm:"not null;index" json:"user_id"`
	TargetID    uint         `gorm:"not null;index" json:"target_id"`
	TargetType  string       `gorm:"type:varchar(30);not null;default:'creation'" json:"target_type"`
	BehaviorType BehaviorType `gorm:"type:varchar(30);not null" json:"behavior_type"`
	Value       float64      `gorm:"not null;default:1.0" json:"value"`
	SessionID   string       `gorm:"type:varchar(64);not null;default:''" json:"session_id"`
	ContextJSON string       `gorm:"type:text;not null;default:'{}'" json:"context_json"`
	CreatedTS   int64        `gorm:"not null;default:0" json:"created_ts"`
}

func (UserBehaviorEvent) TableName() string {
	return "user_behavior_events"
}

// ContentFeature 内容特征：缓存内容的特征向量与元数据
type ContentFeature struct {
	common.BaseModel
	CreationID     uint    `gorm:"not null;uniqueIndex" json:"creation_id"`
	TagIDsJSON     string  `gorm:"type:text;not null;default:'[]'" json:"tag_ids_json"`
	BoardID        uint    `gorm:"not null;default:0" json:"board_id"`
	AuthorID       uint    `gorm:"not null;default:0" json:"author_id"`
	QualityScore   float64 `gorm:"not null;default:0" json:"quality_score"`
	HotScore       float64 `gorm:"not null;default:0" json:"hot_score"`
	FreshnessScore float64 `gorm:"not null;default:1.0" json:"freshness_score"`
	FeatureVector  string  `gorm:"type:text;not null;default:''" json:"feature_vector"`
}

func (ContentFeature) TableName() string {
	return "content_features"
}

// FeedbackType 推荐反馈类型
// enum [impression, click, dwell, dismiss, not_interested]
type FeedbackType string

const (
	FeedbackImpression    FeedbackType = "impression"
	FeedbackClick         FeedbackType = "click"
	FeedbackDwell         FeedbackType = "dwell"
	FeedbackDismiss       FeedbackType = "dismiss"
	FeedbackNotInterested FeedbackType = "not_interested"
)

func (f FeedbackType) IsValid() bool {
	switch f {
	case FeedbackImpression, FeedbackClick, FeedbackDwell, FeedbackDismiss, FeedbackNotInterested:
		return true
	}
	return false
}

// RecommendationFeedback 推荐反馈：记录用户看到推荐内容后的行为反馈
type RecommendationFeedback struct {
	common.BaseModel
	UserID       uint         `gorm:"not null;index" json:"user_id"`
	CreationID   uint         `gorm:"not null;index" json:"creation_id"`
	FeedbackType FeedbackType `gorm:"type:varchar(20);not null;default:'impression'" json:"feedback_type"`
	SourceType   string       `gorm:"type:varchar(30);not null;default:'recommend'" json:"source_type"`
	Position     int          `gorm:"not null;default:0" json:"position"`
	SessionID    string       `gorm:"type:varchar(64);not null;default:''" json:"session_id"`
}

func (RecommendationFeedback) TableName() string {
	return "recommendation_feedbacks"
}

// UserInterestProfile 用户兴趣画像：缓存用户近期兴趣向量
type UserInterestProfile struct {
	common.BaseModel
	UserID          uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	TagWeightsJSON  string    `gorm:"type:text;not null;default:'{}'" json:"tag_weights_json"`
	BoardWeightsJSON string   `gorm:"type:text;not null;default:'{}'" json:"board_weights_json"`
	ActiveTagsJSON  string    `gorm:"type:text;not null;default:'[]'" json:"active_tags_json"`
	LastUpdatedAt   string    `gorm:"type:timestamptz;not null;default:now()" json:"last_updated_at"`
}

func (UserInterestProfile) TableName() string {
	return "user_interest_profiles"
}
