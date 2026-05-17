package post

import (
	"context"
	"time"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	statsRepo "tiny-forum/internal/repository/stats"

	"gorm.io/gorm"
)

// PostRepository 帖子数据访问接口
type PostRepository interface {
	// 基础 CRUD
	Create(post *do.Post) error
	FindByID(id uint) (*do.Post, error)
	Update(post *do.Post) error
	Delete(id uint) error
	// List(page, pageSize int, opts bo.ListPosts) ([]do.Post, int64, error)
	List(ctx context.Context, ListPostsDO *common.PageQuery[do.Post]) ([]do.Post, int64, error)                                // 列表
	ListUserPosts(ctx context.Context, req request.GetUserPostsRequest, userID uint, orderBy string) ([]do.Post, int64, error) //  用户帖子列表

	// 互动
	IncrViewCount(id uint) error            // 增加浏览量
	IncrLikeCount(id uint, delta int) error // 增加点赞量
	AddLike(userID, postID uint) error      // 点赞
	RemoveLike(userID, postID uint) error   // 取消点赞
	IsLiked(userID, postID uint) bool       // 是否已点赞

	// 管理
	// AdminList(page, pageSize int, opts bo.ListPosts) ([]do.Post, int64, error)
	AdminList(ctx context.Context, ListPostsDO *common.PageQuery[do.Post]) ([]do.Post, int64, error) // 管理员列表

	// 板块相关
	GetByBoardID(boardID uint, limit, offset int) ([]do.Post, int64, error) // 根据板块ID获取帖子

	// 问答相关
	GetQuestions(limit, offset int) ([]do.Post, int64, error)           // 获取问题列表
	GetUnansweredQuestions(limit, offset int) ([]do.Post, int64, error) // 获取未回答的问题列表

	// 置顶
	TogglePinInBoard(postID uint, pin bool) error // 在板块内置顶/取消置顶

	// 事务
	CreateWithTx(tx *gorm.DB, post *do.Post) error           // 创建帖子
	AddTags(tx *gorm.DB, post *do.Post, tagIDs []uint) error // 添加标签

	// 统计
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)                                     // 按日期范围统计
	GetHotArticlesByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]*vo.HotArticleRowVO, error) // 按日期范围获取热门文章
}

type postRepository struct {
	db    *gorm.DB
	stats statsRepo.StatsRepository
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db:    db,
		stats: statsRepo.NewStatsRepository(db),
	}
}
