package timeline

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type TimelineRepository interface {
	// event
	CreateEvent(event *do.TimelineEvent) error
	GetUserTimeline(userID uint, limit, offset int) ([]do.TimelineEvent, int64, error)
	GetFollowingTimeline(userID uint, limit, offset int) ([]do.TimelineEvent, int64, error)
	GetEventByTarget(targetType string, targetID uint) ([]do.TimelineEvent, error)
	// subscription
	Subscribe(sub *do.TimelineSubscription) error
	Unsubscribe(subscriberID, targetUserID uint) error
	GetSubscriptions(subscriberID uint) ([]do.TimelineSubscription, error)
	IsSubscribed(subscriberID, targetUserID uint) (bool, error)
	// user
	UpdateLastRead(userID uint, timelineType string) error
	GetLastRead(userID uint, timelineType string) (*do.UserTimeline, error)
}

type timelineRepository struct {
	db *gorm.DB
}

func NewTimelineRepository(db *gorm.DB) TimelineRepository {
	return &timelineRepository{db: db}
}
