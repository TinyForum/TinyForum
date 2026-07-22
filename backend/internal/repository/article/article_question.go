package article

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/logger"
)

func (r *articleRepository) GetQuestions(limit, offset int) ([]do.Article, int64, error) {
	logger.Info("[repository] GetQuestions")
	var posts []do.Article
	var total int64

	query := r.db.Model(&do.Article{}).
		Where("type = ? AND post_status = ?", "question", do.PostStatusPublished)

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

func (r *articleRepository) GetUnansweredQuestions(limit, offset int) ([]do.Article, int64, error) {
	var posts []do.Article
	var total int64

	// 使用 Model 自动映射表名，Where 条件优先使用结构体
	db := r.db.Model(&do.Article{}).
		Joins("LEFT JOIN questions ON posts.id = questions.post_id"). // JOIN 仍需要原生 SQL
		Where(&do.Article{Type: "question", PostStatus: do.PostStatusPublished}).
		Where("questions.accepted_answer_id IS NULL") // IS NULL 条件也保留为原生片段

	// 统计总数（错误处理）
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询并预加载关联
	err := db.
		Offset(offset).
		Limit(limit).
		Preload("Author").
		Preload("Tags").
		Preload("Board").
		Preload("Question").
		Order("posts.created_at DESC").
		Find(&posts).Error

	return posts, total, err
}
