package announcement

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"

	"gorm.io/gorm"
)

// AnnouncementRepository 公告仓库接口
type AnnouncementRepository interface {
	Create(ctx context.Context, announcement *do.Announcement) error
	Update(ctx context.Context, announcement *do.Announcement) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*do.Announcement, error)
	List(ctx context.Context, req *request.ListAnnouncements) ([]do.Announcement, int64, error)
	GetPinned(ctx context.Context, boardID *uint) ([]do.Announcement, error)
	IncrementViewCount(ctx context.Context, id uint) error
	BatchDelete(ctx context.Context, ids []uint) error
	UpdateStatus(ctx context.Context, id uint, status do.AnnouncementStatus) error
}

// AnnouncementListRequest 列表查询参数
// type AnnouncementListRequest struct {
// 	Page      int
// 	PageSize  int
// 	BoardID   *uint
// 	Type      *do.AnnouncementType
// 	Status    *do.AnnouncementStatus
// 	IsPinned  *bool
// 	IsGlobal  *bool
// 	Keyword   string
// 	StartTime *time.Time
// 	EndTime   *time.Time
// }

// 状态常量（用于查询过滤）
// const (
// 	StatusAll       = "all"
// 	StatusDraft     = "draft"
// 	StatusPublished = "published"
// 	StatusExpired   = "expired"
// )

// announcementRepository 实现
type announcementRepository struct {
	db *gorm.DB
}

// NewAnnouncementRepository 构造函数
func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}
