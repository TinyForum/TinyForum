package post

import (
	"tiny-forum/internal/model"
	"tiny-forum/pkg/logger"
)

func (r *postRepository) GetQuestions(limit, offset int) ([]model.Post, int64, error) {
	logger.Info("[repository] GetQuestions")
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).
		Where("type = ? AND status = ?", "question", model.PostStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Preload("Tags").
		Preload("Board").
		Preload("Question").
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) GetUnansweredQuestions(limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Table("posts").
		Select("posts.*").
		Joins("LEFT JOIN questions ON posts.id = questions.post_id").
		Where("posts.type = ? AND posts.status = ? AND questions.accepted_answer_id IS NULL",
			"question", model.PostStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Author").
		Preload("Tags").
		Preload("Board").
		Preload("Question").
		Order("posts.created_at DESC").
		Find(&posts).Error

	return posts, total, err
}
