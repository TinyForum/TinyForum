package admin

import (
	"context"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
	"tiny-forum/internal/model/query"
	"tiny-forum/internal/model/vo"
)

func (s *adminService) ListAnnouncements(ctx context.Context, req *query.ListAnnouncements) (*vo.ListAnnouncements, error) {
	return s.announcementSvc.List(ctx, req)
}

func (s *adminService) CreateAnnouncement(ctx context.Context, req *dto.CreateAnnouncementRequest, userID uint) (*po.Announcement, error) {
	return s.announcementSvc.Create(ctx, req, userID)
}

func (s *adminService) UpdateAnnouncement(ctx context.Context, id uint, req *dto.UpdateAnnouncementRequest, userID uint) error {
	return s.announcementSvc.Update(ctx, id, req, userID)
}
func (s *adminService) DeleteAnnouncement(ctx context.Context, id uint, userID uint) error {
	return s.announcementSvc.Delete(ctx, id, userID)
}
func (s *adminService) PublishAnnouncement(ctx context.Context, id uint, userID uint) error {
	return s.announcementSvc.Publish(ctx, id, userID)
}

func (s *adminService) ArchiveAnnouncement(ctx context.Context, id uint, userID uint) error {
	return s.announcementSvc.Archive(ctx, id, userID)
}

func (s *adminService) PinAnnouncement(ctx context.Context, id uint, pinned bool, userID uint) error {
	return s.announcementSvc.Pin(ctx, id, pinned, userID)
}
