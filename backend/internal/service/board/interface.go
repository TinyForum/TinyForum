package board

import (
	"context"

	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	postRepo "tiny-forum/internal/repository/article"
	boardRepo "tiny-forum/internal/repository/board"
	commentRepo "tiny-forum/internal/repository/comment"
	reportsRepo "tiny-forum/internal/repository/reports"
	"tiny-forum/internal/repository/transaction"
	userRepo "tiny-forum/internal/repository/user"
	attachment "tiny-forum/internal/service/attachment"
	"tiny-forum/internal/service/notification"
)

type BoardService interface {
	// applys
	ApplyModerator(ctx context.Context, input request.ApplyModeratorRequest) error
	CancelApplication(applicationID, userID uint) error
	GetUserApplications(userID uint, page, pageSize int) ([]do.ModeratorApplication, int64, error)
	ReviewApplication(_ context.Context, input request.ReviewApplicationRequest, reviewerID uint) error
	ListApplications(boardID *uint, status do.ApplicationStatus, page, pageSize int) ([]do.ModeratorApplication, int64, error)
	ApplyBoardCreation(ctx context.Context, input request.ApplyBoardCreationRequest) error
	GetUserBoardApplications(userID uint, page, pageSize int) ([]do.BoardCreationApplication, int64, error)
	ListBoardApplications(status do.ApplicationStatus, page, pageSize int) ([]do.BoardCreationApplication, int64, error)
	ReviewBoardCreationApplication(ctx context.Context, input request.ReviewBoardCreationRequest, reviewerID uint) error
	// ban
	BanUser(ctx context.Context, input request.BoardBanUserRequest, bannerID uint) error
	UnbanUser(userID, boardID uint) error
	IsBanned(userID, boardID uint) (bool, error)
	// crud
	Create(input bo.CreateBoardInput) (*do.Board, error)
	Update(id uint, input bo.CreateBoardInput) (*do.Board, error)
	Delete(id uint) error
	GetByID(id uint) (*do.Board, error)
	GetBoardBySlug(slug string) (*do.Board, error)
	GetPostsBySlug(slug string, page, pageSize int) ([]*dto.GetBoardPostsResponse, int64, error)
	List(page, pageSize int) ([]do.Board, int64, error)
	GetTree() ([]do.Board, error)
	GetPosts(boardID uint, page, pageSize int) ([]do.Article, int64, error)
	// moderator
	AddModerator(_ context.Context, input request.AddModeratorRequest, operatorID uint) error
	RemoveModerator(_ context.Context, userID, boardID uint, operatorID uint) error
	GetModerators(boardID uint) ([]do.Moderator, error)
	IsModerator(userID, boardID uint) (bool, error)
	UpdateModeratorPermissions(_ context.Context, input request.UpdateModeratorPermissionsRequest, operatorID uint) error
	CheckModeratorPermission(_ context.Context, userID, boardID uint, permission do.ModeratorPermission) (bool, error)
	GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardWithPerms, error)
	// post
	DeletePost(boardID, postID, userID uint, isAdmin bool) error
	PinPost(boardID, postID uint, pin bool) error
	// moderator panel
	ListModeratorPosts(ctx context.Context, boardID uint, req *request.ListModeratorPostsRequest) ([]vo.ModeratorPostVO, int64, error)   // 版主面板：板块帖子列表
	ListModeratorReports(ctx context.Context, boardID uint, req *request.ListModeratorReportsRequest) ([]vo.ModeratorReportVO, int64, error) // 版主面板：板块举报列表
	ListModeratorBans(ctx context.Context, boardID uint, req *request.ListModeratorBansRequest) ([]vo.BannedUserVO, int64, error)           // 版主面板：板块封禁列表
}
type boardService struct {
	boardRepo     boardRepo.BoardRepository
	userRepo      userRepo.UserRepository
	postRepo      postRepo.ArticleRepository
	commentRepo   commentRepo.CommentRepository
	reportRepo    reportsRepo.ReportsRepository
	notifSvc      notification.NotificationService
	attachmentSvc attachment.AttachmentService
	txManager     transaction.TransactionManager
}

func NewBoardService(
	boardRepo boardRepo.BoardRepository,
	userRepo userRepo.UserRepository,
	postRepo postRepo.ArticleRepository,
	commentRepo commentRepo.CommentRepository,
	reportRepo reportsRepo.ReportsRepository,
	notifSvc notification.NotificationService,
	attachmentSvc attachment.AttachmentService,
	txManager transaction.TransactionManager,
) BoardService {
	return &boardService{
		boardRepo:     boardRepo,
		userRepo:      userRepo,
		postRepo:      postRepo,
		commentRepo:   commentRepo,
		reportRepo:    reportRepo,
		notifSvc:      notifSvc,
		attachmentSvc: attachmentSvc,
		txManager:     txManager,
	}
}
