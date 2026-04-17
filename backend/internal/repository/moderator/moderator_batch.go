package moderator

import (
	"context"
	"tiny-forum/internal/model"
)

func (r *moderatorRepository) DeleteByBoard(ctx context.Context, boardID uint) error {
	return r.db.WithContext(ctx).
		Where("board_id = ?", boardID).
		Delete(&model.Moderator{}).Error
}

func (r *moderatorRepository) DeleteByUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Moderator{}).Error
}
