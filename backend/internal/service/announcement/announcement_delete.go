package announcement

import (
	"context"
	"errors"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

func (s *announcementService) Delete(ctx context.Context, id uint, userID uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrAnnouncementNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
