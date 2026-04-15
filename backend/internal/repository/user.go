package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// MARK: 用户管理
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

// 列出用户
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

// 积分排名
func (r *UserRepository) GetTopUsers(limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Order("score DESC").Limit(limit).Find(&users).Error
	return users, err
}

// 关注排名
func (r *UserRepository) GetTopFollowers(limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Table("users").
		Select("users.*, COUNT(follows.follower_id) as follower_count").
		Joins("LEFT JOIN follows ON users.id = follows.following_id").
		Group("users.id").Order("follower_count DESC").Limit(limit).Find(&users).Error
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

// 用户查询自己的积分
func (r *UserRepository) GetScoreById(userID uint) (int, error) {
	var user model.User
	err := r.db.Model(&model.User{}).Select("score").Where("id = ?", userID).First(&user).Error
	return user.Score, err
}

// MARK: 互动相关操作
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

// MARK: 积分相关操作
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

func (r *UserRepository) AddScore(userID uint, score int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("score", gorm.Expr("score + ?", score)).Error
}

// 设置用户积分
func (r *UserRepository) SetScoreById(id uint, score int) error {
	// 1. 验证ID有效性
	if id == 0 {
		return errors.New("无效的用户ID")
	}

	// 2. 验证积分范围（根据业务需求调整）
	if score < 0 {
		return errors.New("积分不能为负数")
	}
	if score > 999999 {
		return errors.New("积分超出最大限制")
	}

	// 3. 执行更新操作
	result := r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("score", score)

	// 4. 检查更新结果
	if result.Error != nil {
		return fmt.Errorf("更新积分失败: %w", result.Error)
	}

	// 5. 检查是否更新了记录
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}

	return nil
}

// FindByIDs 批量查询用户
func (r *UserRepository) FindByIDs(ids []uint) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	var users []model.User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

// Count 返回用户总数
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

// CountByDateRange 统计指定时间段内新增用户数
func (r *UserRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// CountActiveByDateRange 统计指定时间段内活跃用户数（有发帖或评论行为）
func (r *UserRepository) CountActiveByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64

	// 活跃用户 = 该时间段内有发帖或发评论的去重用户数
	err := r.db.WithContext(ctx).
		Table("users u").
		Where(`u.deleted_at IS NULL AND EXISTS (
			SELECT 1 FROM posts p
			WHERE p.author_id = u.id AND p.deleted_at IS NULL
			  AND p.created_at BETWEEN ? AND ?
		) OR EXISTS (
			SELECT 1 FROM comments c
			WHERE c.author_id = u.id AND c.deleted_at IS NULL
			  AND c.created_at BETWEEN ? AND ?
		)`, startDate, endDate, startDate, endDate).
		Count(&count).Error

	return count, err
}

// ActiveUserRow 活跃用户查询结果行
type ActiveUserRow struct {
	ID       uint
	Username string
	Avatar   string
}

// GetActiveUsersByDateRange 获取指定时间段内活跃用户列表
func (r *UserRepository) GetActiveUsersByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*ActiveUserRow, error) {
	var rows []*ActiveUserRow

	err := r.db.WithContext(ctx).
		Table("users u").
		Select("u.id, u.username, u.avatar").
		Where(`u.deleted_at IS NULL AND (
			EXISTS (
				SELECT 1 FROM posts p
				WHERE p.author_id = u.id AND p.deleted_at IS NULL
				  AND p.created_at BETWEEN ? AND ?
			) OR EXISTS (
				SELECT 1 FROM comments c
				WHERE c.author_id = u.id AND c.deleted_at IS NULL
				  AND c.created_at BETWEEN ? AND ?
			)
		)`, startDate, endDate, startDate, endDate).
		Order("u.score DESC").
		Limit(limit).
		Scan(&rows).Error

	return rows, err
}
