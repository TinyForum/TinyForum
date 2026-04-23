package question

import (
	"tiny-forum/internal/model"
	"tiny-forum/pkg/logger"
)

// FindSimple 获取问题精简列表（旧版，保留兼容）
func (r *questionRepository) FindSimple(pageSize, offset int, boardID *uint) ([]model.QuestionListResponse, int64, error) {
	var questions []model.QuestionListResponse
	var total int64
	logger.Info("[Repository] FindSimple")

	query := r.db.Table("questions").
		Select(`
			questions.id,
			questions.created_at,
			questions.updated_at,
			questions.deleted_at,
			posts.title,
			posts.summary,
			posts.board_id,
			posts.author_id,
			questions.reward_score,
			questions.answer_count
		`).
		Joins("LEFT JOIN posts ON posts.id = questions.post_id").
		Where("posts.deleted_at IS NULL").
		Where("posts.post_status = ?", "published")

	if boardID != nil && *boardID > 0 {
		query = query.Where("posts.board_id = ?", *boardID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("questions.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&questions).Error

	return questions, total, err
}

// FindSimpleQuestions 查询问题基础数据（支持过滤、排序、关键词）
func (r *questionRepository) FindSimpleQuestions(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]QuestionSimpleData, int64, error) {
	var questions []QuestionSimpleData
	var total int64
	logger.Info("[Repository] FindSimpleQuestions")

	// 使用 Model 代替 Table，并预先构建基础查询
	db := r.db.Model(&model.Question{}).
		Joins("LEFT JOIN posts ON posts.id = questions.post_id"). // JOIN 需保留原生 SQL
		Where("posts.deleted_at IS NULL").                        // 软删除条件（posts 表）
		Where("posts.post_status = ?", "published")               // 帖子状态条件

	// 动态筛选：版块 ID
	if boardID != nil && *boardID > 0 {
		db = db.Where("posts.board_id = ?", *boardID)
	}

	// 动态筛选：关键词搜索
	if keyword != "" {
		db = db.Where("posts.title LIKE ? OR posts.summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 动态筛选：回答状态
	switch filter {
	case "unanswered":
		db = db.Where("questions.answer_count = 0")
	case "answered":
		db = db.Where("questions.accepted_answer_id IS NOT NULL")
	}

	// 统计总数（错误处理）
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 选择跨表字段（保留原生 SQL，因为字段多且涉及两个表）
	db = db.Select(`
        questions.id,
        questions.post_id,
        questions.reward_score,
        questions.answer_count,
        questions.accepted_answer_id,
        questions.created_at,
        questions.updated_at,
        posts.title,
        posts.summary,
        posts.view_count,
        posts.board_id,
        posts.author_id
    `)

	// 动态排序
	switch sort {
	case "hot":
		db = db.Order("posts.view_count DESC, questions.answer_count DESC, questions.created_at DESC")
	case "score":
		db = db.Order("questions.reward_score DESC, questions.created_at DESC")
	default:
		db = db.Order("questions.created_at DESC")
	}

	// 分页查询
	err := db.Offset(offset).Limit(pageSize).Find(&questions).Error
	return questions, total, err
}

// FindQuestionSimpleByID 根据ID查询单个问题基础数据
func (r *questionRepository) FindQuestionSimpleByID(questionID uint) (*QuestionSimpleData, error) {
	var question QuestionSimpleData
	err := r.db.Table("questions").
		Select(`
			questions.id,
			questions.post_id,
			questions.reward_score,
			questions.answer_count,
			questions.accepted_answer_id,
			questions.created_at,
			questions.updated_at,
			posts.title,
			posts.summary,
			posts.content,
			posts.view_count,
			posts.board_id,
			posts.author_id
		`).
		Joins("LEFT JOIN posts ON posts.id = questions.post_id").
		Where("questions.id = ?", questionID).
		Where("posts.deleted_at IS NULL").
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}
