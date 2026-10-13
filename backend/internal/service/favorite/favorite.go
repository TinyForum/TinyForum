package favorite

import (
	"context"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	repository "tiny-forum/internal/repository/favorite"
	apperrors "tiny-forum/pkg/errors"
)

// normalizeTargetType 校验并规范化收藏目标类型，默认值为 creation
func normalizeTargetType(targetType string) (string, error) {
	if targetType == "" {
		return repository.FavoriteTargetTypeCreation, nil
	}
	switch targetType {
	case repository.FavoriteTargetTypeCreation:
		return targetType, nil
	default:
		return "", apperrors.ErrInvalidRequest.WithDetail("不支持的收藏对象类型")
	}
}

// AddFavorite 收藏目标内容（幂等友好：重复收藏返回业务错误）
func (s *favoriteService) AddFavorite(ctx context.Context, userID, targetID uint, targetType string, groupID int64) error {
	if userID == 0 || targetID == 0 {
		return apperrors.ErrInvalidRequest
	}
	tt, err := normalizeTargetType(targetType)
	if err != nil {
		return err
	}

	exists, err := s.repo.Exists(ctx, int64(userID), int64(targetID), tt)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.ErrCollectAlready
	}

	return s.repo.Create(ctx, &do.Favorite{
		UserID:     int64(userID),
		TargetID:   int64(targetID),
		TargetType: tt,
		GroupID:    groupID,
	})
}

// RemoveFavorite 取消收藏
func (s *favoriteService) RemoveFavorite(ctx context.Context, userID, targetID uint, targetType string) error {
	if userID == 0 || targetID == 0 {
		return apperrors.ErrInvalidRequest
	}
	tt, err := normalizeTargetType(targetType)
	if err != nil {
		return err
	}

	rows, err := s.repo.Delete(ctx, int64(userID), int64(targetID), tt)
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperrors.ErrCollectNotExist
	}
	return nil
}

// ListMyFavorites 分页获取当前用户的收藏
func (s *favoriteService) ListMyFavorites(ctx context.Context, req request.ListFavoritesRequest, userID uint) (*common.PageResult[vo.UserFavorite], error) {
	tt, err := normalizeTargetType(req.TargetType)
	if err != nil {
		return nil, err
	}

	rows, total, err := s.repo.ListByUser(ctx, int64(userID), tt, req.Offset(), req.Limit())
	if err != nil {
		return nil, err
	}

	list := make([]vo.UserFavorite, 0, len(rows))
	for _, r := range rows {
		list = append(list, vo.UserFavorite{
			ID:        r.FavoriteID,
			CreatedAt: r.CreatedAt,
			TargetID:  r.TargetID,
			Title:     r.Title,
			Summary:   r.Summary,
			LikeCount: r.LikeCount,
			BoardName: r.BoardName,
		})
	}

	return &common.PageResult[vo.UserFavorite]{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		HasMore:  int64(req.Offset()+len(list)) < total,
	}, nil
}
