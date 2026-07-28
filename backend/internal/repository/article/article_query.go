package article

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
func (r *articleRepository) List(ctx context.Context, ListPostsDO *common.PageQuery[do.Article]) ([]do.Article, int64, error) {
	var posts []do.Article
	var total int64

	// 构建基础查询
	baseQuery := r.db.Model(&do.Article{}).
		Joins(`LEFT JOIN "creations" ON "creations"."id" = "articles"."creation_id" AND "creations"."deleted_at" IS NULL`)

	// 用户感知状态过滤
	if ListPostsDO.Data.Creation.CreationStatus != "" {
		baseQuery = baseQuery.Where("creations.creation_status = ?", ListPostsDO.Data.Creation.CreationStatus)
	} else {
		// 默认只查询已发布的
		baseQuery = baseQuery.Where("creations.creation_status = ?", do.CreationStatusPublished)
	}

	// 风控状态过滤
	if ListPostsDO.Data.Creation.ModerationStatus != "" {
		baseQuery = baseQuery.Where("creations.moderation_status = ?", ListPostsDO.Data.Creation.ModerationStatus)
	} else {
		// 默认只查询已审核通过的
		baseQuery = baseQuery.Where("creations.moderation_status = ?", do.ModerationStatusApproved)
	}

	if ListPostsDO.Data.Creation.AuthorID > 0 {
		baseQuery = baseQuery.Where("creations.author_id = ?", ListPostsDO.Data.Creation.AuthorID)
	}

	// 标签过滤：通过 creation_tags 关联
	if len(ListPostsDO.TagNames) > 0 {
		baseQuery = baseQuery.Joins("JOIN creation_tags ON creation_tags.creation_id = creations.id").
			Joins("JOIN tags ON tags.id = creation_tags.tag_id").
			Where("tags.name IN ?", ListPostsDO.TagNames).
			Distinct()
	}

	if ListPostsDO.Data.Creation.Type != "" {
		baseQuery = baseQuery.Where("creations.type = ?", ListPostsDO.Data.Creation.Type)
	}

	if ListPostsDO.Keyword != "" {
		baseQuery = baseQuery.Where("creations.title LIKE ? OR creations.content LIKE ?", "%"+ListPostsDO.Keyword+"%", "%"+ListPostsDO.Keyword+"%")
	}

	// 统计总数：使用 Session 克隆，避免影响后续 Find
	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数
	offset := (ListPostsDO.Page - 1) * ListPostsDO.PageSize
	orderExpr := "creations.pin_top DESC, creations.created_at DESC"
	if ListPostsDO.SortBy == "hot" {
		orderExpr = "creations.pin_top DESC, creations.like_count DESC, creations.view_count DESC"
	}

	// 执行查询
	err := baseQuery.Preload("Creation.Author").Preload("Creation.Tags").Preload("Creation").
		Order(orderExpr).
		Offset(offset).Limit(ListPostsDO.PageSize).
		Find(&posts).Error

	return posts, total, err
}
func (r *articleRepository) ListUserPosts(ctx context.Context, req request.GetUserPostsRequest, userID uint, orderBy string) ([]do.Article, int64, error) {
	db := r.db.WithContext(ctx).Model(&do.Article{}).
		Joins(`LEFT JOIN "creations" ON "creations"."id" = "articles"."creation_id" AND "creations"."deleted_at" IS NULL`).
		Where("creations.author_id = ?", userID)

	// 状态过滤
	if req.Status != "" {
		db = db.Where("creations.creation_status = ?", req.Status)
	}
	if req.ModerationStatus != "" {
		db = db.Where("creations.moderation_status = ?", req.ModerationStatus)
	}
	// 标签过滤（JOIN）
	if req.Tag != "" {
		db = db.Joins("JOIN creation_tags ON creation_tags.creation_id = creations.id").
			Joins("JOIN tags ON tags.id = creation_tags.tag_id").
			Where("tags.name = ?", req.Tag)
	}
	// 板块过滤
	if req.BoardName != "" {
		db = db.Joins("JOIN boards ON boards.id = creations.board_id").
			Where("boards.name = ?", req.BoardName)
	}
	// 关键词搜索
	if req.Keyword != "" {
		pattern := "%" + req.Keyword + "%"
		db = db.Where("creations.title LIKE ? OR creations.content LIKE ?", pattern, pattern)
	}

	// 总数统计
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []do.Article{}, 0, nil
	}

	// 直接使用传入的排序表达式（已由 Service 层保证安全）
	var posts []do.Article
	err := db.Preload("Creation.Tags").Preload("Creation.Board").Preload("Creation").
		Order(orderBy).
		Offset(req.Offset()).
		Limit(req.PageSize).
		Find(&posts).Error
	return posts, total, err
}

