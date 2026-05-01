package user

import (
	"tiny-forum/internal/model/po"
)

func (r *userRepository) Follow(followerID, followingID uint) error {
	follow := po.Follow{FollowerID: followerID, FollowingID: followingID}
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		FirstOrCreate(&follow).Error
}

func (r *userRepository) Unfollow(followerID, followingID uint) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&po.Follow{}).Error
}

func (r *userRepository) IsFollowing(followerID, followingID uint) bool {
	var count int64
	r.db.Model(&po.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count)
	return count > 0
}

func (r *userRepository) GetFollowerCount(userID uint) int64 {
	var count int64
	r.db.Model(&po.Follow{}).Where("following_id = ?", userID).Count(&count)
	return count
}

func (r *userRepository) GetFollowingCount(userID uint) int64 {
	var count int64
	r.db.Model(&po.Follow{}).Where("follower_id = ?", userID).Count(&count)
	return count
}

func (r *userRepository) GetFollowers(userID uint, page, pageSize int) ([]po.User, int64, error) {
	var users []po.User
	var total int64
	offset := (page - 1) * pageSize

	err := r.db.Model(&po.Follow{}).
		Where("following_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Model(&po.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) GetFollowing(userID uint, page, pageSize int) ([]po.User, int64, error) {
	var users []po.User
	var total int64
	offset := (page - 1) * pageSize

	err := r.db.Model(&po.Follow{}).
		Where("follower_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Model(&po.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}
