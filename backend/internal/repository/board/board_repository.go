package board

import (
	"errors"
	"fmt"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository/stats"

	"gorm.io/gorm"
)

type BoardRepository struct {
	db    *gorm.DB
	stats *stats.StatsRepository
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{
		db:    db,
		stats: stats.NewStatsRepository(db),
	}
}

// Create 创建板块
func (r *BoardRepository) Create(board *model.Board) error {
	if board.ParentID != nil && *board.ParentID == 0 {
		board.ParentID = nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if board.ParentID != nil && *board.ParentID != 0 {
			var parent model.Board
			if err := tx.First(&parent, *board.ParentID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("父板块不存在: id=%d", *board.ParentID)
				}
				return err
			}
		}
		return tx.Create(board).Error
	})
}

// Update 更新板块
func (r *BoardRepository) Update(board *model.Board) error {
	return r.db.Save(board).Error
}

// Delete 删除板块
func (r *BoardRepository) Delete(id uint) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&model.Board{})
	return result.RowsAffected, result.Error
}

// FindByID 根据 ID 查找板块
func (r *BoardRepository) FindByID(id uint) (*model.Board, error) {
	var board model.Board
	err := r.db.Preload("Parent").First(&board, id).Error
	return &board, err
}

// FindBySlug 根据 slug 查找板块
func (r *BoardRepository) FindBySlug(slug string) (*model.Board, error) {
	var board model.Board
	err := r.db.Where("slug = ?", slug).First(&board).Error
	return &board, err
}

// GetPostsBySlug 根据板块 slug 获取帖子列表
// repository/board_repo.go
func (r *BoardRepository) GetPostsBySlug(slug string, page, pageSize int) ([]*dto.GetBoardPostsResponse, int64, error) {
	var total int64
	var results []*dto.GetBoardPostsResponse

	// 1. 构建基础查询
	baseQuery := r.db.Table("posts").
		Select("posts.id, posts.title, posts.summary, posts.cover, posts.type, posts.author_id, users.username as author_name, posts.created_at").
		Joins("JOIN boards ON boards.id = posts.board_id").
		Joins("JOIN users ON users.id = posts.author_id").
		Where("boards.slug = ?", slug).
		Where("posts.deleted_at IS NULL").
		Where("posts.post_status = ?", "published") // 只显示已发布的帖子（根据业务调整）

	// 2. 统计总数
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count posts failed: %w", err)
	}

	// 3. 分页查询
	offset := (page - 1) * pageSize
	err := baseQuery.
		Offset(offset).
		Limit(pageSize).
		Order("posts.created_at DESC").
		Scan(&results).Error
	if err != nil {
		return nil, 0, fmt.Errorf("get posts failed: %w", err)
	}

	return results, total, nil
}

// List 分页获取板块列表
func (r *BoardRepository) List(limit, offset int) ([]model.Board, int64, error) {
	var boards []model.Board
	var total int64

	query := r.db.Model(&model.Board{})
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&boards).Error
	return boards, total, err
}

// GetTree 获取板块树形结构
func (r *BoardRepository) GetTree() ([]model.Board, error) {
	var boards []model.Board
	err := r.db.Where("parent_id IS NULL").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		Order("sort_order ASC, id ASC").
		Find(&boards).Error
	return boards, err
}

// IncrementPostCount 增加板块帖子计数
func (r *BoardRepository) IncrementPostCount(boardID uint, delta int) error {
	return r.db.Model(&model.Board{}).Where("id = ?", boardID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}

// IncrementThreadCount 增加板块主题计数
func (r *BoardRepository) IncrementThreadCount(boardID uint) error {
	return r.db.Model(&model.Board{}).Where("id = ?", boardID).
		UpdateColumn("thread_count", gorm.Expr("thread_count + 1")).Error
}

// IncrementTodayCount 增加板块今日计数
func (r *BoardRepository) IncrementTodayCount(boardID uint) error {
	return r.db.Model(&model.Board{}).Where("id = ?", boardID).
		UpdateColumn("today_count", gorm.Expr("today_count + 1")).Error
}
