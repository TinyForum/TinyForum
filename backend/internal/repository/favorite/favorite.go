package favorite

import (
	"context"

	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏数据访问实现
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

// Create 新增收藏记录
func (r *favoriteRepository) Create(ctx context.Context, fav *do.Favorite) error {
	return r.db.WithContext(ctx).Create(fav).Error
}

// Delete 按用户与目标软删除收藏记录，返回受影响行数
func (r *favoriteRepository) Delete(ctx context.Context, userID, targetID int64, targetType string) (int64, error) {
	tx := r.db.WithContext(ctx).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Delete(&do.Favorite{})
	return tx.RowsAffected, tx.Error
}

// Exists 判断用户是否已收藏目标（未软删记录）
func (r *favoriteRepository) Exists(ctx context.Context, userID, targetID int64, targetType string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&do.Favorite{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListByUser 分页查询用户收藏的作品，联表 creations / boards 获取展示信息
func (r *favoriteRepository) ListByUser(ctx context.Context, userID int64, targetType string, offset, limit int) ([]FavoriteRow, int64, error) {
	var total int64

	base := func() *gorm.DB {
		return r.db.WithContext(ctx).
			Table("favorites AS f").
			Joins("JOIN creations AS c ON c.id = f.target_id AND c.deleted_at IS NULL").
			Where("f.user_id = ? AND f.target_type = ? AND f.deleted_at IS NULL", userID, targetType)
	}

	if err := base().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []FavoriteRow
	err := base().
		Select("f.id AS favorite_id, f.created_at AS created_at, f.target_id AS target_id, " +
			"c.title AS title, c.summary AS summary, c.like_count AS like_count, " +
			"COALESCE(b.name, '') AS board_name").
		Joins("LEFT JOIN boards AS b ON b.id = c.board_id AND b.deleted_at IS NULL").
		Order("f.created_at DESC").
		Offset(offset).Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
