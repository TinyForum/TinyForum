// package repository

// import (
// 	"context"
// 	"time"
// 	"tiny-forum/internal/model"

// 	"gorm.io/gorm"
// )

// type AnnouncementRepository interface {
// 	Create(ctx context.Context, announcement *model.Announcement) error
// 	Update(ctx context.Context, announcement *model.Announcement) error
// 	Delete(ctx context.Context, id uint) error
// 	GetByID(ctx context.Context, id uint) (*model.Announcement, error)
// 	List(ctx context.Context, req *AnnouncementListRequest) ([]model.Announcement, int64, error)
// 	GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error)
// 	IncrementViewCount(ctx context.Context, id uint) error
// 	BatchDelete(ctx context.Context, ids []uint) error
// 	UpdateStatus(ctx context.Context, id uint, status model.AnnouncementStatus) error
// }

// type AnnouncementListRequest struct {
// 	Page      int
// 	PageSize  int
// 	BoardID   *uint
// 	Type      *model.AnnouncementType
// 	Status    *model.AnnouncementStatus
// 	IsPinned  *bool
// 	IsGlobal  *bool
// 	Keyword   string
// 	StartTime *time.Time
// 	EndTime   *time.Time
// }

// type announcementRepository struct {
// 	db *gorm.DB
// }

// func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
// 	return &announcementRepository{db: db}
// }

// func (r *announcementRepository) Create(ctx context.Context, announcement *model.Announcement) error {
// 	return r.db.WithContext(ctx).Create(announcement).Error
// }

// func (r *announcementRepository) Update(ctx context.Context, announcement *model.Announcement) error {
// 	return r.db.WithContext(ctx).Save(announcement).Error
// }

// func (r *announcementRepository) Delete(ctx context.Context, id uint) error {
// 	return r.db.WithContext(ctx).Delete(&model.Announcement{}, id).Error
// }

// func (r *announcementRepository) GetByID(ctx context.Context, id uint) (*model.Announcement, error) {
// 	var announcement model.Announcement
// 	err := r.db.WithContext(ctx).
// 		Preload("Board").
// 		Preload("Creator").
// 		First(&announcement, id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &announcement, nil
// }

// const (
// 	AnnouncementStatusAll       = "all" // 查询所有状态
// 	AnnouncementStatusDraft     = "draft"
// 	AnnouncementStatusPublished = "published"
// 	AnnouncementStatusExpired   = "expired"
// )

// func (r *announcementRepository) List(ctx context.Context, req *AnnouncementListRequest) ([]model.Announcement, int64, error) {
// 	var announcements []model.Announcement
// 	var total int64

// 	query := r.db.WithContext(ctx).Model(&model.Announcement{})

// 	// 条件过滤
// 	if req.BoardID != nil {
// 		query = query.Where("board_id = ? OR (board_id IS NULL AND is_global = ?)", *req.BoardID, true)
// 	}
// 	if req.Type != nil {
// 		query = query.Where("type = ?", *req.Type)
// 	}
// 	if req.Status != nil {
// 		// 如果是 "all"，不添加 status 过滤
// 		if *req.Status != AnnouncementStatusAll {
// 			query = query.Where("status = ?", *req.Status)
// 		}
// 	}
// 	if req.IsPinned != nil {
// 		query = query.Where("is_pinned = ?", *req.IsPinned)
// 	}
// 	if req.IsGlobal != nil {
// 		query = query.Where("is_global = ?", *req.IsGlobal)
// 	}
// 	if req.Keyword != "" {
// 		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
// 	}
// 	if req.StartTime != nil {
// 		query = query.Where("published_at >= ?", req.StartTime)
// 	}
// 	if req.EndTime != nil {
// 		query = query.Where("published_at <= ?", req.EndTime)
// 	}

// 	// 时间过滤：只有查询 published 状态或未指定状态时才添加时间过滤
// 	// 如果明确指定了 status="all" 或其他非 published 状态，则不过滤时间
// 	shouldFilterTime := req.Status == nil ||
// 		(req.Status != nil && *req.Status == AnnouncementStatusPublished)

// 	if shouldFilterTime {
// 		query = query.Where("published_at <= ?", time.Now())
// 		query = query.Where("expired_at IS NULL OR expired_at > ?", time.Now())
// 	}

// 	// 统计总数
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	// 分页和排序
// 	offset := (req.Page - 1) * req.PageSize
// 	err := query.
// 		Order("is_pinned DESC").
// 		Order("published_at DESC").
// 		Preload("Board").
// 		Preload("Creator").
// 		Offset(offset).
// 		Limit(req.PageSize).
// 		Find(&announcements).Error

// 	return announcements, total, err
// }

// func (r *announcementRepository) GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error) {
// 	var announcements []model.Announcement
// 	query := r.db.WithContext(ctx).
// 		Where("is_pinned = ?", true).
// 		Where("status = ?", model.AnnouncementStatusPublished).
// 		Where("published_at <= ?", time.Now()).
// 		Where("expired_at IS NULL OR expired_at > ?", time.Now())

// 	if boardID != nil {
// 		query = query.Where("board_id = ? OR (board_id IS NULL AND is_global = ?)", *boardID, true)
// 	} else {
// 		query = query.Where("is_global = ?", true)
// 	}

// 	err := query.Order("published_at DESC").
// 		Preload("Board").
// 		Find(&announcements).Error

// 	return announcements, err
// }

// func (r *announcementRepository) IncrementViewCount(ctx context.Context, id uint) error {
// 	return r.db.WithContext(ctx).
// 		Model(&model.Announcement{}).
// 		Where("id = ?", id).
// 		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
// }

// func (r *announcementRepository) BatchDelete(ctx context.Context, ids []uint) error {
// 	return r.db.WithContext(ctx).Delete(&model.Announcement{}, ids).Error
// }

// func (r *announcementRepository) UpdateStatus(ctx context.Context, id uint, status model.AnnouncementStatus) error {
// 	return r.db.WithContext(ctx).
// 		Model(&model.Announcement{}).
// 		Where("id = ?", id).
// 		Update("status", status).Error
// }

package repository
