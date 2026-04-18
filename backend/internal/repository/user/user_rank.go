package user

import (
	"context"
	"tiny-forum/internal/model"
)

/**
 * @Description 获取积分最高的用户
 * @Param ctx 上下文
 * @Param limit 获取的用户数量
 * @Param excludeBlocked 是否排除被拉黑的
 * @Return []model.User 用户列表
 * @Return error 错误信息
 **/
func (r *UserRepository) GetTopScoreUsers(ctx context.Context, limit int, excludeBlocked bool) ([]model.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	db := r.db.WithContext(ctx).Where("score IS NOT NULL")
	if excludeBlocked {
		db = db.Where("is_blocked = ?", false)
	}

	var users []model.User
	err := db.
		Order("score DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}