// AdminList(ctx context.Context, ListPostsDO *common.PageQuery[do.Post]) ([]do.Post, int64, error)
func (r *articleRepository) AdminList(ctx context.Context, listPostsDO *common.PageQuery[do.Article]) ([]do.Article, int64, error) {
	var posts []do.Article
	var total int64

	// 基础查询
	query := r.db.Model(&do.Article{}).
		Joins(`LEFT JOIN "creations" ON "creations"."id" = "articles"."creation_id" AND "creations"."deleted_at" IS NULL`)

	// 动态添加过滤条件
	// 状态过滤（后台需支持所有状态，仅当传入时才过滤）
	if listPostsDO.Data.Creation.CreationStatus != "" {
		logger.Infof("查询文章状态: %s", listPostsDO.Data.Creation.CreationStatus)
		query = query.Where("creations.creation_status = ?", listPostsDO.Data.Creation.CreationStatus)
	}

	// 类型过滤
	if listPostsDO.Data.Creation.Type != "" {
		logger.Infof("查询类型: %s", listPostsDO.Data.Creation.Type)
		query = query.Where("creations.type = ?", listPostsDO.Data.Creation.Type)
	}

	// 审核状态过滤
	if listPostsDO.Data.Creation.ModerationStatus != "" {
		logger.Infof("查询审核状态: %s", listPostsDO.Data.Creation.ModerationStatus)
		query = query.Where("creations.moderation_status = ?", listPostsDO.Data.Creation.ModerationStatus)
	}

	// 关键词搜索（标题或内容）
	if listPostsDO.Keyword != "" {
		logger.Infof("查询关键字: %s", listPostsDO.Keyword)
		pattern := "%" + listPostsDO.Keyword + "%"
		query = query.Where("creations.title LIKE ? OR creations.content LIKE ?", pattern, pattern)
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
		return []do.Article{}, 0, nil
	}

	// 分页参数
	offset := (listPostsDO.Page - 1) * listPostsDO.PageSize

	// 排序策略
	orderExpr := "creations.created_at DESC"
	switch listPostsDO.SortBy {
	case "hot":
		orderExpr = "creations.like_count DESC, creations.view_count DESC, creations.created_at DESC"
	case "latest":
		orderExpr = "creations.created_at DESC"
	}

	// 执行查询（预加载关联数据）
	err := query.
		Preload("Creation.Author").
		Preload("Creation.Tags").
		Preload("Creation").
		Order(orderExpr).
		Offset(offset).
		Limit(listPostsDO.PageSize).
		Find(&posts).Error

	if err != nil {
		return nil, 0, fmt.Errorf("query posts failed: %w", err)
	}

	return posts, total, nil
}

// 通过文章ID查找文章
func (r *articleRepository) FindByArticleID(id uint) (*do.Article, error) {
	var post do.Article
	err := r.db.Preload("Creation.Author").Preload("Creation.Tags").Preload("Creation").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *articleRepository) FindQuestionByQuestionID(id uint) (*do.Question, error) {
	var question do.Question

	// 构建查询，预加载所有可能用到的关联
	query := r.db.
		Preload("Creation.Author").
		Preload("Creation.Tags").
		Preload("Creation.Board"). // 如果业务需要 Board，可加上
		Preload("AcceptedAnswer")  // 确保 AcceptedAnswer 被加载

	// 如果希望包含软删除的关联（例如管理员查看），可取消注释：
	// query = query.Unscoped() // 注意：这会同时取消主表的软删除过滤，请按需使用

	err := query.First(&question, id).Error
	if err != nil {
		// 可封装为自定义错误，例如：
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		//     return nil, ErrQuestionNotFound
		// }
		return nil, err
	}
	return &question, nil
}
