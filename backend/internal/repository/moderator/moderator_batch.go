package moderator

import (
	"context"
	"tiny-forum/internal/model/po"
)

func (r *moderatorRepository) DeleteByBoard(ctx context.Context, boardID uint) error {
	return r.db.WithContext(ctx).
		Where("board_id = ?", boardID).
		Delete(&po.Moderator{}).Error
}

func (r *moderatorRepository) DeleteByUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&po.Moderator{}).Error
}
