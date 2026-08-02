package user

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
)

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&do.User{}).Count(&count).Error
	return count, err
}

func (r *userRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&do.User{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

func (r *userRepository) CountActiveByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("users u").
		Where(`u.deleted_at IS NULL AND EXISTS (
			SELECT 1 FROM articles a
			JOIN creations cr ON cr.id = a.creation_id AND cr.deleted_at IS NULL
			WHERE cr.author_id = u.id AND a.deleted_at IS NULL
			  AND a.created_at BETWEEN ? AND ?
		) OR EXISTS (
			SELECT 1 FROM comments co
			JOIN replies r ON r.id = co.reply_id AND r.deleted_at IS NULL
			WHERE r.author_id = u.id AND co.deleted_at IS NULL
			  AND co.created_at BETWEEN ? AND ?
		)`, startDate, endDate, startDate, endDate).
		Count(&count).Error
	return count, err
}

func (r *userRepository) GetActiveUsersByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*vo.ActiveUserRowVO, error) {
	var rows []*vo.ActiveUserRowVO
	err := r.db.WithContext(ctx).
		Table("users u").
		Select("u.id, u.username, u.avatar_url").
		Where(`u.deleted_at IS NULL AND (
			EXISTS (
				SELECT 1 FROM articles a
				JOIN creations cr ON cr.id = a.creation_id AND cr.deleted_at IS NULL
				WHERE cr.author_id = u.id AND a.deleted_at IS NULL
				  AND a.created_at BETWEEN ? AND ?
			) OR EXISTS (
				SELECT 1 FROM comments co
				JOIN replies r ON r.id = co.reply_id AND r.deleted_at IS NULL
				WHERE r.author_id = u.id AND co.deleted_at IS NULL
				  AND co.created_at BETWEEN ? AND ?
			)
		)`, startDate, endDate, startDate, endDate).
		Order("u.score DESC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}
