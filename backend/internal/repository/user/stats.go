package user

import (
	"context"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"

	"golang.org/x/sync/errgroup"
)

func (r *userRepository) GetGlobalStatsCount(ctx context.Context, userID uint) (*dto.GlobalStatsCount, error) {
	var (
		stats dto.GlobalStatsCount
		eg    errgroup.Group
	)
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&po.Post{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalPosts = int(count)
		return nil
	})

	// 总评论数
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&po.Comment{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})

	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&po.Favorite{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})
	eg.Go(func() error {
		var count int64
		if err := r.db.WithContext(ctx).Model(&po.Violation{}).Count(&count).Error; err != nil {
			return err
		}
		stats.TotalComments = int(count)
		return nil
	})
	return nil, nil

}
