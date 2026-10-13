package favorite

import (
	"context"
	"time"

	"tiny-forum/internal/model/do"
)

// FavoriteTargetTypeCreation 收藏对象类型：作品（帖子/文章）
const FavoriteTargetTypeCreation = "creation"

// FavoriteRow 收藏列表查询结果（联表 creations / boards）
type FavoriteRow struct {
	FavoriteID uint
	CreatedAt  time.Time
	TargetID   uint
	Title      string
	Summary    string
	LikeCount  int
	BoardName  string
}

// FavoriteRepository 收藏数据访问接口
type FavoriteRepository interface {
	// Create 新增一条收藏记录
	Create(ctx context.Context, fav *do.Favorite) error
	// Delete 按用户与目标软删除收藏记录；未命中返回 ErrNotExist 语义由调用方判断（此处返回删除影响行数）
	Delete(ctx context.Context, userID, targetID int64, targetType string) (int64, error)
	// Exists 判断某用户是否已收藏某目标
	Exists(ctx context.Context, userID, targetID int64, targetType string) (bool, error)
	// ListByUser 分页查询某用户的收藏记录（联表作品信息）
	ListByUser(ctx context.Context, userID int64, targetType string, offset, limit int) ([]FavoriteRow, int64, error)
}
