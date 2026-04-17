package comment

import "tiny-forum/internal/model"

// MarkAsAccepted 标记评论为已采纳答案
func (r *CommentRepository) MarkAsAccepted(commentID uint) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		Update("is_accepted", true).Error
}

// MarkAsAnswer 标记/取消标记评论为答案
func (r *CommentRepository) MarkAsAnswer(commentID uint, isAnswer bool) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		Update("is_answer", isAnswer).Error
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
