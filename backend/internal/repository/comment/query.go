package comment

import (
	"context"
	"tiny-forum/internal/model/do"
)

func (r *commentRepository) BatchCountByPostIDs(ctx context.Context, postIDs []uint) (map[uint]int64, error) {
	type Result struct {
		CreationsID uint
		Count       int64
	}
	var results []Result
	err := r.db.WithContext(ctx).
		Model(&do.Comment{}).
		Select("creations_id, COUNT(*) as count").
		Where("creations_id IN ?", postIDs).
		Group("creations_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	countMap := make(map[uint]int64, len(results))
	for _, r := range results {
		countMap[r.CreationsID] = r.Count
	}
	return countMap, nil
}
