package comment

import (
	"tiny-forum/internal/model/do"
)

// ListByPost 获取帖子的顶级评论（带分页和子评论预加载）
func (r *commentRepository) ListByPost(postID uint, page, pageSize int) ([]do.Comment, int64, error) {
	var comments []do.Comment
	var total int64

	query := r.db.Model(&do.Comment{}).
		Joins("JOIN replies ON replies.id = comments.reply_id").
		Where("comments.works_id = ? AND replies.parent_id IS NULL", postID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := r.db.Model(&do.Comment{}).
		Joins("JOIN replies ON replies.id = comments.reply_id").
		Where("comments.works_id = ? AND replies.parent_id IS NULL", postID).
		Preload("Reply.Author").
		Preload("Reply.Replies.Author").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	return comments, total, err
}

// GetAnswersByPostID 获取帖子的所有答案（按采纳、投票、创建时间排序）
func (r *commentRepository) GetAnswersByQuestionID(questionID uint, limit, offset int) ([]do.Answer, int64, error) {
	var comments []do.Answer
	var total int64

	query := r.db.Model(&do.Answer{}).
		Where("question_id = ?", questionID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Reply.Author").
		Order("is_accepted DESC, vote_count DESC, created_at ASC").
		Find(&comments).Error

	return comments, total, err
}

// GetAnswersByPostIDOrderByNewest 按最新排序获取答案
func (r *commentRepository) GetAnswersByPostIDOrderByNewest(postID uint, limit, offset int) ([]do.Answer, int64, error) {
	var comments []do.Answer
	var total int64

	query := r.db.Model(&do.Answer{}).
		Where("creations_id = ? AND is_answer = ?", postID, true)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("created_at DESC").
		Find(&comments).Error

	return comments, total, err
}

// GetAnswersByPostIDOrderByOldest 按最早排序获取答案
func (r *commentRepository) GetAnswersByPostIDOrderByOldest(postID uint, limit, offset int) ([]do.Answer, int64, error) {
	var comments []do.Answer
	var total int64

	query := r.db.Model(&do.Answer{}).
		Where("creations_id = ? AND is_answer = ?", postID, true)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error

	return comments, total, err
}
