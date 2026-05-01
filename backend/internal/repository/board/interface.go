package board

import (
	"context"
	"time"

	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
	"tiny-forum/internal/repository/stats"
	statsRepo "tiny-forum/internal/repository/stats"

	"gorm.io/gorm"
)

type BoardRepository interface {
	// apply
	CreateApplication(app *po.ModeratorApplication) error
	FindPendingApplication(userID, boardID uint) (*po.ModeratorApplication, error)
	GetApplicationByID(id uint) (*po.ModeratorApplication, error)
	GetApplicationsByUserID(userID uint, page, pageSize int) ([]po.ModeratorApplication, int64, error)
	GetLatestApplicationByUserAndBoard(userID, boardID uint) (*po.ModeratorApplication, error)
	UpdateApplication(app *po.ModeratorApplication) error
	ListApplications(
		boardID *uint,
		status po.ApplicationStatus,
		page, pageSize int,
	) ([]po.ModeratorApplication, int64, error)
	CancelUserApplications(userID, boardID uint) error
	// ban
	BanUser(ban *po.BoardBan) error
	UnbanUser(userID, boardID uint) error
	IsBanned(userID, boardID uint) (bool, error)
	GetBan(userID, boardID uint) (*po.BoardBan, error)
	// moderator
	AddModerator(mod *po.Moderator) error
	UpdateModerator(mod *po.Moderator) error
	RemoveModerator(userID, boardID uint) error
	FindModeratorByUserAndBoard(userID, boardID uint) (*po.Moderator, error)
	GetModerators(boardID uint) ([]po.Moderator, error)
	IsModerator(userID, boardID uint) (bool, error)
	CreateModeratorLog(log *po.ModeratorLog) error
	GetModeratorLogs(boardID uint, limit, offset int) ([]po.ModeratorLog, int64, error)
	GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardInfo, error)
	// repo
	Create(board *po.Board) error
	Update(board *po.Board) error
	Delete(id uint) (int64, error)
	FindByID(id uint) (*po.Board, error)
	FindBySlug(slug string) (*po.Board, error)
	GetPostsBySlug(slug string, page, pageSize int) ([]*dto.GetBoardPostsResponse, int64, error)
	List(limit, offset int) ([]po.Board, int64, error)
	GetTree() ([]po.Board, error)
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
