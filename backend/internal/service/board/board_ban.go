package board

import (
	"context"
	"errors"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (s *boardService) BanUser(ctx context.Context, input request.BoardBanUserRequest, bannerID uint) error {
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}
	isBanned, _ := s.boardRepo.IsBanned(input.UserID, input.BoardID)
	if isBanned {
		return errors.New("用户已被禁言")
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
