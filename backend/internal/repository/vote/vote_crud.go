package vote

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// CreateOrUpdateVote 创建或更新投票
func (r *voteRepository) CreateOrUpdateVote(commentID, userID uint, value int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingVote do.Vote
	err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).
		First(&existingVote).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	if existingVote.ID > 0 {
		oldValue := existingVote.Value
		diff := value - oldValue

		if err := tx.Model(&existingVote).Update("value", value).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&do.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		vote := &do.Vote{
			UserID:    userID,
			CommentID: commentID,
			Value:     value,
		}
		if err := tx.Create(vote).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&do.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", value)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemoveVote 删除投票
func (r *voteRepository) RemoveVote(commentID, userID uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var vote do.Vote
	if err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).
		First(&vote).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&vote).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&do.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("vote_count", gorm.Expr("vote_count - ?", vote.Value)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
