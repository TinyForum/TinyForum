package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type BoardRepository struct {
	db    *gorm.DB
	stats *StatsRepository
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{db: db,
		stats: NewStatsRepository(db)}
}

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

func (r *BoardRepository) Update(board *model.Board) error {
	return r.db.Save(board).Error
}

func (r *BoardRepository) Delete(id uint) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&model.Board{})
	return result.RowsAffected, result.Error
}

func (r *BoardRepository) FindByID(id uint) (*model.Board, error) {
	var board model.Board
	err := r.db.Preload("Parent").First(&board, id).Error
	return &board, err
}

func (r *BoardRepository) FindBySlug(slug string) (*model.Board, error) {
	var board model.Board
	err := r.db.Where("slug = ?", slug).First(&board).Error
	return &board, err
}

func (r *BoardRepository) GetPostsBySlug(slug string, page, pageSize int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	if slug == "" {
		return nil, 0, fmt.Errorf("slug cannot be empty")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := r.db.Model(&model.Post{}).
		Joins("JOIN boards ON boards.id = posts.board_id").
		Where("boards.slug = ?", slug).
		Where("posts.deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count posts failed: %w", err)
	}

	offset := (page - 1) * pageSize
	err := query.
		Offset(offset).
		Limit(pageSize).
		Order("posts.created_at DESC").
		Find(&posts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("get posts failed: %w", err)
	}
	return posts, total, nil
}

func (r *BoardRepository) List(limit, offset int) ([]model.Board, int64, error) {
	var boards []model.Board
	var total int64

	query := r.db.Model(&model.Board{})
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&boards).Error
	return boards, total, err
}

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

func (r *BoardRepository) IncrementPostCount(boardID uint, delta int) error {
	return r.db.Model(&model.Board{}).Where("id = ?", boardID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}

func (r *BoardRepository) IncrementThreadCount(boardID uint) error {
	return r.db.Model(&model.Board{}).Where("id = ?", boardID).
		UpdateColumn("thread_count", gorm.Expr("thread_count + 1")).Error
}

func (r *BoardRepository) IncrementTodayCount(boardID uint) error {
	return r.db.Model(&model.Board{}).Where("id = ?", boardID).
		UpdateColumn("today_count", gorm.Expr("today_count + 1")).Error
}

// ── Moderator ────────────────────────────────────────────────────────────────

// AddModerator 写入新版主记录（含权限 JSON）。
func (r *BoardRepository) AddModerator(mod *model.Moderator) error {
	return r.db.Create(mod).Error
}

// UpdateModerator 保存版主记录（用于权限更新）。
func (r *BoardRepository) UpdateModerator(mod *model.Moderator) error {
	return r.db.Save(mod).Error
}

// RemoveModerator 按 userID + boardID 软删除版主记录。
func (r *BoardRepository) RemoveModerator(userID, boardID uint) error {
	return r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Delete(&model.Moderator{}).Error
}

// FindModeratorByUserAndBoard 查询版主记录（含 User 预加载）。
func (r *BoardRepository) FindModeratorByUserAndBoard(userID, boardID uint) (*model.Moderator, error) {
	var mod model.Moderator
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Preload("User").
		First(&mod).Error
	if err != nil {
		return nil, err
	}
	return &mod, nil
}

// GetModerators 获取板块全部版主（含 User 预加载）。
func (r *BoardRepository) GetModerators(boardID uint) ([]model.Moderator, error) {
	var mods []model.Moderator
	err := r.db.Where("board_id = ?", boardID).
		Preload("User").
		Order("id ASC").
		Find(&mods).Error
	return mods, err
}

// IsModerator 检查用户是否已是版主。
func (r *BoardRepository) IsModerator(userID, boardID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Moderator{}).
		Where("user_id = ? AND board_id = ?", userID, boardID).
		Count(&count).Error
	return count > 0, err
}

// ── Ban ──────────────────────────────────────────────────────────────────────

func (r *BoardRepository) BanUser(ban *model.BoardBan) error {
	return r.db.Create(ban).Error
}

func (r *BoardRepository) UnbanUser(userID, boardID uint) error {
	return r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Delete(&model.BoardBan{}).Error
}

func (r *BoardRepository) IsBanned(userID, boardID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.BoardBan{}).
		Where("user_id = ? AND board_id = ? AND (expires_at IS NULL OR expires_at > ?)",
			userID, boardID, time.Now()).
		Count(&count).Error
	return count > 0, err
}

func (r *BoardRepository) GetBan(userID, boardID uint) (*model.BoardBan, error) {
	var ban model.BoardBan
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).First(&ban).Error
	return &ban, err
}

// ── ModeratorLog ─────────────────────────────────────────────────────────────

func (r *BoardRepository) CreateModeratorLog(log *model.ModeratorLog) error {
	return r.db.Create(log).Error
}

