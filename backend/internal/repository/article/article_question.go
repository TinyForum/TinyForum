package article

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/logger"
)

func (r *articleRepository) GetQuestions(limit, offset int) ([]do.Question, int64, error) {
	logger.Info("[repository] GetQuestions")
	var posts []do.Question
	var total int64

	query := r.db.Model(&do.Question{}).
		Joins(`LEFT JOIN "creations" ON "creations"."id" = "questions"."creation_id" AND "creations"."deleted_at" IS NULL`).
		Where("creations.type = ? AND creations.creation_status = ?", "question", do.CreationStatusPublished)

	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Creation.Author").
		Preload("Creation.Tags").
		Preload("Creation.Board").
		Preload("Creation").
		Preload("Creation.Question").
		Order("creations.created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *articleRepository) GetUnansweredQuestions(limit, offset int) ([]do.Question, int64, error) {
	var posts []do.Question
	var total int64

	// 使用 Model 自动映射表名，Where 条件优先使用结构体
	db := r.db.Model(&do.Question{}).
		Joins(`LEFT JOIN "creations" ON "creations"."id" = "articles"."creation_id" AND "creations"."deleted_at" IS NULL`).
		Joins("LEFT JOIN questions ON questions.post_id = articles.id").
		Where("creations.type = ? AND creations.creation_status = ?", do.CreationTypeQuestion, do.CreationStatusPublished).
		Where("questions.accepted_answer_id IS NULL")

	// 统计总数（错误处理）
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询并预加载关联
	err := db.
		Offset(offset).
		Limit(limit).
		Preload("Creation.Author").
		Preload("Creation.Tags").
		Preload("Creation.Board").
		Preload("Creation").
		Preload("Creation.Question").
		Order("creations.created_at DESC").
		Find(&posts).Error

	return posts, total, err
}
