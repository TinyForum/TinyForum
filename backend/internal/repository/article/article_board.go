package article

import (
	"context"

	"tiny-forum/internal/model/do"
)

func (r *articleRepository) GetByBoardID(boardID uint, limit, offset int) ([]do.Article, int64, error) {
	var posts []do.Article
	var total int64

	query := r.db.Model(&do.Article{}).
		Where("board_id = ? AND status = ?", boardID, do.CreationStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Creation.Author").
		Preload("Creation.Tags").
		Preload("Creation.Board").
		Preload("Creation").
		Order("pin_top DESC, pin_in_board DESC, created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *articleRepository) TogglePinInBoard(postID uint, pin bool) error {
	return r.db.Model(&do.Article{}).Where("id = ?", postID).
		Update("pin_in_board", pin).Error
}

// ListPostsByBoardID 分页获取指定板块下的帖子（版主面板，不做已发布过滤），可按关键词/状态过滤
func (r *articleRepository) ListPostsByBoardID(ctx context.Context, boardID uint, keyword string, status do.CreationStatus, limit, offset int) ([]do.Article, int64, error) {
	base := r.db.WithContext(ctx).Model(&do.Article{}).
		Joins(`LEFT JOIN "creations" ON "creations"."id" = "articles"."creation_id" AND "creations"."deleted_at" IS NULL`).
		Where("creations.board_id = ?", boardID)
	if keyword != "" {
		like := "%" + keyword + "%"
		base = base.Where(`("creations"."title" ILIKE ? OR "creations"."content" ILIKE ?)`, like, like)
	}
	if status != "" {
		base = base.Where("creations.creation_status = ?", status)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var posts []do.Article
	err := base.Select("articles.*").
		Order("creations.pin_in_board DESC, creations.created_at DESC").
		Offset(offset).Limit(limit).
		Preload("Creation").
		Preload("Creation.Author").
		Preload("Creation.Board").
		Preload("Creation.Tags").
		Find(&posts).Error
	return posts, total, err
}

// CountCommentsByPostID 批量统计指定帖子 ID 下的评论数
func (r *articleRepository) CountCommentsByPostID(ctx context.Context, postIDs []uint) (map[uint]int64, error) {
	result := make(map[uint]int64, len(postIDs))
	if len(postIDs) == 0 {
		return result, nil
	}
	type row struct {
		WorksID uint
		Cnt     int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&do.Comment{}).
		Select("works_id, COUNT(*) AS cnt").
		Where("works_id IN ?", postIDs).
		Group("works_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, rw := range rows {
		result[rw.WorksID] = rw.Cnt
	}
	return result, nil
}
