package board

import (
	"errors"
	"time"
	"tiny-forum/internal/model/do"
)

type BanUserInput struct {
	UserID    uint       `json:"user_id"  binding:"required"`
	BoardID   uint       `json:"board_id" binding:"required"`
	Reason    string     `json:"reason"   binding:"required,max=500"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (s *boardService) BanUser(input BanUserInput, bannerID uint) error {
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
