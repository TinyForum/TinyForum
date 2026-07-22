package article

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *articleRepository) CreateWithTx(tx *gorm.DB, post *do.Article) error {
	return tx.Create(post).Error
}

func (r *articleRepository) AddTags(tx *gorm.DB, post *do.Article, tagIDs []uint) error {
	if len(tagIDs) == 0 {
		return nil
	}
	var tags []do.Tag
	if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}
	return tx.Model(post).Association("Tags").Append(&tags)
}
