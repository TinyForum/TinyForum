package violation

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"

	"gorm.io/gorm"
)

// baseQuery 构建用户违规记录的基础查询（含状态、关键词过滤；默认扣除软删记录）
func (r *violationRepository) baseQuery(ctx context.Context, req request.ListUserViolationRequest, userID uint) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&do.Violation{}).Where("user_id = ?", userID)

	if req.Status != "" {
		q = q.Where("status = ?", req.Status)
	}
	if req.Keyword != "" {
		q = q.Where("reason LIKE ?", "%"+req.Keyword+"%")
	}
	return q
}

func (r *violationRepository) ListUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) ([]*do.Violation, error) {
	var violations []*do.Violation

	query := r.baseQuery(ctx, req, userID)

	if req.PageSize > 0 {
		query = query.Offset(req.Offset()).Limit(req.Limit())
	}

	err := query.Order("created_at DESC").Find(&violations).Error
	if err != nil {
		return nil, err
	}
	return violations, nil
}

// CountUserViolationByUserID 统计用户违规记录总数（与列表使用相同过滤条件），用于正确分页。
func (r *violationRepository) CountUserViolationByUserID(ctx context.Context, req request.ListUserViolationRequest, userID uint) (int64, error) {
	var total int64
	err := r.baseQuery(ctx, req, userID).Count(&total).Error
	return total, err
}

// GetViolationByUserAndID 查询指定用户的违规记录
func (r *violationRepository) GetViolationByUserAndID(ctx context.Context, id, userID uint) (*do.Violation, error) {
	var v do.Violation
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// UpdateViolation 更新违规记录（用于申诉）
func (r *violationRepository) UpdateViolation(ctx context.Context, violation *do.Violation) error {
	return r.db.WithContext(ctx).Save(violation).Error
}