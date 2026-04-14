package repository

import (
	"context"
	"errors"
	"tiny-forum/internal/model"

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

func (r *CommentRepository) ValidateParentComment(parentID uint, postID uint) error {
	var comment model.Comment
	err := r.db.Where("id = ? AND post_id = ?", parentID, postID).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("父评论不存在或不属于当前帖子")
		}
		return err
	}
	return nil
}

// 获取指定评论的子评论
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

// 在现有 CommentRepository 中添加以下方法

func (r *CommentRepository) MarkAsAccepted(commentID uint) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		Update("is_accepted", true).Error
}

func (r *CommentRepository) MarkAsAnswer(commentID uint, isAnswer bool) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		Update("is_answer", isAnswer).Error
}

func (r *CommentRepository) GetAnswersByPostID(postID uint, limit, offset int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND is_answer = ?", postID, true)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("is_accepted DESC, vote_count DESC, created_at ASC").
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepository) UpdateVoteCount(commentID uint, voteCount int) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		Update("vote_count", voteCount).Error
}

// 在 CommentRepository 中添加以下方法

// GetAnswersByPostIDOrderByNewest 按最新排序获取答案
func (r *CommentRepository) GetAnswersByPostIDOrderByNewest(postID uint, limit, offset int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND is_answer = ?", postID, true)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("created_at DESC").
		Find(&comments).Error

	return comments, total, err
}

// GetAnswersByPostIDOrderByOldest 按最早排序获取答案
func (r *CommentRepository) GetAnswersByPostIDOrderByOldest(postID uint, limit, offset int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND is_answer = ?", postID, true)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error

	return comments, total, err
}

// internal/repository/comment_repository.go

// Count 获取评论总数
func (r *CommentRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Count(&count).Error
	return count, err
}

// CountByDateRange 根据日期范围统计评论数
func (r *CommentRepository) CountByDateRange(ctx context.Context, startDate, endDate string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}
