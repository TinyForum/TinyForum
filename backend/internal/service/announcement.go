package service

import (
	"context"
	"errors"
	"time"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrAnnouncementNotFound = errors.New("公告不存在")
	ErrInvalidPublishTime   = errors.New("发布时间无效")
	ErrExpiredTimeInvalid   = errors.New("过期时间必须晚于发布时间")
	ErrPermissionDenied     = errors.New("权限不足")
)

type AnnouncementService interface {
	Create(ctx context.Context, req *CreateAnnouncementRequest, userID uint) (*model.Announcement, error)
	Update(ctx context.Context, id uint, req *UpdateAnnouncementRequest, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
	GetByID(ctx context.Context, id uint) (*model.Announcement, error)
	List(ctx context.Context, req *ListAnnouncementRequest) (*ListAnnouncementResponse, error)
	GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error)
	Publish(ctx context.Context, id uint, userID uint) error
	Archive(ctx context.Context, id uint, userID uint) error
	Pin(ctx context.Context, id uint, pinned bool, userID uint) error
}

type CreateAnnouncementRequest struct {
	Title       string                   `json:"title" binding:"required,min=1,max=200"`
	Content     string                   `json:"content" binding:"required"`
	Summary     string                   `json:"summary"`
	Cover       string                   `json:"cover"`
	Type        model.AnnouncementType   `json:"type" binding:"required,oneof=normal important emergency event"`
	IsPinned    bool                     `json:"is_pinned"`
	IsGlobal    bool                     `json:"is_global"`
	Status      model.AnnouncementStatus `json:"status"`
	BoardID     *uint                    `json:"board_id"`
	PublishedAt *time.Time               `json:"published_at"`
	ExpiredAt   *time.Time               `json:"expired_at"`
}

type UpdateAnnouncementRequest struct {
	Title       *string                 `json:"title"`
	Content     *string                 `json:"content"`
	Summary     *string                 `json:"summary"`
	Cover       *string                 `json:"cover"`
	Type        *model.AnnouncementType `json:"type"`
	IsPinned    *bool                   `json:"is_pinned"`
	IsGlobal    *bool                   `json:"is_global"`
	BoardID     *uint                   `json:"board_id"`
	PublishedAt *time.Time              `json:"published_at"`
	ExpiredAt   *time.Time              `json:"expired_at"`
}

type ListAnnouncementRequest struct {
	Page      int                       `form:"page" binding:"min=1"`
	PageSize  int                       `form:"page_size" binding:"min=1,max=100"`
	BoardID   *uint                     `form:"board_id"`
	Type      *model.AnnouncementType   `form:"type"`
	Status    *model.AnnouncementStatus `form:"status"`
	IsPinned  *bool                     `form:"is_pinned"`
	IsGlobal  *bool                     `form:"is_global"`
	Keyword   string                    `form:"keyword"`
	StartTime *time.Time                `form:"start_time"`
	EndTime   *time.Time                `form:"end_time"`
}

type ListAnnouncementResponse struct {
	Total         int64                `json:"total"`
	Page          int                  `json:"page"`
	PageSize      int                  `json:"page_size"`
	Announcements []model.Announcement `json:"announcements"`
}

type announcementService struct {
	repo repository.AnnouncementRepository
}

func NewAnnouncementService(repo repository.AnnouncementRepository) AnnouncementService {
	return &announcementService{repo: repo}
}

