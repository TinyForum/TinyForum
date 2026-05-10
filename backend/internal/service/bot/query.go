package bot

import (
	"context"
	"tiny-forum/internal/model/vo"
)

func (s *service) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*vo.BotResponse, int64, error) {
	offset := (page - 1) * pageSize
	bots, total, err := s.repo.ListByUser(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var res []*vo.BotResponse
	for _, b := range bots {
		res = append(res, s.toResponse(b))
	}
	return res, total, nil
}

func (s *service) List(ctx context.Context, page, pageSize int) ([]*vo.BotResponse, int64, error) {
	offset := (page - 1) * pageSize
	bots, total, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var res []*vo.BotResponse
	for _, b := range bots {
		res = append(res, s.toResponse(b))
	}
	return res, total, nil
}
