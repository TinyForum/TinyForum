package user

import (
	"context"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"

	"golang.org/x/sync/errgroup"
)

func (r *userRepository) GetGlobalStatsCount(ctx context.Context, userID uint) (*dto.GlobalStatsCount, error) {
	var (
		stats dto.GlobalStatsCount
		eg    errgroup.Group
	)
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&model.Post{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalPosts = int(count)
		return nil
	})

	// 总评论数
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&model.Comment{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})

	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&model.Favorite{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&model.Violation{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})
	return nil, nil

}
