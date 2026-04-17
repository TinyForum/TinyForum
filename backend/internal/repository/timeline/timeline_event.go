package timeline

import (
	"tiny-forum/internal/model"
)

// CreateEvent 创建时间线事件
func (r *TimelineRepository) CreateEvent(event *model.TimelineEvent) error {
	return r.db.Create(event).Error
}

// GetUserTimeline 获取用户时间线（包含用户自己的事件和与自己相关的事件）
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

// GetFollowingTimeline 获取关注用户的时间线（仅关注用户的事件）
func (r *TimelineRepository) GetFollowingTimeline(userID uint, limit, offset int) ([]model.TimelineEvent, int64, error) {
	var events []model.TimelineEvent
	var total int64

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

// GetEventByTarget 根据目标类型和目标ID查询时间线事件
func (r *TimelineRepository) GetEventByTarget(targetType string, targetID uint) ([]model.TimelineEvent, error) {
	var events []model.TimelineEvent
	err := r.db.Where("target_type = ? AND target_id = ?", targetType, targetID).
		Find(&events).Error
	return events, err
}
