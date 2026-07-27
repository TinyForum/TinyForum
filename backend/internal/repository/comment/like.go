package comment

import (
	"errors"
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

func (r *commentRepository) IncrViewCount(id uint) error {
	return r.db.Model(&do.Creation{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *commentRepository) IncrLikeCount(id uint, delta int) error {
	return r.db.Model(&do.Creation{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

func (r *commentRepository) AddLike(userID, articleID uint) error {
	// 检查是否已存在
	var existing do.Like
	err := r.db.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, articleID, do.LikeTargetPost).First(&existing).Error
	if err == nil {
		return apperrors.ErrLikeAlready
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // 其他数据库错误
	}
	// 不存在，创建
	like := &do.Like{UserID: userID, TargetType: do.LikeTargetPost, TargetID: articleID}
	return r.db.Create(like).Error
}

func (r *commentRepository) RemoveLike(userID, postID uint) error {
	// 1. 检查点赞记录是否存在
	var like do.Like
	err := r.db.Where("user_id = ? AND target_id = ? AND target_type = ?",
		userID, postID, do.LikeTargetPost).
		First(&like).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperrors.ErrLikeNotExist // 已取消点赞
	}
	if err != nil {
		return err // 其他数据库错误
	}

	// 2. 存在则执行删除
	return r.db.Delete(&like).Error
}
func (r *commentRepository) IsLiked(userID, postID uint) bool {
	var count int64
	r.db.Model(&do.Like{}).Where("user_id = ? AND target_id = ? AND target_type = ?", userID, postID, do.LikeTargetPost).Count(&count)
	return count > 0
}
