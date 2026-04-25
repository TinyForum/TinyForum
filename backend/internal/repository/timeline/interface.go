package timeline

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type TimelineRepository interface {
	// event
	CreateEvent(event *model.TimelineEvent) error
	GetUserTimeline(userID uint, limit, offset int) ([]model.TimelineEvent, int64, error)
	GetFollowingTimeline(userID uint, limit, offset int) ([]model.TimelineEvent, int64, error)
	GetEventByTarget(targetType string, targetID uint) ([]model.TimelineEvent, error)
	// subscription
	Subscribe(sub *model.TimelineSubscription) error
	Unsubscribe(subscriberID, targetUserID uint) error
	GetSubscriptions(subscriberID uint) ([]model.TimelineSubscription, error)
	IsSubscribed(subscriberID, targetUserID uint) (bool, error)
	// user
	UpdateLastRead(userID uint, timelineType string) error
	GetLastRead(userID uint, timelineType string) (*model.UserTimeline, error)
}

type timelineRepository struct {
	db *gorm.DB
}

func NewTimelineRepository(db *gorm.DB) TimelineRepository {
	return &timelineRepository{db: db}
}
