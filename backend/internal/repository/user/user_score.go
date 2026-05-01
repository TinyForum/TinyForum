package user

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *userRepository) GetScoreById(userID uint) (int, error) {
	var user do.User
	err := r.db.Model(&do.User{}).Select("score").Where("id = ?", userID).First(&user).Error
	return user.Score, err
}

func (r *userRepository) GetUsersScoreTotal() (int, error) {
	var totalScore int64
	err := r.db.Model(&do.User{}).Select("COALESCE(SUM(score), 0)").Scan(&totalScore).Error
	return int(totalScore), err
}

func (r *userRepository) GetEveryoneUsersScore() ([]do.User, error) {
	var users []do.User
	err := r.db.Model(&do.User{}).Select("id, score").Find(&users).Error
	return users, err
}

func (r *userRepository) AddScore(userID uint, score int) error {
	return r.db.Model(&do.User{}).Where("id = ?", userID).
		UpdateColumn("score", gorm.Expr("score + ?", score)).Error
}

func (r *userRepository) DeductScore(tx *gorm.DB, userID uint, score int) error {
	if score <= 0 {
		return nil
	}
	var user do.User
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}
	if user.Score < score {
		return errors.New("积分不足")
	}
	return tx.Model(&user).Update("score", gorm.Expr("score - ?", score)).Error
}

func (r *userRepository) SetScoreById(id uint, score int) error {
	if id == 0 {
		return errors.New("无效的用户ID")
	}
	if score < 0 {
		return errors.New("积分不能为负数")
	}
	if score > 999999 {
		return errors.New("积分超出最大限制")
	}
	result := r.db.Model(&do.User{}).
		Where("id = ?", id).
		Update("score", score)
	if result.Error != nil {
		return fmt.Errorf("更新积分失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
