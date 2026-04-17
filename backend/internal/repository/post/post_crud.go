package post

import (
	"tiny-forum/internal/model"
)

func (r *postRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("Author").Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err
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
