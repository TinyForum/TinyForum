package tag

import (
	"context"
	"time"
	"tiny-forum/internal/model/po"
)

type TagRepository interface {
	Create(tag *po.Tag) error
	FindByID(id uint) (*po.Tag, error)
	FindByName(name string) (*po.Tag, error)
	List() ([]po.Tag, error)
	Update(tag *po.Tag) error
	Delete(id uint) error
	IncrPostCount(id uint, delta int) error
	// post
	FindTagsByPostIDs(postIDs []uint) (map[uint][]po.Tag, error)
	FindTagsByPostID(postID uint) ([]po.Tag, error)
	// stats
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
}
