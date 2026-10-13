package violation

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/repository/violation"
)

type ViolationService interface {
	ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error)
	// CountUserViolationByUserID 统计用户违规记录总数
	CountUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) (int64, error)
	// GetViolation 获取用户自己的违规记录详情
	GetViolation(ctx context.Context, id, userID uint) (*do.Violation, error)
	// AppealViolation 提交违规申诉
	AppealViolation(ctx context.Context, id, userID uint, reason string) (*do.Violation, error)
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
