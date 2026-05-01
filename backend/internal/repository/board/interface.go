package board

import (
	"context"
	"time"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/repository/stats"
	statsRepo "tiny-forum/internal/repository/stats"

	"gorm.io/gorm"
)

type BoardRepository interface {
	// apply
	CreateApplication(app *do.ModeratorApplication) error
	FindPendingApplication(userID, boardID uint) (*do.ModeratorApplication, error)
	GetApplicationByID(id uint) (*do.ModeratorApplication, error)
	GetApplicationsByUserID(userID uint, page, pageSize int) ([]do.ModeratorApplication, int64, error)
	GetLatestApplicationByUserAndBoard(userID, boardID uint) (*do.ModeratorApplication, error)
	UpdateApplication(app *do.ModeratorApplication) error
	ListApplications(
		boardID *uint,
		status do.ApplicationStatus,
		page, pageSize int,
	) ([]do.ModeratorApplication, int64, error)
	CancelUserApplications(userID, boardID uint) error
	// ban
	BanUser(ban *do.BoardBan) error
	UnbanUser(userID, boardID uint) error
	IsBanned(userID, boardID uint) (bool, error)
	GetBan(userID, boardID uint) (*do.BoardBan, error)
	// moderator
	AddModerator(mod *do.Moderator) error
	UpdateModerator(mod *do.Moderator) error
	RemoveModerator(userID, boardID uint) error
	FindModeratorByUserAndBoard(userID, boardID uint) (*do.Moderator, error)
	GetModerators(boardID uint) ([]do.Moderator, error)
	IsModerator(userID, boardID uint) (bool, error)
	CreateModeratorLog(log *do.ModeratorLog) error
	GetModeratorLogs(boardID uint, limit, offset int) ([]do.ModeratorLog, int64, error)
	GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardInfo, error)
	// repo
	Create(board *do.Board) error
	Update(board *do.Board) error
	Delete(id uint) (int64, error)
	FindByID(id uint) (*do.Board, error)
	FindBySlug(slug string) (*do.Board, error)
	GetPostsBySlug(slug string, page, pageSize int) ([]*dto.GetBoardPostsResponse, int64, error)
	List(limit, offset int) ([]do.Board, int64, error)
	GetTree() ([]do.Board, error)
	IncrementPostCount(boardID uint, delta int) error
	IncrementThreadCount(boardID uint) error
	IncrementTodayCount(boardID uint) error
	// stats
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	GetHotBoardsByDateRange(
		ctx context.Context,
		startDate, endDate time.Time,
		limit int,
	) ([]*statsRepo.HotBoardRow, error)
}

type boardRepository struct {
	db    *gorm.DB
	stats stats.StatsRepository
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{
		db:    db,
		stats: stats.NewStatsRepository(db),
	}
}
