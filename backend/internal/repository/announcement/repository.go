package announcement

import (
	"context"
	"tiny-forum/internal/model/po"
	"tiny-forum/internal/model/query"

	"gorm.io/gorm"
)

// AnnouncementRepository 公告仓库接口
type AnnouncementRepository interface {
	Create(ctx context.Context, announcement *po.Announcement) error
	Update(ctx context.Context, announcement *po.Announcement) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*po.Announcement, error)
	List(ctx context.Context, req *query.ListAnnouncements) ([]po.Announcement, int64, error)
	GetPinned(ctx context.Context, boardID *uint) ([]po.Announcement, error)
	IncrementViewCount(ctx context.Context, id uint) error
	BatchDelete(ctx context.Context, ids []uint) error
	UpdateStatus(ctx context.Context, id uint, status po.AnnouncementStatus) error
}

// AnnouncementListRequest 列表查询参数
// type AnnouncementListRequest struct {
// 	Page      int
// 	PageSize  int
// 	BoardID   *uint
// 	Type      *po.AnnouncementType
// 	Status    *po.AnnouncementStatus
// 	IsPinned  *bool
// 	IsGlobal  *bool
// 	Keyword   string
// 	StartTime *time.Time
// 	EndTime   *time.Time
// }

// 状态常量（用于查询过滤）
const (
	StatusAll       = "all"
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusExpired   = "expired"
)

// announcementRepository 实现
type announcementRepository struct {
	db *gorm.DB
}

// NewAnnouncementRepository 构造函数
func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}
