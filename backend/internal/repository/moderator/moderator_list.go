package moderator

import (
	"context"
	"tiny-forum/internal/model/po"
)

func (r *moderatorRepository) List(ctx context.Context, page, pageSize int, boardID *uint) ([]po.Moderator, int64, error) {
	var moderators []po.Moderator
	var total int64

	query := r.db.WithContext(ctx).Model(&po.Moderator{})

	if boardID != nil {
		query = query.Where("board_id = ?", *boardID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("User").
		Preload("Board").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&moderators).Error

	return moderators, total, err
}
