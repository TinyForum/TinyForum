package article

import "tiny-forum/internal/model/do"

func (r *articleRepository) GetByBoardID(boardID uint, limit, offset int) ([]do.Article, int64, error) {
	var posts []do.Article
	var total int64

	query := r.db.Model(&do.Article{}).
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

func (r *articleRepository) TogglePinInBoard(postID uint, pin bool) error {
	return r.db.Model(&do.Article{}).Where("id = ?", postID).
		Update("pin_in_board", pin).Error
}
