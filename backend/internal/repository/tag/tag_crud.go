package tag

import (
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

func (r *tagRepository) Create(tag *po.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) FindByID(id uint) (*po.Tag, error) {
	var tag po.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *tagRepository) FindByName(name string) (*po.Tag, error) {
	var tag po.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

func (r *tagRepository) List() ([]po.Tag, error) {
	var tags []po.Tag
	err := r.db.Order("post_count DESC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) Update(tag *po.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&po.Tag{}, id).Error
}

func (r *tagRepository) IncrPostCount(id uint, delta int) error {
	return r.db.Model(&po.Tag{}).Where("id = ?", id).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", delta)).Error
}
