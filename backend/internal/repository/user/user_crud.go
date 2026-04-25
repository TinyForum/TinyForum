package user

import (
	"context"
	"tiny-forum/internal/model"
)

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByResetToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("reset_password_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *userRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.Model(&model.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *userRepository) FindByIDs(ids []uint) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	var users []model.User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (r *userRepository) GetUserBasicInfo(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Select("id, username, avatar").First(&user, id).Error
	return &user, err
}

func (r *userRepository) GetUserBasicInfoById(userID uint) (*model.User, error) {
	var user model.User
	err := r.db.Model(&model.User{}).
		Select("id, username, role").
		Where("id = ?", userID).
		First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserRoleById(userID uint) (string, error) {
	var role string
	err := r.db.Model(&model.User{}).
		Select("role").
		Where("id = ?", userID).
		Scan(&role).Error
	if err != nil {
		return "", err
	}
	return role, nil
}

// FindByEmailUnscoped 查找用户（包括已软删除的）
func (r *userRepository) FindByEmailUnscoped(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Unscoped().Where("email = ?", email).First(&user).Error
	return &user, err
}
