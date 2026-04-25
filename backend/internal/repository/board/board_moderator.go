package board

import (
	"tiny-forum/internal/model"
)

// AddModerator 添加版主
func (r *boardRepository) AddModerator(mod *model.Moderator) error {
	return r.db.Create(mod).Error
}

// UpdateModerator 更新版主权限
func (r *boardRepository) UpdateModerator(mod *model.Moderator) error {
	return r.db.Save(mod).Error
}

// RemoveModerator 移除版主
func (r *boardRepository) RemoveModerator(userID, boardID uint) error {
	return r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Delete(&model.Moderator{}).Error
}

// FindModeratorByUserAndBoard 根据用户和板块查询版主记录
func (r *boardRepository) FindModeratorByUserAndBoard(userID, boardID uint) (*model.Moderator, error) {
	var mod model.Moderator
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Preload("User").
		First(&mod).Error
	if err != nil {
		return nil, err
	}
	return &mod, nil
}

// GetModerators 获取板块的所有版主
func (r *boardRepository) GetModerators(boardID uint) ([]model.Moderator, error) {
	var mods []model.Moderator
	err := r.db.Where("board_id = ?", boardID).
		Preload("User").
		Order("id ASC").
		Find(&mods).Error
	return mods, err
}

// IsModerator 判断用户是否为版主
func (r *boardRepository) IsModerator(userID, boardID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Moderator{}).
		Where("user_id = ? AND board_id = ?", userID, boardID).
		Count(&count).Error
	return count > 0, err
}

// CreateModeratorLog 创建版主操作日志
func (r *boardRepository) CreateModeratorLog(log *model.ModeratorLog) error {
	return r.db.Create(log).Error
}

// GetModeratorLogs 获取版主操作日志
func (r *boardRepository) GetModeratorLogs(boardID uint, limit, offset int) ([]model.ModeratorLog, int64, error) {
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

// ModeratorBoardInfo 用户管理的板块信息（含权限）
type ModeratorBoardInfo struct {
	model.Board
	Permissions string `gorm:"column:permissions" json:"permissions"`
}

// GetModeratorBoardsWithPermissions 获取用户管理的板块及权限
func (r *boardRepository) GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardInfo, error) {
	var results []ModeratorBoardInfo

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
		return nil, err
	}
	return results, nil
}
