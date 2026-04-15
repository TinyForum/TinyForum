package repository

import (
	"context"
	"time"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *model.Post) error
	FindByID(id uint) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id uint) error
	List(page, pageSize int, opts PostListOptions) ([]model.Post, int64, error)

	IncrViewCount(id uint) error
	IncrLikeCount(id uint, delta int) error
	AddLike(userID, postID uint) error
	RemoveLike(userID, postID uint) error
	IsLiked(userID, postID uint) bool

	AdminList(page, pageSize int, keyword string) ([]model.Post, int64, error)
	GetByBoardID(boardID uint, limit, offset int) ([]model.Post, int64, error)
	GetQuestions(limit, offset int) ([]model.Post, int64, error)
	GetUnansweredQuestions(limit, offset int) ([]model.Post, int64, error)
	TogglePinInBoard(postID uint, pin bool) error

	CreateWithTx(tx *gorm.DB, post *model.Post) error
	AddTags(tx *gorm.DB, post *model.Post, tagIDs []uint) error
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	GetHotArticlesByDateRange(
		ctx context.Context,
		startDate, endDate time.Time,
		limit int,
	) ([]*HotArticleRow, error)
}

type postRepository struct {
	db    *gorm.DB
	stats *StatsRepository
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db,
		stats: NewStatsRepository(db)}
}

func (r *postRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("Author").Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err // 直接返回 nil 和错误
	}
	return &post, nil
}

func (r *postRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

func (r *postRepository) List(page, pageSize int, opts PostListOptions) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("status = ?", model.PostStatusPublished)

	if opts.AuthorID > 0 {
		query = query.Where("author_id = ?", opts.AuthorID)
	}
	if opts.TagID > 0 {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id = ?", opts.TagID)
	}
	if opts.PostType != "" {
		query = query.Where("type = ?", opts.PostType)
	}
	if opts.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+opts.Keyword+"%", "%"+opts.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	orderExpr := "pin_top DESC, created_at DESC"
	if opts.SortBy == "hot" {
		orderExpr = "pin_top DESC, like_count DESC, view_count DESC"
	}

	err := query.Preload("Author").Preload("Tags").
		Order(orderExpr).
		Offset(offset).Limit(pageSize).
		Find(&posts).Error
	return posts, total, err
}

type PostListOptions struct {
	AuthorID uint
	TagID    uint
	PostType string
	Keyword  string
	SortBy   string // "" = latest, "hot" = popular
}

func (r *postRepository) IncrViewCount(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *postRepository) IncrLikeCount(id uint, delta int) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

// Like
func (r *postRepository) AddLike(userID, postID uint) error {
	like := &model.Like{UserID: userID, PostID: &postID}
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		FirstOrCreate(like).Error
}

// RemoveLike 是一个用于移除用户对帖子的点赞的方法
// 它接收两个参数：用户ID(userID)和帖子ID(postID)，都是无符号整数类型
// 该方法会删除数据库中对应的点赞记录，并返回可能的错误
func (r *postRepository) RemoveLike(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&model.Like{}).Error
}

func (r *postRepository) IsLiked(userID, postID uint) bool {
	var count int64
	r.db.Model(&model.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// Admin: all posts including draft/hidden
func (r *postRepository) AdminList(page, pageSize int, keyword string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	query := r.db.Model(&model.Post{})
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Preload("Author").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}

// 在现有 PostRepository 中添加以下方法

func (r *postRepository) GetByBoardID(boardID uint, limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).
		Where("board_id = ? AND status = ?", boardID, model.PostStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Preload("Tags").
		Preload("Board").
		Order("pin_top DESC, pin_in_board DESC, created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

// 获取问题
func (r *postRepository) GetQuestions(limit, offset int) ([]model.Post, int64, error) {
	logger.Info("[repository] GetQuestions")
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).
		Where("type = ? AND status = ?", "question", model.PostStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Preload("Tags").
		Preload("Board").
		Preload("Question").
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) GetUnansweredQuestions(limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Table("posts").
		Select("posts.*").
		Joins("LEFT JOIN questions ON posts.id = questions.post_id").
		Where("posts.type = ? AND posts.status = ? AND questions.accepted_answer_id IS NULL",
			"question", model.PostStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Preload("Tags").
		Preload("Board").
		Preload("Question").
		Order("posts.created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) TogglePinInBoard(postID uint, pin bool) error {
	return r.db.Model(&model.Post{}).Where("id = ?", postID).
		Update("pin_in_board", pin).Error
}

// CreateWithTx 使用事务创建帖子
func (r *postRepository) CreateWithTx(tx *gorm.DB, post *model.Post) error {
	return tx.Create(post).Error
}

// AddTags 添加标签关联
func (r *postRepository) AddTags(tx *gorm.DB, post *model.Post, tagIDs []uint) error {
	if len(tagIDs) == 0 {
		return nil
	}

	var tags []model.Tag
	if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}

	return tx.Model(post).Association("Tags").Append(&tags)
}

// Count 返回帖子总数
func (r *postRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Post{}).Count(&count).Error
	return count, err
}

// CountByDateRange 统计指定时间段内新增帖子数
func (r *postRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Post{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetHotArticlesByDateRange 委托给 StatsRepository 执行复合查询
func (r *postRepository) GetHotArticlesByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*HotArticleRow, error) {
	return r.stats.GetHotArticlesByDateRange(ctx, startDate, endDate, limit)
}
