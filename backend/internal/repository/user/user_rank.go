package user

import (
	"context"
	"fmt"
	"tiny-forum/internal/model"
)

type TopUsersQuery struct {
	Limit          int
	ExcludeBlocked bool
	Fields         []string
}

func (r *UserRepository) GetTopUsers(ctx context.Context, query TopUsersQuery) ([]model.User, error) {
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	dbQuery := r.db.WithContext(ctx).Model(&model.User{})

	if len(query.Fields) > 0 {
		dbQuery = dbQuery.Select(query.Fields)
	}
	if query.ExcludeBlocked {
		dbQuery = dbQuery.Where("is_blocked = ?", false)
	}

	var users []model.User
	err := dbQuery.
		Order("COALESCE(score, 0) DESC").
		Limit(query.Limit).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("获取排行榜失败: %w", err)
	}
	if users == nil {
		return []model.User{}, nil
	}
	return users, nil
}

func (r *UserRepository) GetTopFollowers(limit int) ([]model.User, error) {
	var users []model.User
	err := r.db.Table("users").
		Select("users.*, COUNT(follows.follower_id) as follower_count").
		Joins("LEFT JOIN follows ON users.id = follows.following_id").
		Group("users.id").
		Order("follower_count DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}
