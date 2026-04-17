package moderator

import (
	"context"
	"errors"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

func (r *moderatorRepository) Create(ctx context.Context, moderator *model.Moderator) error {
	return r.db.WithContext(ctx).Create(moderator).Error
}

func (r *moderatorRepository) Update(ctx context.Context, moderator *model.Moderator) error {
	return r.db.WithContext(ctx).Save(moderator).Error
}

func (r *moderatorRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Moderator{}, id).Error
}

func (r *moderatorRepository) GetByID(ctx context.Context, id uint) (*model.Moderator, error) {
	var moderator model.Moderator
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Board").
		First(&moderator, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &moderator, nil
}
