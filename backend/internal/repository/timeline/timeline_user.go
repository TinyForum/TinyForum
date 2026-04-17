package timeline

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// UpdateLastRead 更新用户某个时间线的最后阅读时间
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

// GetLastRead 获取用户某个时间线的最后阅读时间
func (r *TimelineRepository) GetLastRead(userID uint, timelineType string) (*model.UserTimeline, error) {
	var userTimeline model.UserTimeline
	err := r.db.Where("user_id = ? AND timeline_type = ?", userID, timelineType).
		First(&userTimeline).Error
	return &userTimeline, err
}
