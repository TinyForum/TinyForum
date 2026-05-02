package post

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
)

func (r *postRepository) Create(post *do.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*do.Post, error) {
	var post do.Post
	err := r.db.Preload("Author").Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post *do.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&do.Post{}, id).Error
}

// MARK: list
// 获取文章列表
// 查询 published + approved
func (r *postRepository) List(page, pageSize int, opts dto.PostListOptions) ([]do.Post, int64, error) {
	var posts []do.Post
	var total int64

	query := r.db.Model(&do.Post{})

	// 用户感知状态过滤
	if opts.Status != "" {
		query = query.Where("post_status = ?", opts.Status)
	} else {
		// 默认只查询已发布的
		query = query.Where("post_status = ?", do.PostStatusPublished)
	}

	// 风控状态过滤
	if opts.ModerationStatus != "" {
		query = query.Where("moderation_status = ?", opts.ModerationStatus)
	} else {
		// 默认只查询已审核通过的
		query = query.Where("moderation_status = ?", do.ModerationStatusApproved)
	}
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

func (r *postRepository) ListUserPosts(ctx context.Context, req request.GetUserPostsRequest, userID uint, orderBy string) ([]do.Post, int64, error) {
	db := r.db.WithContext(ctx).Model(&do.Post{}).Where("author_id = ?", userID)

	// 状态过滤
	if req.Status != "" {
		db = db.Where("post_status = ?", req.Status)
	}
	if req.ModerationStatus != "" {
		db = db.Where("moderation_status = ?", req.ModerationStatus)
	}
	// 标签过滤（JOIN）
	if req.Tag != "" {
		db = db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name = ?", req.Tag)
	}
	// 板块过滤
	if req.BoardName != "" {
		db = db.Joins("JOIN boards ON boards.id = posts.board_id").
			Where("boards.name = ?", req.BoardName)
	}
	// 关键词搜索
	if req.Keyword != "" {
		pattern := "%" + req.Keyword + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", pattern, pattern)
	}

	// 总数统计
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []do.Post{}, 0, nil
	}

	// 直接使用传入的排序表达式（已由 Service 层保证安全）
	var posts []do.Post
	err := db.Preload("Tags").Preload("Board").
		Order(orderBy).
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&posts).Error
	return posts, total, err
}

func (r *postRepository) AdminList(page, pageSize int, opts dto.PostListOptions) ([]do.Post, int64, error) {
	var posts []do.Post
	var total int64

	query := r.db.Model(&do.Post{})

	// 状态过滤：不设默认，只有传入时才过滤（后台需要看所有状态）
	if opts.Status != "" {
		query = query.Where("status = ?", opts.Status)
	}

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
	if opts.ModerationStatus != "" {
		query = query.Where("moderation_status = ?", opts.ModerationStatus)
	}

	// 统计
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	// 排序：默认按创建时间倒序，支持 hot 和 latest（latest 与默认相同）
	orderExpr := "created_at DESC"
	switch opts.SortBy {
	case "hot":
		orderExpr = "like_count DESC, view_count DESC"
	case "latest":
		orderExpr = "created_at DESC"
	}
	// 如果后台需要全局置顶，可以加上：orderExpr = "pin_top DESC, " + orderExpr

	err := query.Preload("Author").Preload("Tags").
		Order(orderExpr).
		Offset(offset).Limit(pageSize).
		Find(&posts).Error
	return posts, total, err
}
