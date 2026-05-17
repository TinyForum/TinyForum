package vote

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// GetUserVote 获取用户对评论的投票类型
// 返回值：voteType（可能为 nil 表示未投票），error
func (r *voteRepository) GetUserVote(commentID, userID uint) (*do.AnswerVoteType, error) {
	var vote do.AnswerVote
	err := r.db.Where("comment_id = ? AND user_id = ?", commentID, userID).
		Select("vote_type").
		First(&vote).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 无投票记录
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户投票失败: %w", err)
	}
	return vote.VoteType, nil
}

// GetVoteCount 获取评论的净投票数（up 数量 - down 数量）
// 直接从 Comment 表的 vote_count 字段读取，高效
func (r *voteRepository) GetVoteCount(commentID uint) (int, error) {
	var comment do.Comment
	err := r.db.Select("vote_count").Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // 评论不存在，返回 0
		}
		return 0, fmt.Errorf("查询评论投票数失败: %w", err)
	}
	return comment.VoteCount, nil
}

// GetVoteUsers 获取对某评论投了指定票型的用户列表
func (r *voteRepository) GetVoteUsers(commentID uint, voteType do.AnswerVoteType) ([]do.User, error) {
	var users []do.User
	err := r.db.Table("users").
		Joins("INNER JOIN answer_votes ON answer_votes.user_id = users.id").
		Where("answer_votes.comment_id = ? AND answer_votes.vote_type = ?", commentID, voteType).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("查询投票用户列表失败: %w", err)
	}
	return users, nil
}
