package vote

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// GetUserVote 获取用户对评论的投票
func (r *voteRepository) GetUserVote(commentID, userID uint) (int, error) {
	var vote do.Vote
	err := r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).
		First(&vote).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return vote.Value, nil
}

// GetVoteCount 获取评论的投票数（从 vote_count 字段获取更高效）
func (r *voteRepository) GetVoteCount(commentID uint) (int, error) {
	var comment do.Comment
	err := r.db.Select("vote_count").First(&comment, commentID).Error
	if err != nil {
		return 0, err
	}
	return comment.VoteCount, nil
}

// GetVoteUsers 获取投票的用户列表
func (r *voteRepository) GetVoteUsers(commentID uint, voteType int) ([]do.User, error) {
	var users []do.User
	err := r.db.Table("users").
		Joins("INNER JOIN votes ON votes.user_id = users.id").
		Where("votes.comment_id = ? AND votes.value = ?", commentID, voteType).
		Find(&users).Error
	return users, err
}
