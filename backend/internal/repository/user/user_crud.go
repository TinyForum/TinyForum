package user

import (
	"tiny-forum/internal/model"
)

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

func (r *UserRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
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

func (r *UserRepository) FindByIDs(ids []uint) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	var users []model.User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUserBasicInfo(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Select("id, username, avatar").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetUserBasicInfoById(userID uint) (*model.User, error) {
	var user model.User
	err := r.db.Model(&model.User{}).
		Select("id, username, role").
		Where("id = ?", userID).
		First(&user).Error
	return &user, err
}

func (r *UserRepository) GetUserRoleById(userID uint) (string, error) {
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
