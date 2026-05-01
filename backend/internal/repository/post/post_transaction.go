package post

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *postRepository) CreateWithTx(tx *gorm.DB, post *do.Post) error {
	return tx.Create(post).Error
}

func (r *postRepository) AddTags(tx *gorm.DB, post *do.Post, tagIDs []uint) error {
	if len(tagIDs) == 0 {
		return nil
	}
	var tags []do.Tag
	if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}
	return tx.Model(post).Association("Tags").Append(&tags)
}
