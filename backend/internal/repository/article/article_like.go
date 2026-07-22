package article

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *articleRepository) IncrViewCount(id uint) error {
	return r.db.Model(&do.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *articleRepository) IncrLikeCount(id uint, delta int) error {
	return r.db.Model(&do.Article{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

func (r *articleRepository) AddLike(userID, postID uint) error {
	like := &do.Like{UserID: userID, TargetType: do.LikeTargetPost, TargetID: postID}
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		FirstOrCreate(like).Error
}

func (r *articleRepository) RemoveLike(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&do.Like{}).Error
}

func (r *articleRepository) IsLiked(userID, postID uint) bool {
	var count int64
	r.db.Model(&do.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}
