package article

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
type ArticleRepository interface {
	// 基础 CRUD
	Create(post *do.Article) error
	FindByArticleID(id uint) (*do.Article, error)           // 根据帖子ID获取帖子
	FindQuestionByQuestionID(id uint) (*do.Question, error) // 根据问题ID获取问题
	Update(post *do.Article) error
	Delete(id uint) error
	// List(page, pageSize int, opts bo.ListPosts) ([]do.Post, int64, error)
	List(ctx context.Context, ListPostsDO *common.PageQuery[do.Article]) ([]do.Article, int64, error)                             // 列表
	ListUserPosts(ctx context.Context, req request.GetUserPostsRequest, userID uint, orderBy string) ([]do.Article, int64, error) //  用户帖子列表

	// 互动
	IncrViewCount(id uint) error            // 增加浏览量
	IncrLikeCount(id uint, delta int) error // 增加点赞量
	AddLike(userID, postID uint) error      // 点赞
	RemoveLike(userID, postID uint) error   // 取消点赞
	IsLiked(userID, postID uint) bool       // 是否已点赞

	// 管理
	// AdminList(page, pageSize int, opts bo.ListPosts) ([]do.Post, int64, error)
	AdminList(ctx context.Context, ListPostsDO *common.PageQuery[do.Article]) ([]do.Article, int64, error) // 管理员列表

	// 板块相关
	GetByBoardID(boardID uint, limit, offset int) ([]do.Article, int64, error) // 根据板块ID获取帖子
	// ListPostsByBoardID 分页获取指定板块下的帖子（版主面板，不做已发布过滤），可按关键词/状态过滤
	ListPostsByBoardID(ctx context.Context, boardID uint, keyword string, status do.CreationStatus, limit, offset int) ([]do.Article, int64, error)
	// CountCommentsByPostID 批量统计指定帖子 ID 下的评论数
	CountCommentsByPostID(ctx context.Context, postIDs []uint) (map[uint]int64, error)

	// 问答相关
	GetQuestions(limit, offset int) ([]do.Question, int64, error)           // 获取问题列表
	GetUnansweredQuestions(limit, offset int) ([]do.Question, int64, error) // 获取未回答的问题列表

	// 置顶
	TogglePinInBoard(postID uint, pin bool) error // 在板块内置顶/取消置顶

	// 事务
	CreateWithTx(tx *gorm.DB, post *do.Article) error           // 创建帖子
	AddTags(tx *gorm.DB, post *do.Article, tagIDs []uint) error // 添加标签

	// 统计
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)                                     // 按日期范围统计
	GetHotArticlesByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]*vo.HotArticleRowVO, error) // 按日期范围获取热门文章
}

type articleRepository struct {
	db    *gorm.DB
	stats statsRepo.StatsRepository
}

func NewPostRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{
		db:    db,
		stats: statsRepo.NewStatsRepository(db),
	}
}
