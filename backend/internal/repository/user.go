package repository

import (
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