func (s *announcementService) Create(ctx context.Context, req *CreateAnnouncementRequest, userID uint) (*model.Announcement, error) {
	// 验证时间
	if err := s.validateTime(req.PublishedAt, req.ExpiredAt); err != nil {
		return nil, err
	}

	now := time.Now()
	announcement := &model.Announcement{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Cover:       req.Cover,
		Type:        req.Type,
		IsPinned:    req.IsPinned,
		IsGlobal:    req.IsGlobal,
		BoardID:     req.BoardID,
		PublishedAt: req.PublishedAt,
		ExpiredAt:   req.ExpiredAt,
		Status:      model.AnnouncementStatusDraft,
		ViewCount:   0,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	// 如果没有设置发布时间，默认为当前时间
	if announcement.PublishedAt == nil {
		announcement.PublishedAt = &now
	}

	if err := s.repo.Create(ctx, announcement); err != nil {
		return nil, err
	}

	return announcement, nil
}

func (s *announcementService) Update(ctx context.Context, id uint, req *UpdateAnnouncementRequest, userID uint) error {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAnnouncementNotFound
		}
		return err
	}

	// 更新字段
	if req.Title != nil {
		announcement.Title = *req.Title
	}
	if req.Content != nil {
		announcement.Content = *req.Content
	}
	if req.Summary != nil {
		announcement.Summary = *req.Summary
	}
	if req.Cover != nil {
		announcement.Cover = *req.Cover
	}
	if req.Type != nil {
		announcement.Type = *req.Type
	}
	if req.IsPinned != nil {
		announcement.IsPinned = *req.IsPinned
	}
	if req.IsGlobal != nil {
		announcement.IsGlobal = *req.IsGlobal
	}
	if req.BoardID != nil {
		announcement.BoardID = req.BoardID
	}
	if req.PublishedAt != nil {
		announcement.PublishedAt = req.PublishedAt
	}
	if req.ExpiredAt != nil {
		announcement.ExpiredAt = req.ExpiredAt
	}

	// 验证时间
	if err := s.validateTime(announcement.PublishedAt, announcement.ExpiredAt); err != nil {
		return err
	}

	announcement.UpdatedBy = userID

	return s.repo.Update(ctx, announcement)
}

func (s *announcementService) Delete(ctx context.Context, id uint, userID uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAnnouncementNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *announcementService) GetByID(ctx context.Context, id uint) (*model.Announcement, error) {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAnnouncementNotFound
		}
		return nil, err
	}

	// 异步增加浏览次数
	go s.repo.IncrementViewCount(context.Background(), id)

	return announcement, nil
}

func (s *announcementService) List(ctx context.Context, req *ListAnnouncementRequest) (*ListAnnouncementResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	repoReq := &repository.AnnouncementListRequest{
		Page:      req.Page,
		PageSize:  req.PageSize,
		BoardID:   req.BoardID,
		Type:      req.Type,
		Status:    req.Status,
		IsPinned:  req.IsPinned,
		IsGlobal:  req.IsGlobal,
		Keyword:   req.Keyword,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	announcements, total, err := s.repo.List(ctx, repoReq)
	if err != nil {
		return nil, err
	}

	return &ListAnnouncementResponse{
		Total:         total,
		Page:          req.Page,
		PageSize:      req.PageSize,
		Announcements: announcements,
	}, nil
}

func (s *announcementService) GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error) {
	return s.repo.GetPinned(ctx, boardID)
}

func (s *announcementService) Publish(ctx context.Context, id uint, userID uint) error {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if announcement.PublishedAt == nil || announcement.PublishedAt.After(time.Now()) {
		return ErrInvalidPublishTime
	}

	return s.repo.UpdateStatus(ctx, id, model.AnnouncementStatusPublished)
}

func (s *announcementService) Archive(ctx context.Context, id uint, userID uint) error {
	return s.repo.UpdateStatus(ctx, id, model.AnnouncementStatusArchived)
}

func (s *announcementService) Pin(ctx context.Context, id uint, pinned bool, userID uint) error {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	announcement.IsPinned = pinned
	return s.repo.Update(ctx, announcement)
}

func (s *announcementService) validateTime(publishedAt, expiredAt *time.Time) error {
	if publishedAt != nil && expiredAt != nil {
		if !expiredAt.After(*publishedAt) {
			return ErrExpiredTimeInvalid
		}
	}
	return nil
}
