package moderator

import (
	"context"
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *moderatorRepository) GetByUserAndBoard(ctx context.Context, userID, boardID uint) (*do.Moderator, error) {
	var moderator do.Moderator
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND board_id = ?", userID, boardID).
		First(&moderator).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &moderator, nil
}

func (r *moderatorRepository) GetByBoard(ctx context.Context, boardID uint) ([]do.Moderator, error) {
	var moderators []do.Moderator
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("board_id = ?", boardID).
		Find(&moderators).Error
	return moderators, err
}

func (r *moderatorRepository) GetByUser(ctx context.Context, userID uint) ([]do.Moderator, error) {
	var moderators []do.Moderator
	err := r.db.WithContext(ctx).
		Preload("Board").
		Where("user_id = ?", userID).
		Find(&moderators).Error
	return moderators, err
}

func (r *moderatorRepository) Exists(ctx context.Context, userID, boardID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&do.Moderator{}).
		Where("user_id = ? AND board_id = ?", userID, boardID).
		Count(&count).Error
	return count > 0, err
}

func (r *moderatorRepository) IsModerator(ctx context.Context, userID, boardID uint) (bool, error) {
	return r.Exists(ctx, userID, boardID)
}
