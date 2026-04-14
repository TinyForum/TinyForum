package repository

import (
	"errors"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *UserRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.Model(&model.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// Follow/unfollow
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

// Score
func (r *UserRepository) AddScore(userID uint, score int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("score", gorm.Expr("score + ?", score)).Error
}

// Ranking
func (r *UserRepository) GetTopUsers(limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Order("score DESC").Limit(limit).Find(&users).Error
	return users, err
}

// GetFollowing 获取用户关注的列表
func (r *UserRepository) GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * pageSize

	// 获取关注总数
	err := r.db.Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取关注的用户列表
	err = r.db.Model(&model.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// GetFollowing 获取关注用户的列表
// GetFollowers 获取用户的粉丝列表（谁关注了该用户）
func (r *UserRepository) GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * pageSize

	// 获取粉丝总数
	err := r.db.Model(&model.Follow{}).
		Where("following_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取粉丝用户列表
	err = r.db.Model(&model.Follow{}).
		Select("users.*").
		Joins("JOIN users ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// DeductScore 扣减用户积分（使用事务）
func (r *UserRepository) DeductScore(tx *gorm.DB, userID uint, score int) error {
	if score <= 0 {
		return nil
	}

	var user model.User
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	if user.Score < score {
		return errors.New("积分不足")
	}

	return tx.Model(&user).Update("score", gorm.Expr("score - ?", score)).Error
}

// internal/repository/user_repository.go

// FindByIDs 批量查询用户
func (r *UserRepository) FindByIDs(ids []uint) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	var users []model.User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}
