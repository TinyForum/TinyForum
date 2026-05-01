package moderator

import (
	"context"
	"tiny-forum/internal/model/do"
)

func (r *moderatorRepository) DeleteByBoard(ctx context.Context, boardID uint) error {
	return r.db.WithContext(ctx).
		Where("board_id = ?", boardID).
		Delete(&do.Moderator{}).Error
}

func (r *moderatorRepository) DeleteByUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&do.Moderator{}).Error
}
