package violation

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"

	"gorm.io/gorm"
)

type ViolationRepository interface {
	ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error)
	// CountUserViolationByUserID 统计用户违规记录总数（与 ListUserViolationByUserID 使用相同的过滤条件）
	CountUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) (int64, error)
	// GetViolationByUserAndID 查询指定用户的违规记录
	GetViolationByUserAndID(ctx context.Context, id, userID uint) (*do.Violation, error)
	// UpdateViolation 更新违规记录（用于申诉）
	UpdateViolation(ctx context.Context, violation *do.Violation) error
}

type violationRepository struct {
	db *gorm.DB
}

func NewViolationRepository(db *gorm.DB) ViolationRepository {
	return &violationRepository{
		db: db,
	}
}
