package user

import (
	"context"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

func (s *userService) ListUserViolation(ctx context.Context, req request.ListUserViolationRequest, userID uint) (*common.PageResult[vo.ViolationVO], error) {
	// 1. 获取违规记录列表
	violations, err := s.violationSvc.ListUserViolationByUserID(ctx, req, userID)
	if err != nil {
		return nil, err
	}

	total, err := s.violationSvc.CountUserViolationByUserID(ctx, req, userID)
	if err != nil {
		return nil, err
	}

	// 2. 转换为 VO
	voList := make([]vo.ViolationVO, 0, len(violations))
	for _, v := range violations {
		voList = append(voList, s.convertViolationToVO(v))
	}

	return &common.PageResult[vo.ViolationVO]{
		List:     voList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  int64(req.Page*req.PageSize) < total,
	}, nil
}

// convertViolationToVO 单个违规记录转换（提取复用）
func (s *userService) convertViolationToVO(v *do.Violation) vo.ViolationVO {
	var operatorIDPtr *uint
	if v.OperatorID != 0 {
		operatorIDPtr = &v.OperatorID
	}

	return vo.ViolationVO{
		ID:             v.ID,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
		UserID:         v.UserID,
		OperatorID:     operatorIDPtr,
		ViolationType:  v.ViolationType,
		Reason:         v.Reason,
		Source:         v.Source,
		Status:         v.Status,
		PunishType:     v.PunishType,
		PunishExpireAt: v.PunishExpireAt,
		AppealStatus:   v.AppealStatus,
		AppealReason:   v.AppealReason,
		AppealTime:     v.AppealTime,
		AppealResult:   v.AppealResult,
	}
}

// GetUserViolation 获取当前用户的违规记录详情
func (s *userService) GetUserViolation(ctx context.Context, id, userID uint) (*vo.ViolationVO, error) {
	v, err := s.violationSvc.GetViolation(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	ret := s.convertViolationToVO(v)
	return &ret, nil
}

// AppealUserViolation 当前用户对违规记录提交申诉
func (s *userService) AppealUserViolation(ctx context.Context, id, userID uint, reason string) (*vo.ViolationVO, error) {
	v, err := s.violationSvc.AppealViolation(ctx, id, userID, reason)
	if err != nil {
		return nil, err
	}
	ret := s.convertViolationToVO(v)
	return &ret, nil
}
