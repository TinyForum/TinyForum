package board

import (
	"context"
	"time"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository/stats"
	statsRepo "tiny-forum/internal/repository/stats"

	"gorm.io/gorm"
)

type BoardRepository interface {
	// apply
	CreateApplication(app *model.ModeratorApplication) error
	FindPendingApplication(userID, boardID uint) (*model.ModeratorApplication, error)
	GetApplicationByID(id uint) (*model.ModeratorApplication, error)
	GetApplicationsByUserID(userID uint, page, pageSize int) ([]model.ModeratorApplication, int64, error)
	GetLatestApplicationByUserAndBoard(userID, boardID uint) (*model.ModeratorApplication, error)
	UpdateApplication(app *model.ModeratorApplication) error
	ListApplications(
		boardID *uint,
		status model.ApplicationStatus,
		page, pageSize int,
	) ([]model.ModeratorApplication, int64, error)
	CancelUserApplications(userID, boardID uint) error
	// ban
	BanUser(ban *model.BoardBan) error
	UnbanUser(userID, boardID uint) error
	IsBanned(userID, boardID uint) (bool, error)
	GetBan(userID, boardID uint) (*model.BoardBan, error)
	// moderator
	AddModerator(mod *model.Moderator) error
	UpdateModerator(mod *model.Moderator) error
	RemoveModerator(userID, boardID uint) error
	FindModeratorByUserAndBoard(userID, boardID uint) (*model.Moderator, error)
	GetModerators(boardID uint) ([]model.Moderator, error)
	IsModerator(userID, boardID uint) (bool, error)
	CreateModeratorLog(log *model.ModeratorLog) error
	GetModeratorLogs(boardID uint, limit, offset int) ([]model.ModeratorLog, int64, error)
	GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardInfo, error)
	// repo
	Create(board *model.Board) error
	Update(board *model.Board) error
	Delete(id uint) (int64, error)
	FindByID(id uint) (*model.Board, error)
	FindBySlug(slug string) (*model.Board, error)
	GetPostsBySlug(slug string, page, pageSize int) ([]*dto.GetBoardPostsResponse, int64, error)
	List(limit, offset int) ([]model.Board, int64, error)
	GetTree() ([]model.Board, error)
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
