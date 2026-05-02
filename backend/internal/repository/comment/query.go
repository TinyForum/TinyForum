package comment

import (
	"context"
	"tiny-forum/internal/model/do"
)

func (r *commentRepository) BatchCountByPostIDs(ctx context.Context, postIDs []uint) (map[uint]int64, error) {
	type Result struct {
		PostID uint
		Count  int64
	}
	var results []Result
	err := r.db.WithContext(ctx).
		Model(&do.Comment{}).
		Select("post_id, COUNT(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	countMap := make(map[uint]int64, len(results))
	for _, r := range results {
		countMap[r.PostID] = r.Count
	}
	return countMap, nil
}
