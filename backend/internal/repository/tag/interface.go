package tag

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
)

type TagRepository interface {
	Create(tag *do.Tag) error
	FindByID(id uint) (*do.Tag, error)
	FindByName(name string) (*do.Tag, error)
	List() ([]do.Tag, error)
	Update(tag *do.Tag) error
	Delete(id uint) error
	IncrPostCount(id uint, delta int) error
	// post
	FindTagsByPostIDs(postIDs []uint) (map[uint][]do.Tag, error)
	FindTagsByPostID(postID uint) ([]do.Tag, error)
	// stats
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
}
