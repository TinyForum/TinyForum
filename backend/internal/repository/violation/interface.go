package violation

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"

	"gorm.io/gorm"
)

type ViolationRepository interface {
	ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error)
	// CreateViolation(ctx context.Context, violation *model.Violation) (*model.Violation, error)
	// GetViolation(ctx context.Context, id string) (*model.Violation, error)
	// UpdateViolation(ctx context.Context, violation *model.Violation) (*model.Violation, error)
	// DeleteViolation(ctx context.Context, id string) error
}

type violationRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) ViolationRepository {
	return &violationRepository{
		db: db,
	}
}
