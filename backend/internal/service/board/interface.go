package board

import (
	"context"

	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	boardRepo "tiny-forum/internal/repository/board"
	postRepo "tiny-forum/internal/repository/post"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
)

type BoardService interface {
	// applys
	ApplyModerator(input model.ApplyModeratorInput) error
	CancelApplication(applicationID, userID uint) error
	GetUserApplications(userID uint, page, pageSize int) ([]model.ModeratorApplication, int64, error)
	ReviewApplication(_ context.Context, input ReviewApplicationInput, reviewerID uint) error
	ListApplications(boardID *uint, status model.ApplicationStatus, page, pageSize int) ([]model.ModeratorApplication, int64, error)
	// ban
	BanUser(input BanUserInput, bannerID uint) error
	UnbanUser(userID, boardID uint) error
	IsBanned(userID, boardID uint) (bool, error)
	// crud
	Create(input CreateBoardInput) (*model.Board, error)
	Update(id uint, input CreateBoardInput) (*model.Board, error)
	Delete(id uint) error
	GetByID(id uint) (*model.Board, error)
	GetBoardBySlug(slug string) (*model.Board, error)
	GetPostsBySlug(slug string, page, pageSize int) ([]*dto.GetBoardPostsResponse, int64, error)
	List(page, pageSize int) ([]model.Board, int64, error)
	GetTree() ([]model.Board, error)
	GetPosts(boardID uint, page, pageSize int) ([]model.Post, int64, error)
	// moderator
	AddModerator(_ context.Context, input AddModeratorInput, operatorID uint) error
	RemoveModerator(_ context.Context, userID, boardID uint, operatorID uint) error
	GetModerators(boardID uint) ([]model.Moderator, error)
	IsModerator(userID, boardID uint) (bool, error)
	UpdateModeratorPermissions(_ context.Context, input UpdateModeratorPermissionsInput, operatorID uint) error
	CheckModeratorPermission(_ context.Context, userID, boardID uint, permission string) (bool, error)
	GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardWithPerms, error)
	// post
	DeletePost(boardID, postID, userID uint, isAdmin bool) error
	PinPost(boardID, postID uint, pin bool) error
}
type boardService struct {
	boardRepo boardRepo.BoardRepository
	userRepo  userRepo.UserRepository
	postRepo  postRepo.PostRepository
	notifSvc  notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
}

func NewBoardService(
	boardRepo boardRepo.BoardRepository,
	userRepo userRepo.UserRepository,
	postRepo postRepo.PostRepository,
	notifSvc notification.NotificationService,
) BoardService {
	return &boardService{
		boardRepo: boardRepo,
		userRepo:  userRepo,
		postRepo:  postRepo,
		notifSvc:  notifSvc,
	}
}
