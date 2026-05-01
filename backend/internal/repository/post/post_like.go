package post

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *postRepository) IncrViewCount(id uint) error {
	return r.db.Model(&do.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *postRepository) IncrLikeCount(id uint, delta int) error {
	return r.db.Model(&do.Post{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

func (r *postRepository) AddLike(userID, postID uint) error {
	like := &do.Like{UserID: userID, PostID: &postID}
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		FirstOrCreate(like).Error
}

func (r *postRepository) RemoveLike(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&do.Like{}).Error
}

func (r *postRepository) IsLiked(userID, postID uint) bool {
	var count int64
	r.db.Model(&do.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}
