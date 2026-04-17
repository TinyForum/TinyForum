package board

import (
	"errors"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// CreateApplication 创建版主申请
func (r *BoardRepository) CreateApplication(app *model.ModeratorApplication) error {
	return r.db.Create(app).Error
}

// FindPendingApplication 查找用户在某板块的待审核申请
func (r *BoardRepository) FindPendingApplication(userID, boardID uint) (*model.ModeratorApplication, error) {
	var app model.ModeratorApplication
	err := r.db.Where("user_id = ? AND board_id = ? AND status = ?",
		userID, boardID, model.ApplicationPending).First(&app).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetApplicationByID 根据 ID 获取申请
func (r *BoardRepository) GetApplicationByID(id uint) (*model.ModeratorApplication, error) {
	var app model.ModeratorApplication
	err := r.db.Preload("User").Preload("Board").Preload("Reviewer").
		First(&app, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetApplicationsByUserID 获取用户的所有申请记录
func (r *BoardRepository) GetApplicationsByUserID(userID uint, page, pageSize int) ([]model.ModeratorApplication, int64, error) {
	var applications []model.ModeratorApplication
	var total int64

	query := r.db.Model(&model.ModeratorApplication{}).Where("user_id = ?", userID)

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

// UpdateApplication 更新申请状态
func (r *BoardRepository) UpdateApplication(app *model.ModeratorApplication) error {
	return r.db.Save(app).Error
}

// ListApplications 分页列出申请（可按板块和状态过滤）
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

// CancelUserApplications 撤销用户在某板块的所有 pending 申请
func (r *BoardRepository) CancelUserApplications(userID, boardID uint) error {
	return r.db.Model(&model.ModeratorApplication{}).
		Where("user_id = ? AND board_id = ? AND status = ?", userID, boardID, model.ApplicationPending).
		Update("status", model.ApplicationCanceled).Error
}
