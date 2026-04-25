package post

import (
	"context"
	"time"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	statsRepo "tiny-forum/internal/repository/stats"

	"gorm.io/gorm"
)

// PostRepository 帖子数据访问接口
type PostRepository interface {
	// 基础 CRUD
	Create(post *model.Post) error
	FindByID(id uint) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id uint) error
	List(page, pageSize int, opts dto.PostListOptions) ([]model.Post, int64, error)

	// 互动
	IncrViewCount(id uint) error
	IncrLikeCount(id uint, delta int) error
	AddLike(userID, postID uint) error
	RemoveLike(userID, postID uint) error
	IsLiked(userID, postID uint) bool

	// 管理
	AdminList(page, pageSize int, opts dto.PostListOptions) ([]model.Post, int64, error)

	// 板块相关
	GetByBoardID(boardID uint, limit, offset int) ([]model.Post, int64, error)

	// 问答相关
	GetQuestions(limit, offset int) ([]model.Post, int64, error)
	GetUnansweredQuestions(limit, offset int) ([]model.Post, int64, error)

	// 置顶
	TogglePinInBoard(postID uint, pin bool) error

	// 事务
	CreateWithTx(tx *gorm.DB, post *model.Post) error
	AddTags(tx *gorm.DB, post *model.Post, tagIDs []uint) error

	// 统计
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	GetHotArticlesByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]*statsRepo.HotArticleRow, error)
}

type postRepository struct {
	db    *gorm.DB
	stats statsRepo.StatsRepository // 假设 StatsRepository 已在同一包或需要导入
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db:    db,
		stats: statsRepo.NewStatsRepository(db),
	}
}
