// internal/repository/vote_repository.go
package repository

import (
	"errors"

	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type VoteRepository struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

// CreateOrUpdateVote 创建或更新投票
func (r *VoteRepository) CreateOrUpdateVote(commentID, userID uint, value int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找现有投票
	var existingVote model.Vote
	err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).
		First(&existingVote).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}

	if existingVote.ID > 0 {
		// 更新现有投票
		oldValue := existingVote.Value
		diff := value - oldValue

		// 更新投票记录
		if err := tx.Model(&existingVote).Update("value", value).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 更新评论的投票计数
		if err := tx.Model(&model.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", diff)).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// 创建新投票
		vote := &model.Vote{
			UserID:    userID,
			CommentID: commentID,
			Value:     value,
		}
		if err := tx.Create(vote).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 更新评论的投票计数
		if err := tx.Model(&model.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", value)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// RemoveVote 删除投票
func (r *VoteRepository) RemoveVote(commentID, userID uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找投票记录
	var vote model.Vote
	if err := tx.Where("comment_id = ? AND user_id = ?", commentID, userID).
		First(&vote).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除投票记录
	if err := tx.Delete(&vote).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新评论的投票计数
	if err := tx.Model(&model.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("vote_count", gorm.Expr("vote_count - ?", vote.Value)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetUserVote 获取用户对评论的投票
func (r *VoteRepository) GetUserVote(commentID, userID uint) (int, error) {
	var vote model.Vote
	err := r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).
		First(&vote).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil // 未投票
	}
	if err != nil {
		return 0, err
	}

	return vote.Value, nil
}

// GetVoteCount 获取评论的投票数（从 vote_count 字段获取更高效）
func (r *VoteRepository) GetVoteCount(commentID uint) (int, error) {
	var comment model.Comment
	err := r.db.Select("vote_count").First(&comment, commentID).Error
	if err != nil {
		return 0, err
	}
	return comment.VoteCount, nil
}

// GetVoteUsers 获取投票的用户列表
func (r *VoteRepository) GetVoteUsers(commentID uint, voteType int) ([]model.User, error) {
	var users []model.User
	err := r.db.Table("users").
		Joins("INNER JOIN votes ON votes.user_id = users.id").
		Where("votes.comment_id = ? AND votes.value = ?", commentID, voteType).
		Find(&users).Error
	return users, err
}

// UnacceptAnswer 取消接受答案
func (r *CommentRepository) UnacceptAnswer(commentID uint) error {
	// 使用事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 更新评论的 is_accepted 字段
	if err := tx.Model(&model.Comment{}).
		Where("id = ?", commentID).
		Update("is_accepted", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 如果有单独的 Question 表存储 accepted_answer_id，也需要更新
	// 这里假设 Comment 表有 PostID，而 Post 表可能有 accepted_answer_id
	if err := tx.Model(&model.Post{}).
		Where("accepted_answer_id = ?", commentID).
		Update("accepted_answer_id", nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetAcceptedAnswer 获取问题已接受的答案
func (r *CommentRepository) GetAcceptedAnswer(postID uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Where("post_id = ? AND is_accepted = ?", postID, true).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
