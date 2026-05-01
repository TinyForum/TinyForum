package user

import (
	"context"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
)

/*
*
- @Description 获取积分最高的用户
- @Param ctx 上下文
- @Param limit 获取的用户数量
- @Param excludeBlocked 是否排除被拉黑的
- @Return []po.User 用户列表
- @Return user,error 错误信息
*
*/
func (r *userRepository) GetTopScoreUsersSimple(ctx context.Context, limit int, excludeBlocked bool) ([]dto.LeaderboardUserSimple, error) {
	var users []dto.LeaderboardUserSimple
	db := r.db.WithContext(ctx).Model(&po.User{}).
		Where("score IS NOT NULL").
		Select("id", "username", "avatar", "score").
		Order("score DESC").
		Limit(limit)
	if excludeBlocked {
		db = db.Where("is_blocked = ?", false)
	}
	err := db.Find(&users).Error
	return users, err
}

// GetTopScoreUsersDetail 返回排行榜详细信息（用于 /leaderboard/detail 端点）
func (r *userRepository) GetTopScoreUsersDetail(ctx context.Context, limit int, excludeBlocked bool) ([]dto.LeaderboardUserDetail, error) {

	var users []dto.LeaderboardUserDetail
	db := r.db.WithContext(ctx).Model(&po.User{}).
		Where("score IS NOT NULL").
		Select("id", "username", "avatar", "score", "email", "role").
		Order("score DESC").
		Limit(limit)
	if excludeBlocked {
		db = db.Where("is_blocked = ?", false)
	}
	err := db.Find(&users).Error
	return users, err
}
