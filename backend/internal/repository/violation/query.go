package violation

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (r *violationRepository) ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error) {
	var violations []*do.Violation

	// 基础查询：根据用户ID过滤
	query := r.db.WithContext(ctx).Model(&do.Violation{}).Where("user_id = ?", userID)

	// 可选：支持分页
	if req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	// 可选：按时间倒序排列
	err := query.Order("created_at DESC").Find(&violations).Error
	if err != nil {
		return nil, err
	}
	return violations, nil
}
