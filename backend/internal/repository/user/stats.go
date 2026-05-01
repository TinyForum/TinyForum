package user

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"

	"golang.org/x/sync/errgroup"
)

func (r *userRepository) GetGlobalStatsCount(ctx context.Context, userID uint) (*dto.GlobalStatsCount, error) {
	var (
		stats dto.GlobalStatsCount
		eg    errgroup.Group
	)
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&do.Post{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalPosts = int(count)
		return nil
	})

	// 总评论数
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&do.Comment{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})

	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&do.Favorite{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&do.Violation{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})
	return nil, nil

}
