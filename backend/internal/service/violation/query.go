package violation

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (s *violationService) ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error) {
	return s.repo.ListUserViolationByUserID(ctx, req, userID)
}
