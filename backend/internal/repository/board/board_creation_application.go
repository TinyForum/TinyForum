package board

import (
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// CreateBoardApplication 创建板块创建申请
func (r *boardRepository) CreateBoardApplication(app *do.BoardCreationApplication) error {
	return r.db.Create(app).Error
}

// FindPendingBoardApplicationBySlug 查找指定 slug 的待审核板块创建申请
func (r *boardRepository) FindPendingBoardApplicationBySlug(slug string) (*do.BoardCreationApplication, error) {
	var app do.BoardCreationApplication
	err := r.db.Where("slug = ? AND status = ?", slug, do.ApplicationPending).First(&app).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetBoardApplicationByID 根据 ID 获取板块创建申请
func (r *boardRepository) GetBoardApplicationByID(id uint) (*do.BoardCreationApplication, error) {
	var app do.BoardCreationApplication
	err := r.db.Preload("User").Preload("Reviewer").Preload("CreatedBoard").
		First(&app, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// GetBoardApplicationsByUserID 获取用户的所有板块创建申请记录
func (r *boardRepository) GetBoardApplicationsByUserID(userID uint, page, pageSize int) ([]do.BoardCreationApplication, int64, error) {
	var applications []do.BoardCreationApplication
	var total int64

	query := r.db.Model(&do.BoardCreationApplication{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Reviewer").Preload("CreatedBoard").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&applications).Error

	return applications, total, err
}

// UpdateBoardApplication 更新板块创建申请状态
func (r *boardRepository) UpdateBoardApplication(app *do.BoardCreationApplication) error {
	return r.db.Save(app).Error
}

// ListBoardApplications 管理员分页列出板块创建申请（可按状态过滤）
func (r *boardRepository) ListBoardApplications(
	status do.ApplicationStatus,
	page, pageSize int,
) ([]do.BoardCreationApplication, int64, error) {
	var apps []do.BoardCreationApplication
	var total int64

	query := r.db.Model(&do.BoardCreationApplication{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.
		Preload("User").Preload("Reviewer").Preload("CreatedBoard").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&apps).Error
	return apps, total, err
}
