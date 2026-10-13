package violation

import (
	"context"
	"errors"
	"time"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

func (s *violationService) ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error) {
	return s.repo.ListUserViolationByUserID(ctx, req, userID)
}

// CountUserViolationByUserID 统计用户违规记录总数（与列表同过滤条件）
func (s *violationService) CountUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) (int64, error) {
	return s.repo.CountUserViolationByUserID(ctx, req, userID)
}

// GetViolation 获取用户自己的违规记录详情
func (s *violationService) GetViolation(ctx context.Context, id, userID uint) (*do.Violation, error) {
	v, err := s.repo.GetViolationByUserAndID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithDetail("违规记录不存在")
		}
		return nil, err
	}
	return v, nil
}

// AppealViolation 提交违规申诉（状态机：only none/rejected 可再次申诉）
func (s *violationService) AppealViolation(ctx context.Context, id, userID uint, reason string) (*do.Violation, error) {
	v, err := s.GetViolation(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	switch v.AppealStatus {
	case do.AppealStatusPending:
		return nil, apperrors.ErrInvalidRequest.WithDetail("申诉处理中，请勿重复提交")
	case do.AppealStatusApproved:
		return nil, apperrors.ErrInvalidRequest.WithDetail("该违规已申诉通过")
	}
	now := time.Now()
	v.AppealStatus = do.AppealStatusPending
	v.AppealReason = reason
	v.AppealTime = &now
	if err := s.repo.UpdateViolation(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}
