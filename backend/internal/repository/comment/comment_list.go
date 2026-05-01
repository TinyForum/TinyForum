package comment

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// ListByPost 获取帖子的顶级评论（带分页和子评论预加载）
func (r *commentRepository) ListByPost(postID uint, page, pageSize int) ([]do.Comment, int64, error) {
	var comments []do.Comment
	var total int64

	query := r.db.Model(&do.Comment{}).Where("post_id = ? AND parent_id IS NULL", postID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("Author").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Author").Order("created_at ASC")
		}).
		Order("created_at ASC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	return comments, total, err
}

// GetAnswersByPostID 获取帖子的所有答案（按采纳、投票、创建时间排序）
func (r *commentRepository) GetAnswersByPostID(postID uint, limit, offset int) ([]do.Comment, int64, error) {
	var comments []do.Comment
	var total int64

	query := r.db.Model(&do.Comment{}).
		Where("post_id = ? AND is_answer = ?", postID, true)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("is_accepted DESC, vote_count DESC, created_at ASC").
		Find(&comments).Error

	return comments, total, err
}

// GetAnswersByPostIDOrderByNewest 按最新排序获取答案
func (r *commentRepository) GetAnswersByPostIDOrderByNewest(postID uint, limit, offset int) ([]do.Comment, int64, error) {
	var comments []do.Comment
	var total int64

	query := r.db.Model(&do.Comment{}).
		Where("post_id = ? AND is_answer = ?", postID, true)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("created_at DESC").
		Find(&comments).Error

	return comments, total, err
}

// GetAnswersByPostIDOrderByOldest 按最早排序获取答案
func (r *commentRepository) GetAnswersByPostIDOrderByOldest(postID uint, limit, offset int) ([]do.Comment, int64, error) {
	var comments []do.Comment
	var total int64

	query := r.db.Model(&do.Comment{}).
		Where("post_id = ? AND is_answer = ?", postID, true)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error

	return comments, total, err
}
