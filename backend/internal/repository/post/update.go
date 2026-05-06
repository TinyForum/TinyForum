package post

import "tiny-forum/internal/model/do"

func (r *postRepository) Create(post *do.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*do.Post, error) {
	var post do.Post
	err := r.db.Preload("Author").Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post *do.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&do.Post{}, id).Error
}
