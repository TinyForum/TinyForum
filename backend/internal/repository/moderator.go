// // internal/repository/moderator_repo.go
// package repository

// import (
// 	"context"
// 	"tiny-forum/internal/model"

// 	"encoding/json"
// 	"errors"

// 	"gorm.io/gorm"
// )

// type ModeratorRepository interface {
// 	// 基础 CRUD
// 	Create(ctx context.Context, moderator *model.Moderator) error
// 	Update(ctx context.Context, moderator *model.Moderator) error
// 	Delete(ctx context.Context, id uint) error
// 	GetByID(ctx context.Context, id uint) (*model.Moderator, error)

// 	// 查询
// 	GetByUserAndBoard(ctx context.Context, userID, boardID uint) (*model.Moderator, error)
// 	GetByBoard(ctx context.Context, boardID uint) ([]model.Moderator, error)
// 	GetByUser(ctx context.Context, userID uint) ([]model.Moderator, error)

// 	// 权限相关
// 	UpdatePermissions(ctx context.Context, moderatorID uint, permissions model.Permission) error
// 	HasPermission(ctx context.Context, userID, boardID uint, permission string) (bool, error)

// 	// 批量操作
// 	DeleteByBoard(ctx context.Context, boardID uint) error
// 	DeleteByUser(ctx context.Context, userID uint) error

// 	// 检查
// 	Exists(ctx context.Context, userID, boardID uint) (bool, error)
// 	IsModerator(ctx context.Context, userID, boardID uint) (bool, error)

// 	// 分页列表
// 	List(ctx context.Context, page, pageSize int, boardID *uint) ([]model.Moderator, int64, error)
// }

// type moderatorRepository struct {
// 	db *gorm.DB
// }

// func NewModeratorRepository(db *gorm.DB) ModeratorRepository {
// 	return &moderatorRepository{db: db}
// }

// // Create 创建版主
// func (r *moderatorRepository) Create(ctx context.Context, moderator *model.Moderator) error {
// 	return r.db.WithContext(ctx).Create(moderator).Error
// }

// // Update 更新版主
// func (r *moderatorRepository) Update(ctx context.Context, moderator *model.Moderator) error {
// 	return r.db.WithContext(ctx).Save(moderator).Error
// }

// // Delete 删除版主
// func (r *moderatorRepository) Delete(ctx context.Context, id uint) error {
// 	return r.db.WithContext(ctx).Delete(&model.Moderator{}, id).Error
// }

// // GetByID 根据ID获取版主
// func (r *moderatorRepository) GetByID(ctx context.Context, id uint) (*model.Moderator, error) {
// 	var moderator model.Moderator
// 	err := r.db.WithContext(ctx).
// 		Preload("User").
// 		Preload("Board").
// 		First(&moderator, id).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &moderator, nil
// }

// // GetByUserAndBoard 根据用户ID和板块ID获取版主
// func (r *moderatorRepository) GetByUserAndBoard(ctx context.Context, userID, boardID uint) (*model.Moderator, error) {
// 	var moderator model.Moderator
// 	err := r.db.WithContext(ctx).
// 		Where("user_id = ? AND board_id = ?", userID, boardID).
// 		First(&moderator).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &moderator, nil
// }

// // GetByBoard 获取板块下的所有版主
// func (r *moderatorRepository) GetByBoard(ctx context.Context, boardID uint) ([]model.Moderator, error) {
// 	var moderators []model.Moderator
// 	err := r.db.WithContext(ctx).
// 		Preload("User").
// 		Where("board_id = ?", boardID).
// 		Find(&moderators).Error
// 	return moderators, err
// }

// // GetByUser 获取用户管理的所有板块
// func (r *moderatorRepository) GetByUser(ctx context.Context, userID uint) ([]model.Moderator, error) {
// 	var moderators []model.Moderator
// 	err := r.db.WithContext(ctx).
// 		Preload("Board").
// 		Where("user_id = ?", userID).
// 		Find(&moderators).Error
// 	return moderators, err
// }

// // UpdatePermissions 更新版主权限
// func (r *moderatorRepository) UpdatePermissions(ctx context.Context, moderatorID uint, permissions model.Permission) error {
// 	permsJSON, err := json.Marshal(permissions)
// 	if err != nil {
// 		return err
// 	}

// 	return r.db.WithContext(ctx).
// 		Model(&model.Moderator{}).
// 		Where("id = ?", moderatorID).
// 		Update("permissions", permsJSON).Error
// }

// // HasPermission 检查版主是否有特定权限
// func (r *moderatorRepository) HasPermission(ctx context.Context, userID, boardID uint, permission string) (bool, error) {
// 	var moderator model.Moderator
// 	err := r.db.WithContext(ctx).
// 		Where("user_id = ? AND board_id = ?", userID, boardID).
// 		First(&moderator).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return false, nil
// 		}
// 		return false, err
// 	}

// 	perms, err := moderator.GetPermissions()
// 	if err != nil {
// 		return false, err
// 	}

// 	switch permission {
// 	case "delete_post":
// 		return perms.CanDeletePost, nil
// 	case "pin_post":
// 		return perms.CanPinPost, nil
// 	case "edit_any_post":
// 		return perms.CanEditAnyPost, nil
// 	case "manage_moderator":
// 		return perms.CanManageModerator, nil
// 	case "ban_user":
// 		return perms.CanBanUser, nil
// 	default:
// 		return false, nil
// 	}
// }

// // DeleteByBoard 删除板块下的所有版主
// func (r *moderatorRepository) DeleteByBoard(ctx context.Context, boardID uint) error {
// 	return r.db.WithContext(ctx).
// 		Where("board_id = ?", boardID).
// 		Delete(&model.Moderator{}).Error
// }

// // DeleteByUser 删除用户的所有版主记录
// func (r *moderatorRepository) DeleteByUser(ctx context.Context, userID uint) error {
// 	return r.db.WithContext(ctx).
// 		Where("user_id = ?", userID).
// 		Delete(&model.Moderator{}).Error
// }

// // Exists 检查版主是否存在
// func (r *moderatorRepository) Exists(ctx context.Context, userID, boardID uint) (bool, error) {
// 	var count int64
// 	err := r.db.WithContext(ctx).
// 		Model(&model.Moderator{}).
// 		Where("user_id = ? AND board_id = ?", userID, boardID).
// 		Count(&count).Error
// 	return count > 0, err
// }

// // IsModerator 检查用户是否是版主（别名方法）
// func (r *moderatorRepository) IsModerator(ctx context.Context, userID, boardID uint) (bool, error) {
// 	return r.Exists(ctx, userID, boardID)
// }

// // List 分页获取版主列表
// func (r *moderatorRepository) List(ctx context.Context, page, pageSize int, boardID *uint) ([]model.Moderator, int64, error) {
// 	var moderators []model.Moderator
// 	var total int64

// 	query := r.db.WithContext(ctx).Model(&model.Moderator{})

// 	if boardID != nil {
// 		query = query.Where("board_id = ?", *boardID)
// 	}

// 	// 获取总数
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	// 分页查询
// 	offset := (page - 1) * pageSize
// 	err := query.
// 		Preload("User").
// 		Preload("Board").
// 		Offset(offset).
// 		Limit(pageSize).
// 		Order("created_at DESC").
// 		Find(&moderators).Error

//		return moderators, total, err
//	}
package repository
