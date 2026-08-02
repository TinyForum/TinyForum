package topic

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *topicRepository) Follow(follow *do.TopicFollow) error {
	// 1. 查询所有记录（包括软删除）
	var existing do.TopicFollow
	err := r.db.Unscoped().Where("user_id = ? AND topic_id = ?", follow.UserID, follow.TopicID).
		First(&existing).Error

	if err == nil {
		// 记录存在（无论是有效还是软删除）
		if existing.DeletedAt.Valid {
			// 如果已软删除，恢复（设置 deleted_at = NULL）
			return r.db.Unscoped().Model(&existing).Update("deleted_at", nil).Error
		}
		// 否则已经有效关注，幂等返回 nil（follow.ID 保持为 0 表示未创建新记录）
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // 其他数据库错误
	}

	// 2. 完全不存在，创建新记录
	return r.db.Create(follow).Error
}
func (r *topicRepository) Unfollow(userID, topicID uint) error {
	result := r.db.Where("user_id = ? AND topic_id = ?", userID, topicID).
		Delete(&do.TopicFollow{})
	if result.Error != nil {
		return result.Error
	}
	// 幂等：取消不存在的关注不报错
	return nil
}

func (r *topicRepository) IsFollowing(userID, topicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&do.TopicFollow{}).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Count(&count).Error
	return count > 0, err
}

func (r *topicRepository) GetFollowers(topicID uint, limit, offset int) ([]do.TopicFollow, int64, error) {
	var follows []do.TopicFollow
	var total int64

	query := r.db.Model(&do.TopicFollow{}).Where("topic_id = ?", topicID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("User").
		Order("created_at DESC").
		Find(&follows).Error
	return follows, total, err
}
