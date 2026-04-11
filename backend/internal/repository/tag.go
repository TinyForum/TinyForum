package repository

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *TagRepository) FindByName(name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

func (r *TagRepository) List() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("post_count DESC").Find(&tags).Error
	return tags, err
}

func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

func (r *TagRepository) IncrPostCount(id uint, delta int) error {
	return r.db.Model(&model.Tag{}).Where("id = ?", id).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}
