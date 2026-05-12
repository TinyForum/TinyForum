package post

import (
	"context"
	"fmt"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
)

// MARK: list
// 获取文章列表
// 查询 published + approved
func (r *postRepository) List(ctx context.Context, ListPostsDO *common.PageQuery[do.Post]) ([]do.Post, int64, error) {
	var posts []do.Post
	var total int64

	// 构建基础查询
	baseQuery := r.db.Model(&do.Post{})

	// 用户感知状态过滤
	if ListPostsDO.Data.PostStatus != "" {
		baseQuery = baseQuery.Where("post_status = ?", ListPostsDO.Data.PostStatus)
	} else {
		// 默认只查询已发布的
		baseQuery = baseQuery.Where("post_status = ?", do.PostStatusPublished)
	}

	// 风控状态过滤
	if ListPostsDO.Data.ModerationStatus != "" {
		baseQuery = baseQuery.Where("moderation_status = ?", ListPostsDO.Data.ModerationStatus)
	} else {
		// 默认只查询已审核通过的
		baseQuery = baseQuery.Where("moderation_status = ?", do.ModerationStatusApproved)
	}

	if ListPostsDO.Data.AuthorID > 0 {
		baseQuery = baseQuery.Where("author_id = ?", ListPostsDO.Data.AuthorID)
	}

	// 标签过滤：修复为通过标签名称匹配（因为 TagNames 是 []string 标签名）
	if len(ListPostsDO.TagNames) > 0 {
		baseQuery = baseQuery.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name IN ?", ListPostsDO.TagNames).
			Distinct() // 避免因多个标签导致同一帖子重复
	}

	if ListPostsDO.Data.Type != "" {
		baseQuery = baseQuery.Where("type = ?", ListPostsDO.Data.Type)
	}

	if ListPostsDO.Keyword != "" {
		baseQuery = baseQuery.Where("title LIKE ? OR content LIKE ?", "%"+ListPostsDO.Keyword+"%", "%"+ListPostsDO.Keyword+"%")
	}

	// 统计总数：使用 Session 克隆，避免影响后续 Find
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数
	offset := (ListPostsDO.Page - 1) * ListPostsDO.PageSize
	orderExpr := "pin_top DESC, created_at DESC"
	if ListPostsDO.SortBy == "hot" {
		orderExpr = "pin_top DESC, like_count DESC, view_count DESC"
	}

	// 执行查询
	err := baseQuery.Preload("Author").Preload("Tags").
		Order(orderExpr).
		Offset(offset).Limit(ListPostsDO.PageSize).
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

// AdminList(ctx context.Context, ListPostsDO *common.PageQuery[do.Post]) ([]do.Post, int64, error)
func (r *postRepository) AdminList(ctx context.Context, listPostsDO *common.PageQuery[do.Post]) ([]do.Post, int64, error) {
	var posts []do.Post
	var total int64

	// 基础查询
	query := r.db.Model(&do.Post{})

	// 动态添加过滤条件
	// 状态过滤（后台需支持所有状态，仅当传入时才过滤）
	if listPostsDO.Data.PostStatus != "" {
		logger.Infof("查询文章状态: %s", listPostsDO.Data.PostStatus)
		query = query.Where("post_status = ?", listPostsDO.Data.PostStatus)
	}

	// 类型过滤
	if listPostsDO.Data.Type != "" {
		logger.Infof("查询类型: %s", listPostsDO.Data.Type)
		query = query.Where("type = ?", listPostsDO.Data.Type)
	}

	// 审核状态过滤
	if listPostsDO.Data.ModerationStatus != "" {
		logger.Infof("查询审核状态: %s", listPostsDO.Data.ModerationStatus)
		query = query.Where("moderation_status = ?", listPostsDO.Data.ModerationStatus)
	}

	// 关键词搜索（标题或内容）
	if listPostsDO.Keyword != "" {
		logger.Infof("查询关键字: %s", listPostsDO.Keyword)
		pattern := "%" + listPostsDO.Keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", pattern, pattern)
	}

	// 可选：作者ID过滤（注释保留，按需开启）
	// if listPostsDO.Data.AuthorID > 0 {
	//     logger.Infof("查询作者ID: %d", listPostsDO.Data.AuthorID)
	//     query = query.Where("author_id = ?", listPostsDO.Data.AuthorID)
	// }

	// 可选：标签过滤（注意多表关联可能影响性能）
	// if len(listPostsDO.Data.Tags) > 0 {
	//     query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
	//         Where("post_tags.tag_id IN ?", listPostsDO.Data.Tags)
	// }

	// 计算总数（在 offset/limit 之前）
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count posts failed: %w", err)
	}

	// 无数据时直接返回，避免不必要的查询
	if total == 0 {
		return []do.Post{}, 0, nil
	}

	// 分页参数
	offset := (listPostsDO.Page - 1) * listPostsDO.PageSize

	// 排序策略
	orderExpr := "created_at DESC"
	switch listPostsDO.SortBy {
	case "hot":
		orderExpr = "like_count DESC, view_count DESC, created_at DESC"
	case "latest":
		orderExpr = "created_at DESC"
	}

	// 执行查询（预加载关联数据）
	err := query.
		Preload("Author").
		Preload("Tags").
		Order(orderExpr).
		Offset(offset).
		Limit(listPostsDO.PageSize).
		Find(&posts).Error

	if err != nil {
		return nil, 0, fmt.Errorf("query posts failed: %w", err)
	}

	return posts, total, nil
}
