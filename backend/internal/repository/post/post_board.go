package post

import "tiny-forum/internal/model/do"

func (r *postRepository) GetByBoardID(boardID uint, limit, offset int) ([]do.Post, int64, error) {
	var posts []do.Post
	var total int64

	query := r.db.Model(&do.Post{}).
		Where("board_id = ? AND status = ?", boardID, do.PostStatusPublished)

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
	return r.db.Model(&do.Post{}).Where("id = ?", postID).
		Update("pin_in_board", pin).Error
}
