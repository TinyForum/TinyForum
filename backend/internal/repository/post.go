package repository

import (
	"tiny-forum/internal/model"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("Author").Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err // 直接返回 nil 和错误
	}
	return &post, nil
}

func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

func (r *PostRepository) List(page, pageSize int, opts PostListOptions) ([]model.Post, int64, error) {
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

func (r *PostRepository) IncrViewCount(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *PostRepository) IncrLikeCount(id uint, delta int) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

// Like
func (r *PostRepository) AddLike(userID, postID uint) error {
	like := &model.Like{UserID: userID, PostID: &postID}
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		FirstOrCreate(like).Error
}

func (r *PostRepository) RemoveLike(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&model.Like{}).Error
}

func (r *PostRepository) IsLiked(userID, postID uint) bool {
	var count int64
	r.db.Model(&model.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// Admin: all posts including draft/hidden
func (r *PostRepository) AdminList(page, pageSize int, keyword string) ([]model.Post, int64, error) {
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

func (r *PostRepository) GetByBoardID(boardID uint, limit, offset int) ([]model.Post, int64, error) {
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
func (r *PostRepository) GetQuestions(limit, offset int) ([]model.Post, int64, error) {
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

func (r *PostRepository) GetUnansweredQuestions(limit, offset int) ([]model.Post, int64, error) {
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

func (r *PostRepository) TogglePinInBoard(postID uint, pin bool) error {
	return r.db.Model(&model.Post{}).Where("id = ?", postID).
		Update("pin_in_board", pin).Error
}

// CreateWithTx 使用事务创建帖子
func (r *PostRepository) CreateWithTx(tx *gorm.DB, post *model.Post) error {
	return tx.Create(post).Error
}

// AddTags 添加标签关联
func (r *PostRepository) AddTags(tx *gorm.DB, post *model.Post, tagIDs []uint) error {
	if len(tagIDs) == 0 {
		return nil
	}

	var tags []model.Tag
	if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}

	return tx.Model(post).Association("Tags").Append(&tags)
}
