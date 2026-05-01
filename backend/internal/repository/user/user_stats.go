package user

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
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
			SELECT 1 FROM posts p
			WHERE p.author_id = u.id AND p.deleted_at IS NULL
			  AND p.created_at BETWEEN ? AND ?
		) OR EXISTS (
			SELECT 1 FROM comments c
			WHERE c.author_id = u.id AND c.deleted_at IS NULL
			  AND c.created_at BETWEEN ? AND ?
		)`, startDate, endDate, startDate, endDate).
		Count(&count).Error
	return count, err
}

type ActiveUserRow struct {
	ID       uint
	Username string
	Avatar   string
}

func (r *userRepository) GetActiveUsersByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*ActiveUserRow, error) {
	var rows []*ActiveUserRow
	err := r.db.WithContext(ctx).
		Table("users u").
		Select("u.id, u.username, u.avatar").
		Where(`u.deleted_at IS NULL AND (
			EXISTS (
				SELECT 1 FROM posts p
				WHERE p.author_id = u.id AND p.deleted_at IS NULL
				  AND p.created_at BETWEEN ? AND ?
			) OR EXISTS (
				SELECT 1 FROM comments c
				WHERE c.author_id = u.id AND c.deleted_at IS NULL
				  AND c.created_at BETWEEN ? AND ?
			)
		)`, startDate, endDate, startDate, endDate).
		Order("u.score DESC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}
