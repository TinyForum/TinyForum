package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// AnnouncementRepository 公告仓库接口
type AnnouncementRepository interface {
	Create(ctx context.Context, announcement *model.Announcement) error
	Update(ctx context.Context, announcement *model.Announcement) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.Announcement, error)
	List(ctx context.Context, req *AnnouncementListRequest) ([]model.Announcement, int64, error)
	GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error)
	IncrementViewCount(ctx context.Context, id uint) error
	BatchDelete(ctx context.Context, ids []uint) error
	UpdateStatus(ctx context.Context, id uint, status model.AnnouncementStatus) error
}

// AnnouncementListRequest 列表查询参数
type AnnouncementListRequest struct {
	Page      int
	PageSize  int
	BoardID   *uint
	Type      *model.AnnouncementType
	Status    *model.AnnouncementStatus
	IsPinned  *bool
	IsGlobal  *bool
	Keyword   string
	StartTime *time.Time
	EndTime   *time.Time
}

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