func (r *BoardRepository) GetModeratorLogs(boardID uint, limit, offset int) ([]model.ModeratorLog, int64, error) {
	var logs []model.ModeratorLog
	var total int64

	query := r.db.Model(&model.ModeratorLog{}).Where("board_id = ?", boardID)
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).
		Preload("Moderator").
		Order("created_at DESC").
		Find(&logs).Error
	return logs, total, err
}

// ── ModeratorApplication ─────────────────────────────────────────────────────

// CreateApplication 创建版主申请。
func (r *BoardRepository) CreateApplication(app *model.ModeratorApplication) error {
	return r.db.Create(app).Error
}

// FindPendingApplication 查找用户在某板块的待审核申请（用于去重）。
func (r *BoardRepository) FindPendingApplication(userID, boardID uint) (*model.ModeratorApplication, error) {
	var app model.ModeratorApplication
	err := r.db.Where("user_id = ? AND board_id = ? AND status = ?",
		userID, boardID, model.ApplicationPending).First(&app).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetApplicationByID 按 ID 获取申请（含关联预加载）。
func (r *BoardRepository) GetApplicationByID(id uint) (*model.ModeratorApplication, error) {
	var app model.ModeratorApplication
	err := r.db.Preload("User").Preload("Board").Preload("Reviewer").
		First(&app, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetApplicationsByUserID 根据用户ID获取申请记录（分页）
func (r *BoardRepository) GetApplicationsByUserID(userID uint, page, pageSize int) ([]model.ModeratorApplication, int64, error) {
	var applications []model.ModeratorApplication
	var total int64

	query := r.db.Model(&model.ModeratorApplication{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	err := query.Preload("Board").Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&applications).Error

	return applications, total, err
}

func (r *BoardRepository) GetLatestApplicationByUserAndBoard(userID, boardID uint) (*model.ModeratorApplication, error) {
	var app model.ModeratorApplication
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Order("created_at DESC").
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &app, nil
}

// UpdateApplication 保存申请（用于状态流转：审批/撤销）。
func (r *BoardRepository) UpdateApplication(app *model.ModeratorApplication) error {
	return r.db.Save(app).Error
}

// ListApplications 分页列出申请，可按板块和状态过滤。
// boardID == nil 则不过滤板块；status == "" 则不过滤状态。
func (r *BoardRepository) ListApplications(
	boardID *uint,
	status model.ApplicationStatus,
	page, pageSize int,
) ([]model.ModeratorApplication, int64, error) {
	var apps []model.ModeratorApplication
	var total int64

	query := r.db.Model(&model.ModeratorApplication{})
	if boardID != nil {
		query = query.Where("board_id = ?", *boardID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("User").Preload("Board").Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&apps).Error
	return apps, total, err
}

// CancelUserApplications 撤销用户在某板块所有 pending 申请（当用户被直接任命为版主时调用）。
func (r *BoardRepository) CancelUserApplications(userID, boardID uint) error {
	return r.db.Model(&model.ModeratorApplication{}).
		Where("user_id = ? AND board_id = ? AND status = ?", userID, boardID, model.ApplicationPending).
		Update("status", model.ApplicationCanceled).Error
}

// ── Stats helpers ─────────────────────────────────────────────────────────────

func (r *BoardRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Board{}).Count(&count).Error
	return count, err
}

func (r *BoardRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Board{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

func (r *BoardRepository) GetHotBoardsByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*HotBoardRow, error) {
	return r.stats.GetHotBoardsByDateRange(ctx, startDate, endDate, limit)
}

// repository/board_repository.go

// ModeratorBoardInfo 临时结构体，用于接收带权限的板块数据
type ModeratorBoardInfo struct {
	model.Board
	Permissions string `gorm:"column:permissions" json:"permissions"`
}

// GetModeratorBoardsWithPermissions 获取用户管理的板块及权限
func (r *BoardRepository) GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardInfo, error) {
	var results []ModeratorBoardInfo

	// 使用原始 SQL 查询，直接从 moderators 表获取 permissions 字段
	err := r.db.Raw(`
        SELECT 
            b.id,
            b.name,
            b.slug,
            b.description,
            b.icon,
            b.cover,
            b.parent_id,
            b.sort_order,
            b.view_role,
            b.post_role,
            b.reply_role,
            b.post_count,
            b.thread_count,
            b.today_count,
            b.created_at,
            b.updated_at,
            b.deleted_at,
            m.permissions
        FROM boards b
        INNER JOIN moderators m ON m.board_id = b.id
        WHERE m.user_id = ? 
        AND b.deleted_at IS NULL
        ORDER BY b.sort_order ASC, b.id ASC
    `, userID).Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("获取管理的板块失败: %w", err)
	}

	return results, nil
}
