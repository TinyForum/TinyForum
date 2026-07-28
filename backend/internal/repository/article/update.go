package article

import "tiny-forum/internal/model/do"

func (r *articleRepository) Create(post *do.Article) error {
	if err := r.db.Create(&post.Creation).Error; err != nil {
		return err
	}
	post.CreationID = post.Creation.ID
	return r.db.Create(post).Error
}

func (r *articleRepository) Update(post *do.Article) error {
	return r.db.Save(post).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&do.Article{}, id).Error
}
