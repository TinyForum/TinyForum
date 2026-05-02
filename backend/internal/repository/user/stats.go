package user

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"

	"golang.org/x/sync/errgroup"
)

// func (r *userRepository) GetGlobalStatsCount(ctx context.Context, userID uint) (*dto.GlobalStatsCount, error) {
// 	var (
// 		stats dto.GlobalStatsCount
// 		eg    errgroup.Group
// 	)
// 	eg.Go(func() error {
// 		var count int64
// 		if err := r.db.WithContext(ctx).Model(&do.Post{}).Count(&count).Error; err != nil {
// 			return err
// 		}
// 		stats.TotalCountPosts = int(count)
// 		return nil
// 	})

// 	// 总评论数
// 	eg.Go(func() error {
// 		var count int64
// 		if err := r.db.WithContext(ctx).Model(&do.Comment{}).Count(&count).Error; err != nil {
// 			return err
// 		}
// 		stats.TotalCountComments = int(count)
// 		return nil
// 	})

// 	eg.Go(func() error {
// 		var count int64
// 		if err := r.db.WithContext(ctx).Model(&do.Favorite{}).Count(&count).Error; err != nil {
// 			return err
// 		}
// 		stats.TotalCountFavorites = int(count)
// 		return nil
// 	})
// 	eg.Go(func() error {
// 		var count int64
// 		if err := r.db.WithContext(ctx).Model(&do.Violation{}).Count(&count).Error; err != nil {
// 			return err
// 		}
// 		stats.TotalCountViolation = int(count)
// 		return nil
// 	})
// 	return nil, nil

// }

// 获取用户统计
func (r *userRepository) GetGlobalStatsCount(ctx context.Context, userID uint) (*dto.StatsInfo, error) {
	var stats dto.StatsInfo

	// 使用 errgroup 并发统计各类数量
	eg, ctx := errgroup.WithContext(ctx)

	// 用户发帖数
	eg.Go(func() error {
		var count int64
		err := r.db.WithContext(ctx).Model(&do.Post{}).Where("author_id = ?", userID).Count(&count).Error
		stats.TotalPost = int64(count)
		return err
	})

	// 用户评论数
	eg.Go(func() error {
		var count int64
		err := r.db.WithContext(ctx).Model(&do.Comment{}).Where("author_id = ?", userID).Count(&count).Error
		stats.TotalComment = int64(count)
		return err
	})

	// 用户收到的点赞数（对用户所有帖子和评论的点赞）
	eg.Go(func() error {
		// 方式：联合查询 Post 和 Comment 的 author_id
		var likeCount int64
		// 子查询：帖子点赞
		postLikeSub := r.db.WithContext(ctx).Model(&do.Like{}).
			Joins("JOIN posts ON likes.post_id = posts.id").
			Where("posts.author_id = ?", userID).
			Select("COUNT(*)")
		// 子查询：评论点赞
		commentLikeSub := r.db.WithContext(ctx).Model(&do.Like{}).
			Joins("JOIN comments ON likes.comment_id = comments.id").
			Where("comments.author_id = ?", userID).
			Select("COUNT(*)")

		// 执行两个计数并相加
		var postLikeCount, commentLikeCount int64
		if err := postLikeSub.Scan(&postLikeCount).Error; err != nil {
			return err
		}
		if err := commentLikeSub.Scan(&commentLikeCount).Error; err != nil {
			return err
		}
		likeCount = postLikeCount + commentLikeCount
		stats.TotalLike = int64(likeCount)
		return nil
	})

	// 用户违规数
	eg.Go(func() error {
		var count int64
		err := r.db.WithContext(ctx).Model(&do.Violation{}).Where("user_id = ?", userID).Count(&count).Error
		stats.TotalViolation = int64(count)
		return err
	})

	// 用户收藏数
	eg.Go(func() error {
		var count int64
		err := r.db.WithContext(ctx).Model(&do.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
		stats.TotalFavorite = int64(count)
		return err
	})

	// 用户关注数（我关注的人）
	eg.Go(func() error {
		var count int64
		err := r.db.WithContext(ctx).Model(&do.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
		stats.TotalFollowing = int64(count)
		return err
	})

	// 用户粉丝数（关注我的人）
	eg.Go(func() error {
		var count int64
		err := r.db.WithContext(ctx).Model(&do.Follow{}).Where("following_id = ?", userID).Count(&count).Error
		stats.TotalFollowing = int64(count)
		return err
	})

	// 用户积分（直接从 User 表中取）
	eg.Go(func() error {
		var user do.User
		err := r.db.WithContext(ctx).Select("score").Where("id = ?", userID).First(&user).Error
		if err == nil {
			stats.TotalScore = user.Score
		}
		return err // 即使无记录也返回错误
	})

	if err := eg.Wait(); err != nil {
		// 如果用户不存在，可能 first 报错，可单独处理
		return nil, err
	}

	return &stats, nil
}
