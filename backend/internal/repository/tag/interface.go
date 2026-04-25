package tag

import (
	"context"
	"time"
	"tiny-forum/internal/model"
)

type TagRepository interface {
	Create(tag *model.Tag) error
	FindByID(id uint) (*model.Tag, error)
	FindByName(name string) (*model.Tag, error)
	List() ([]model.Tag, error)
	Update(tag *model.Tag) error
	Delete(id uint) error
	IncrPostCount(id uint, delta int) error
	// post
	FindTagsByPostIDs(postIDs []uint) (map[uint][]model.Tag, error)
	FindTagsByPostID(postID uint) ([]model.Tag, error)
	// stats
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
}
