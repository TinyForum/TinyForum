package tag

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *tagRepository) Create(tag *do.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) FindByID(id uint) (*do.Tag, error) {
	var tag do.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *tagRepository) FindByName(name string) (*do.Tag, error) {
	var tag do.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

func (r *tagRepository) List() ([]do.Tag, error) {
	var tags []do.Tag
	err := r.db.Order("post_count DESC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) Update(tag *do.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&do.Tag{}, id).Error
}

func (r *tagRepository) IncrPostCount(id uint, delta int) error {
	return r.db.Model(&do.Tag{}).Where("id = ?", id).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}
