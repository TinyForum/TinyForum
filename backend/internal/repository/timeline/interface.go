package timeline

import (
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

type TimelineRepository interface {
	// event
	CreateEvent(event *po.TimelineEvent) error
	GetUserTimeline(userID uint, limit, offset int) ([]po.TimelineEvent, int64, error)
	GetFollowingTimeline(userID uint, limit, offset int) ([]po.TimelineEvent, int64, error)
	GetEventByTarget(targetType string, targetID uint) ([]po.TimelineEvent, error)
	// subscription
	Subscribe(sub *po.TimelineSubscription) error
	Unsubscribe(subscriberID, targetUserID uint) error
	GetSubscriptions(subscriberID uint) ([]po.TimelineSubscription, error)
	IsSubscribed(subscriberID, targetUserID uint) (bool, error)
	// user
	UpdateLastRead(userID uint, timelineType string) error
	GetLastRead(userID uint, timelineType string) (*po.UserTimeline, error)
}

type timelineRepository struct {
	db *gorm.DB
}

func NewTimelineRepository(db *gorm.DB) TimelineRepository {
	return &timelineRepository{db: db}
}
