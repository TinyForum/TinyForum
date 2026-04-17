package user

import (
	"tiny-forum/internal/model"
)

func (r *UserRepository) Follow(followerID, followingID uint) error {
	follow := model.Follow{FollowerID: followerID, FollowingID: followingID}
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		FirstOrCreate(&follow).Error
}

func (r *UserRepository) Unfollow(followerID, followingID uint) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&model.Follow{}).Error
}

func (r *UserRepository) IsFollowing(followerID, followingID uint) bool {
	var count int64
	r.db.Model(&model.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count)
	return count > 0
}

func (r *UserRepository) GetFollowerCount(userID uint) int64 {
	var count int64
	r.db.Model(&model.Follow{}).Where("following_id = ?", userID).Count(&count)
	return count
}

func (r *UserRepository) GetFollowingCount(userID uint) int64 {
	var count int64
	r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&count)
	return count
}

func (r *UserRepository) GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	offset := (page - 1) * pageSize

	err := r.db.Model(&model.Follow{}).
		Where("following_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Model(&model.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	offset := (page - 1) * pageSize

	err := r.db.Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Model(&model.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}
