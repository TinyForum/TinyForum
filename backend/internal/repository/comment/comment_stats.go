package comment

import (
	"context"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// CountByPost 统计帖子的评论总数
func (r *commentRepository) CountByPost(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// Count 统计所有评论总数
func (r *commentRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Comment{}).Count(&count).Error
	return count, err
}

// CountByDateRange 统计指定时间段内新增评论数
func (r *commentRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// IncrLikeCount 增加/减少评论的点赞数
func (r *commentRepository) IncrLikeCount(id uint, delta int) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

// UpdateVoteCount 更新评论的投票数（用于问答答案）
func (r *commentRepository) UpdateVoteCount(commentID uint, voteCount int) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		Update("vote_count", voteCount).Error
}
