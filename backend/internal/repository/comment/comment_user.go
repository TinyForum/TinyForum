package comment

import (
	"context"

	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// ListByAuthor 查询指定用户发表的评论（依附于 Reply），返回评论记录与总数。
// 通过 replies.author_id 关联，过滤软删除的回复，按回复创建时间倒序。
func (r *commentRepository) ListByAuthor(ctx context.Context, authorID uint, offset, limit int) ([]do.Comment, int64, error) {
	var comments []do.Comment
	var total int64

	base := func() *gorm.DB {
		return r.db.WithContext(ctx).
			Model(&do.Comment{}).
			Joins("JOIN replies ON replies.id = comments.reply_id").
			Where("replies.author_id = ? AND replies.deleted_at IS NULL", authorID)
	}

	if err := base().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := base().
		Preload("Reply.Author").
		Order("replies.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&comments).Error

	return comments, total, err
}

// UpdateReplyStatusByCommentID 更新评论对应回复的显示状态（审核通过=visible，拒绝=hidden）
func (r *commentRepository) UpdateReplyStatusByCommentID(ctx context.Context, commentID uint, status do.ReplyStatus) error {
	var comment do.Comment
	if err := r.db.WithContext(ctx).First(&comment, commentID).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Model(&do.Reply{}).
		Where("id = ?", comment.ReplyID).
		Update("status", status).Error
}
