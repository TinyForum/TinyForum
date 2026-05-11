package user

import (
	"context"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

func (s *userService) ListUserViolation(ctx context.Context, req request.ListUserViolationRequest, userID uint) (*common.PageResult[vo.ViolationVO], error) {
	// 1. 获取违规记录列表
	violations, err := s.violationSvc.ListUserViolationByUserID(ctx, req, userID)
	if err != nil {
		return nil, err
	}

	// 2. 获取总记录数（分页用）

	if err != nil {
		return nil, err
	}
	total := int64(len(violations))

	// 3. 转换为 VO
	voList := make([]vo.ViolationVO, 0, len(violations))
	for _, v := range violations {
		voList = append(voList, vo.ViolationVO{
			ID:             v.ID,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			UserID:         v.UserID,
			OperatorID:     v.OperatorID,
			ViolationType:  string(v.ViolationType),
			Reason:         v.Reason,
			Source:         string(v.Source),
			Status:         string(v.Status),
			PunishType:     string(v.PunishType),
			PunishExpireAt: v.PunishExpireAt,
			AppealStatus:   string(v.AppealStatus),
			AppealReason:   v.AppealReason,
			AppealTime:     v.AppealTime,
			AppealResult:   v.AppealResult,
		})
	}

	return &common.PageResult[vo.ViolationVO]{
		List:     voList,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  int64(req.Page*req.PageSize) < total,
	}, nil
}
