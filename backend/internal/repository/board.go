package repository

import (
	"errors"
	"fmt"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type BoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(board *model.Board) error {
	// 确保 ParentID 的处理
	if board.ParentID != nil && *board.ParentID == 0 {
		board.ParentID = nil
	}
	// 使用事务创建
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 再次验证父板块存在（防止并发问题）
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

// Moderator methods
func (r *BoardRepository) AddModerator(mod *model.Moderator) error {
	return r.db.Create(mod).Error
}

func (r *BoardRepository) RemoveModerator(userID, boardID uint) error {
	return r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Delete(&model.Moderator{}).Error
}

func (r *BoardRepository) FindModeratorByUserAndBoard(userID, boardID uint) (*model.Moderator, error) {
	var mod model.Moderator
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Preload("User").
		First(&mod).Error
	return &mod, err
}

func (r *BoardRepository) GetModerators(boardID uint) ([]model.Moderator, error) {
	var mods []model.Moderator
	err := r.db.Where("board_id = ?", boardID).
		Preload("User").
		Order("id ASC").
		Find(&mods).Error
	return mods, err
}

func (r *BoardRepository) IsModerator(userID, boardID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Moderator{}).
		Where("user_id = ? AND board_id = ?", userID, boardID).
		Count(&count).Error
	return count > 0, err
}

// Ban methods
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
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		First(&ban).Error
	return &ban, err
}

// ModeratorLog methods
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
