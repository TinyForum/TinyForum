package repository

import (
	"bbs-forum/internal/model"

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
	return &post, err
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
