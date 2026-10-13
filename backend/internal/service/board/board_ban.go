package board

import (
	"context"
	"errors"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

func (s *boardService) BanUser(ctx context.Context, input request.BoardBanUserRequest, bannerID uint) error {
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}
	isBanned, err := s.boardRepo.IsBanned(input.UserID, input.BoardID)
	if err != nil {
		return err
	}
	if isBanned {
		return apperrors.ErrInvalidRequest.WithDetail("用户已被禁言")
	}
	ban := &do.BoardBan{
		UserID:    input.UserID,
		BoardID:   input.BoardID,
		BannedBy:  bannerID,
		Reason:    input.Reason,
		ExpiresAt: input.ExpiresAt,
	}
	if err := s.boardRepo.BanUser(ban); err != nil {
		return err
	}
	s.notifSvc.Create(user.ID, &bannerID, do.NotifySystem,
		"你在板块中被禁言", &input.BoardID, "board")
	return nil
}

func (s *boardService) UnbanUser(userID, boardID uint) error {
	return s.boardRepo.UnbanUser(userID, boardID)
}

func (s *boardService) IsBanned(userID, boardID uint) (bool, error) {
	return s.boardRepo.IsBanned(userID, boardID)
}
