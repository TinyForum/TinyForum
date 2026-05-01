package admin

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

func (s *adminService) ListAnnouncements(ctx context.Context, req *request.ListAnnouncements) (*vo.ListAnnouncements, error) {
	return s.announcementSvc.List(ctx, req)
}

func (s *adminService) CreateAnnouncement(ctx context.Context, req *request.CreateAnnouncement, userID uint) (*do.Announcement, error) {
	return s.announcementSvc.Create(ctx, req, userID)
}

func (s *adminService) UpdateAnnouncement(ctx context.Context, id uint, req *request.UpdateAnnouncement, userID uint) error {
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
