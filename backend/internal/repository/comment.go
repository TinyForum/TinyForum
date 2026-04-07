package repository

import (
	"bbs-forum/internal/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("Author").First(&comment, id).Error
	return &comment, err
}

func (r *CommentRepository) ListByPost(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{}).Where("post_id = ? AND parent_id IS NULL", postID)

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

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

func (r *CommentRepository) Update(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepository) IncrLikeCount(id uint, delta int) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

func (r *CommentRepository) CountByPost(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}
