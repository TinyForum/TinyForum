package board

import (
	"errors"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

// CreateApplication 创建版主申请
func (r *boardRepository) CreateApplication(app *po.ModeratorApplication) error {
	return r.db.Create(app).Error
}

// FindPendingApplication 查找用户在某板块的待审核申请
func (r *boardRepository) FindPendingApplication(userID, boardID uint) (*po.ModeratorApplication, error) {
	var app po.ModeratorApplication
	err := r.db.Where("user_id = ? AND board_id = ? AND status = ?",
		userID, boardID, po.ApplicationPending).First(&app).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetApplicationByID 根据 ID 获取申请
func (r *boardRepository) GetApplicationByID(id uint) (*po.ModeratorApplication, error) {
	var app po.ModeratorApplication
	err := r.db.Preload("User").Preload("Board").Preload("Reviewer").
		First(&app, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetApplicationsByUserID 获取用户的所有申请记录
func (r *boardRepository) GetApplicationsByUserID(userID uint, page, pageSize int) ([]po.ModeratorApplication, int64, error) {
	var applications []po.ModeratorApplication
	var total int64

	query := r.db.Model(&po.ModeratorApplication{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Board").Preload("Reviewer").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&applications).Error

	return applications, total, err
}

// GetLatestApplicationByUserAndBoard 获取用户在某板块的最新申请
func (r *boardRepository) GetLatestApplicationByUserAndBoard(userID, boardID uint) (*po.ModeratorApplication, error) {
	var app po.ModeratorApplication
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

// UpdateApplication 更新申请状态
func (r *boardRepository) UpdateApplication(app *po.ModeratorApplication) error {
	return r.db.Save(app).Error
}

// ListApplications 分页列出申请（可按板块和状态过滤）
func (r *boardRepository) ListApplications(
	boardID *uint,
	status po.ApplicationStatus,
	page, pageSize int,
) ([]po.ModeratorApplication, int64, error) {
	var apps []po.ModeratorApplication
	var total int64

	query := r.db.Model(&po.ModeratorApplication{})
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

// CancelUserApplications 撤销用户在某板块的所有 pending 申请
func (r *boardRepository) CancelUserApplications(userID, boardID uint) error {
	return r.db.Model(&po.ModeratorApplication{}).
		Where("user_id = ? AND board_id = ? AND status = ?", userID, boardID, po.ApplicationPending).
		Update("status", po.ApplicationCanceled).Error
}
