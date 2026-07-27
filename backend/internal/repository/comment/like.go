package comment

import (
	"errors"
	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

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

func (r *commentRepository) AddLike(userID, commentID uint) error {
	var existing do.Like
	// 使用 Unscoped() 忽略软删除过滤，查找所有记录（包括已删除）
	err := r.db.Unscoped().Where("user_id = ? AND target_id = ? AND target_type = ?", userID, commentID, do.LikeTargetComment).First(&existing).Error
	logger.Info("查询点赞记录")
	if err == nil {
		logger.Info("查询到点赞记录")
		// 如果找到记录，判断是否被软删除
		if existing.BaseModel.DeletedAt.Valid { // 或 existing.DeletedAt != nil
			// 恢复（更新 deleted_at 为 NULL）
			logger.Info("恢复点赞记录")
			return r.db.Unscoped().Model(&existing).Update("deleted_at", nil).Error
		}
		logger.Info("已存在且未删除，提示已点赞")
		return apperrors.ErrLikeAlready // 已存在且未删除，提示已点赞
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Info("查询点赞记录失败")
		return err
	}
	// 确实不存在，创建新记录
	like := &do.Like{UserID: userID, TargetType: do.LikeTargetComment, TargetID: commentID}
	return r.db.Create(like).Error
}

func (r *commentRepository) RemoveLike(userID, commentID uint) error {
	// 1. 检查点赞记录是否存在
	var like do.Like
	err := r.db.Where("user_id = ? AND target_id = ? AND target_type = ?",
		userID, commentID, do.LikeTargetComment).
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
	r.db.Model(&do.Like{}).Where("user_id = ? AND target_id = ? AND target_type = ?", userID, postID, do.LikeTargetComment).Count(&count)
	return count > 0
}
