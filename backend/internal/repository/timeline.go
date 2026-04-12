package repository

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type TimelineRepository struct {
	db *gorm.DB
}

func NewTimelineRepository(db *gorm.DB) *TimelineRepository {
	return &TimelineRepository{db: db}
}

// Event methods
func (r *TimelineRepository) CreateEvent(event *model.TimelineEvent) error {
	return r.db.Create(event).Error
}

func (r *TimelineRepository) GetUserTimeline(userID uint, limit, offset int) ([]model.TimelineEvent, int64, error) {
	var events []model.TimelineEvent
	var total int64

	query := r.db.Model(&model.TimelineEvent{}).
		Where("user_id = ? OR actor_id = ?", userID, userID)

	query.Count(&total)

	err := query.Limit(limit).Offset(offset).
		Preload("User").
		Preload("Actor").
		Order("score DESC, created_at DESC").
		Find(&events).Error

	return events, total, err
}

func (r *TimelineRepository) GetFollowingTimeline(userID uint, limit, offset int) ([]model.TimelineEvent, int64, error) {
	var events []model.TimelineEvent
	var total int64

	// Get followed user IDs
	subQuery := r.db.Table("timeline_subscriptions").
		Select("target_user_id").
		Where("subscriber_id = ? AND target_type = ? AND is_active = ?", userID, "user", true)

	query := r.db.Model(&model.TimelineEvent{}).
		Where("user_id IN (?)", subQuery)

	query.Count(&total)

	err := query.Limit(limit).Offset(offset).
		Preload("User").
		Preload("Actor").
		Order("score DESC, created_at DESC").
		Find(&events).Error

	return events, total, err
}

func (r *TimelineRepository) GetEventByTarget(targetType string, targetID uint) ([]model.TimelineEvent, error) {
	var events []model.TimelineEvent
	err := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Find(&events).Error
	return events, err
}

// Subscription methods
func (r *TimelineRepository) Subscribe(sub *model.TimelineSubscription) error {
	// Check if already exists
	var existing model.TimelineSubscription
	err := r.db.Where("subscriber_id = ? AND target_user_id = ?", sub.SubscriberID, sub.TargetUserID).
		First(&existing).Error

	if err == nil {
		// Already exists, reactivate
		return r.db.Model(&existing).Update("is_active", true).Error
	}

	return r.db.Create(sub).Error
}

func (r *TimelineRepository) Unsubscribe(subscriberID, targetUserID uint) error {
	return r.db.Model(&model.TimelineSubscription{}).
		Where("subscriber_id = ? AND target_user_id = ?", subscriberID, targetUserID).
		Update("is_active", false).Error
}

func (r *TimelineRepository) GetSubscriptions(subscriberID uint) ([]model.TimelineSubscription, error) {
	var subs []model.TimelineSubscription
	err := r.db.Where("subscriber_id = ? AND is_active = ?", subscriberID, true).
		Find(&subs).Error
	return subs, err
}

func (r *TimelineRepository) IsSubscribed(subscriberID, targetUserID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.TimelineSubscription{}).
		Where("subscriber_id = ? AND target_user_id = ? AND is_active = ?",
			subscriberID, targetUserID, true).
		Count(&count).Error
	return count > 0, err
}

// UserTimeline methods
func (r *TimelineRepository) UpdateLastRead(userID uint, timelineType string) error {
	var userTimeline model.UserTimeline

	err := r.db.Where("user_id = ? AND timeline_type = ?", userID, timelineType).
		First(&userTimeline).Error

	if err == gorm.ErrRecordNotFound {
		userTimeline = model.UserTimeline{
			UserID:       userID,
			TimelineType: timelineType,
		}
		return r.db.Create(&userTimeline).Error
	}

	return r.db.Model(&userTimeline).Update("last_read_at", gorm.Expr("NOW()")).Error
}

func (r *TimelineRepository) GetLastRead(userID uint, timelineType string) (*model.UserTimeline, error) {
	var userTimeline model.UserTimeline
	err := r.db.Where("user_id = ? AND timeline_type = ?", userID, timelineType).
		First(&userTimeline).Error
	return &userTimeline, err
}
