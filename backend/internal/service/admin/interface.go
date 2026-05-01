package admin

import (
	"context"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
	"tiny-forum/internal/model/query"
	"tiny-forum/internal/model/vo"
	announcementSvc "tiny-forum/internal/service/announcement"
	userSvc "tiny-forum/internal/service/user"
)

type AdminService interface {
	ListAnnouncements(ctx context.Context, req *query.ListAnnouncements) (*vo.ListAnnouncements, error)
	CreateAnnouncement(ctx context.Context, req *dto.CreateAnnouncementRequest, userID uint) (*po.Announcement, error)
	UpdateAnnouncement(ctx context.Context, id uint, req *dto.UpdateAnnouncementRequest, userID uint) error
	DeleteAnnouncement(ctx context.Context, id uint, userID uint) error
	PublishAnnouncement(ctx context.Context, id uint, userID uint) error
	ArchiveAnnouncement(ctx context.Context, id uint, userID uint) error
	PinAnnouncement(ctx context.Context, id uint, pinned bool, userID uint) error
	// users
	ListUsers(page, pageSize int, keyword string) ([]po.User, int64, error)
	SetActiveUser(targetID uint, operatorID uint, isActive bool) error
	SetBlockedUser(targetID uint, operatorID uint, isBlocked bool) error
	DeleteUser(operatorID uint, targetID uint) error
	SetRoleUser(operatorID, targetID uint, newRole string) error
}

type adminService struct {
	// commentRepo     commentRepo.CommentRepository
	// postRepo        postRepo.PostRepository
	// userRepo        userRepo.UserRepository
	// voteRepo        voteRepo.VoteRepository
	announcementSvc announcementSvc.AnnouncementService
	userSvc         userSvc.UserService
}

func NewAdminService(
	// commentRepo commentRepo.CommentRepository,
	// postRepo postRepo.PostRepository,
	// userRepo userRepo.UserRepository,
	announcementSvc announcementSvc.AnnouncementService,
	// voteRepo voteRepo.VoteRepository,
	userSvc userSvc.UserService,
) AdminService {
	return &adminService{
		// 	postRepo:        postRepo,
		// 	userRepo:        userRepo,
		// 	voteRepo:        voteRepo,
		// 	commentRepo:     commentRepo,
		announcementSvc: announcementSvc,
		userSvc:         userSvc,
	}
}
