package favorite

import (
	"context"

	"tiny-forum/internal/model/common"
	repo "tiny-forum/internal/repository/favorite"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

// FavoriteService 收藏业务接口
type FavoriteService interface {
	// AddFavorite 收藏目标内容
	AddFavorite(ctx context.Context, userID, targetID uint, targetType string, groupID int64) error
	// RemoveFavorite 取消收藏
	RemoveFavorite(ctx context.Context, userID, targetID uint, targetType string) error
	// ListMyFavorites 分页获取当前用户的收藏
	ListMyFavorites(ctx context.Context, req request.ListFavoritesRequest, userID uint) (*common.PageResult[vo.UserFavorite], error)
}

type favoriteService struct {
	repo repo.FavoriteRepository
}

// NewFavoriteService 创建收藏服务
func NewFavoriteService(repo repo.FavoriteRepository) FavoriteService {
	return &favoriteService{repo: repo}
}
