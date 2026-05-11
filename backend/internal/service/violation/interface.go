package violation

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/repository/violation"
)

type ViolationService interface {
	ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error)
	// CreateViolation(ctx context.Context, violation *model.Violation) (*model.Violation, error)
	// GetViolation(ctx context.Context, id string) (*model.Violation, error)
	// UpdateViolation(ctx context.Context, violation *model.Violation) (*model.Violation, error)
	// DeleteViolation(ctx context.Context, id string) error
}

type violationService struct {
	repo violation.ViolationRepository
}

func NewViolationService(
	repo violation.ViolationRepository,
) ViolationService {
	return &violationService{
		repo: repo,
	}
}
