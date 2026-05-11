package violation

import (
	"gorm.io/gorm"
)

type ViolationRepository interface {
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
