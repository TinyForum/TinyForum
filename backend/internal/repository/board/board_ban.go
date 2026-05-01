package board

import (
	"time"
	"tiny-forum/internal/model/po"
)

// BanUser 禁言用户
func (r *boardRepository) BanUser(ban *po.BoardBan) error {
	return r.db.Create(ban).Error
}

// UnbanUser 解除禁言
func (r *boardRepository) UnbanUser(userID, boardID uint) error {
	return r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Delete(&po.BoardBan{}).Error
}

// IsBanned 检查用户是否被禁言
func (r *boardRepository) IsBanned(userID, boardID uint) (bool, error) {
	var count int64
	err := r.db.Model(&po.BoardBan{}).
		Where("user_id = ? AND board_id = ? AND (expires_at IS NULL OR expires_at > ?)",
			userID, boardID, time.Now()).
		Count(&count).Error
	return count > 0, err
}

// GetBan 获取禁言记录
func (r *boardRepository) GetBan(userID, boardID uint) (*po.BoardBan, error) {
	var ban po.BoardBan
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).First(&ban).Error
	return &ban, err
}
