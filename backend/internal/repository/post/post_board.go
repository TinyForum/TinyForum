package post

import (
	"tiny-forum/internal/model/po"
)

func (r *postRepository) GetByBoardID(boardID uint, limit, offset int) ([]po.Post, int64, error) {
	var posts []po.Post
	var total int64

	query := r.db.Model(&po.Post{}).
		Where("board_id = ? AND status = ?", boardID, po.PostStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Preload("Tags").
		Preload("Board").
		Order("pin_top DESC, pin_in_board DESC, created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) TogglePinInBoard(postID uint, pin bool) error {
	return r.db.Model(&po.Post{}).Where("id = ?", postID).
		Update("pin_in_board", pin).Error
}
