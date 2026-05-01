package admin

import (
	"context"
	"tiny-forum/internal/model/query"
	"tiny-forum/internal/model/vo"
	announcementRepo "tiny-forum/internal/repository/announcement"
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	userRepo "tiny-forum/internal/repository/user"
	voteRepo "tiny-forum/internal/repository/vote"
)

type AdminService interface {
	ListAnnouncements(ctx context.Context, req *query.ListAnnouncements) (*vo.ListAnnouncements, error)
}

// emailService 邮件服务实现（私有）
type adminService struct {
	commentRepo      commentRepo.CommentRepository
	postRepo         postRepo.PostRepository
	userRepo         userRepo.UserRepository
	voteRepo         voteRepo.VoteRepository
	announcementRepo announcementRepo.AnnouncementRepository
}

func NewAdminService(commentRepo commentRepo.CommentRepository,
	postRepo postRepo.PostRepository,
	userRepo userRepo.UserRepository,
	announcementRepo announcementRepo.AnnouncementRepository,
	voteRepo voteRepo.VoteRepository) AdminService {
	return &adminService{
		postRepo:         postRepo,
		userRepo:         userRepo,
		voteRepo:         voteRepo,
		commentRepo:      commentRepo,
		announcementRepo: announcementRepo,
	}
}
