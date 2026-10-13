package board

import (
	"context"
	"errors"
	"time"

	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// BanUser 禁言用户
func (r *boardRepository) BanUser(ban *do.BoardBan) error {
	return r.db.Create(ban).Error
}

// UnbanUser 解除禁言
func (r *boardRepository) UnbanUser(userID, boardID uint) error {
	return r.db.Where("user_id = ? AND board_id = ?", userID, boardID).
		Delete(&do.BoardBan{}).Error
}

// IsBanned 检查用户是否被禁言
func (r *boardRepository) IsBanned(userID, boardID uint) (bool, error) {
	var count int64
	err := r.db.Model(&do.BoardBan{}).
		Where("user_id = ? AND board_id = ? AND (expires_at IS NULL OR expires_at > ?)",
			userID, boardID, time.Now()).
		Count(&count).Error
	return count > 0, err
}

// GetBan 获取禁言记录
func (r *boardRepository) GetBan(userID, boardID uint) (*do.BoardBan, error) {
	var ban do.BoardBan
	err := r.db.Where("user_id = ? AND board_id = ?", userID, boardID).First(&ban).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ban, nil
}

// ListBans 分页获取板块封禁记录（预加载被禁用户与封禁操作者）
func (r *boardRepository) ListBans(ctx context.Context, boardID uint, limit, offset int) ([]do.BoardBan, int64, error) {
	base := r.db.WithContext(ctx).Model(&do.BoardBan{}).Where("board_id = ?", boardID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var bans []do.BoardBan
	err := base.Preload("User").Preload("Board").Preload("Banner").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&bans).Error
	return bans, total, err
}
