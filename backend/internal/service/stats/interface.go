package stats

import (
	"context"
	"time"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
	boardRepo "tiny-forum/internal/repository/board"
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	statsRepo "tiny-forum/internal/repository/stats"
	tagRepo "tiny-forum/internal/repository/tag"
	userRepo "tiny-forum/internal/repository/user"
)

type StatsService interface {
	GetStatsByDate(ctx context.Context, date time.Time, statsType string) (*po.StatsTodayInfo, error)
	GetTotalStats(ctx context.Context, startDate, endDate time.Time, statsType string) (*po.StatsInfoResp, error)
	GetTrendStats(ctx context.Context, startDate, endDate time.Time, statsType, intervals string) ([]*po.TrendData, error)
	GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time, statsType string) ([]dto.DailyStatResponse, error)
}

type statsService struct {
	statsRepo   statsRepo.StatsRepository
	postRepo    postRepo.PostRepository
	tagRepo     tagRepo.TagRepository
	boardRepo   boardRepo.BoardRepository
	userRepo    userRepo.UserRepository
	commentRepo commentRepo.CommentRepository
}

func NewStatsService(
	statsRepo statsRepo.StatsRepository,
	postRepo postRepo.PostRepository,
	tagRepo tagRepo.TagRepository,
	boardRepo boardRepo.BoardRepository,
	userRepo userRepo.UserRepository,
	commentRepo commentRepo.CommentRepository,
) StatsService {
	return &statsService{
		statsRepo:   statsRepo,
		postRepo:    postRepo,
		tagRepo:     tagRepo,
		boardRepo:   boardRepo,
		userRepo:    userRepo,
		commentRepo: commentRepo,
	}
}
