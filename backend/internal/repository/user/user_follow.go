package user

import (
	"time"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm/clause"
)

// 关注：如果有软删除记录则恢复，否则新建
func (r *userRepository) Follow(followerID, followingID uint) error {
	// 尝试恢复软删除记录（仅当存在 deleted_at IS NOT NULL 时）
	result := r.db.Unscoped().Model(&do.Follow{}).
		Where("follower_id = ? AND following_id = ? AND deleted_at IS NOT NULL", followerID, followingID).
		Updates(map[string]interface{}{
			"deleted_at": nil,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}

	// 使用 OnConflict 处理已存在活跃记录的情况（忽略插入）
	follow := do.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	insertResult := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&follow)
	return insertResult.Error
}

// 取消关注：软删除活跃记录
func (r *userRepository) Unfollow(followerID, followingID uint) error {
	return r.db.Model(&do.Follow{}).
		Where("follower_id = ? AND following_id = ? AND deleted_at IS NULL", followerID, followingID).
		Update("deleted_at", time.Now()).Error
}

// 删除关注：物理删除（不区分软删除或活跃，彻底删除记录）
func (r *userRepository) DeleteFollow(followerID, followingID uint) error {
	return r.db.Unscoped().Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&do.Follow{}).Error
}

func (r *userRepository) IsFollowing(followerID, followingID uint) bool {
	var count int64
	r.db.Model(&do.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count)
	return count > 0
}

func (r *userRepository) GetFollowerCount(userID uint) int64 {
	var count int64
	r.db.Model(&do.Follow{}).Where("following_id = ?", userID).Count(&count)
	return count
}

func (r *userRepository) GetFollowingCount(userID uint) int64 {
	var count int64
	r.db.Model(&do.Follow{}).Where("follower_id = ?", userID).Count(&count)
	return count
}

func (r *userRepository) GetFollowers(userID uint, page, pageSize int) ([]do.User, int64, error) {
	var users []do.User
	var total int64
	offset := (page - 1) * pageSize

	err := r.db.Model(&do.Follow{}).
		Where("following_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Model(&do.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) GetFollowing(userID uint, page, pageSize int) ([]do.User, int64, error) {
	var users []do.User
	var total int64
	offset := (page - 1) * pageSize

	err := r.db.Model(&do.Follow{}).
		Where("follower_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Model(&do.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}
