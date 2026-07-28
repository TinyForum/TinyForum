package question

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	"tiny-forum/pkg/logger"
)

// FindSimple 获取问题精简列表（旧版，保留兼容）
func (r *questionRepository) FindSimple(pageSize, offset int, boardID *uint) ([]vo.QuestionListResponse, int64, error) {
	var questions []vo.QuestionListResponse
	var total int64
	logger.Info("[Repository] FindSimple")

	query := r.db.Table("questions").
		Select(`
			questions.id,
			questions.created_at,
			questions.updated_at,
			questions.deleted_at,
			creations.title,
			creations.summary,
			creations.board_id,
			creations.author_id,
			questions.reward_score,
			questions.answer_count
		`).
		Joins("LEFT JOIN creations ON creations.id = questions.creation_id").
		Where("creations.deleted_at IS NULL").
		Where("creations.creation_status = ?", "published")

	if boardID != nil && *boardID > 0 {
		query = query.Where("creations.board_id = ?", *boardID)
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
func (r *questionRepository) FindSimpleQuestions(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]vo.QuestionSimpleDataVO, int64, error) {
	var questions []vo.QuestionSimpleDataVO
	var total int64
	logger.Info("[Repository] FindSimpleQuestions")

	// 使用 Model 代替 Table，并预先构建基础查询
	db := r.db.Model(&do.Question{}).
		Joins("LEFT JOIN creations ON creations.id = questions.creation_id").
		Where("creations.deleted_at IS NULL").
		Where("creations.creation_status = ?", "published")

	// 动态筛选：版块 ID
	if boardID != nil && *boardID > 0 {
		db = db.Where("creations.board_id = ?", *boardID)
	}

	// 动态筛选：关键词搜索
	if keyword != "" {
		db = db.Where("creations.title LIKE ? OR creations.summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
        questions.creation_id,
        questions.reward_score,
        questions.answer_count,
        questions.accepted_answer_id,
        questions.created_at,
        questions.updated_at,
        creations.title,
        creations.summary,
        creations.view_count,
        creations.board_id,
        creations.author_id
    `)

	// 动态排序
	switch sort {
	case "hot":
		db = db.Order("creations.view_count DESC, questions.answer_count DESC, questions.created_at DESC")
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
func (r *questionRepository) FindQuestionSimpleByID(questionID uint) (*vo.QuestionSimpleDataVO, error) {
	var question vo.QuestionSimpleDataVO
	err := r.db.Table("questions").
		Select(`
			questions.id,
			questions.creation_id,
			questions.reward_score,
			questions.answer_count,
			questions.accepted_answer_id,
			questions.created_at,
			questions.updated_at,
			creations.title,
			creations.summary,
			creations.content,
			creations.view_count,
			creations.board_id,
			creations.author_id
		`).
		Joins("LEFT JOIN creations ON creations.id = questions.creation_id").
		Where("questions.id = ?", questionID).
		Where("creations.deleted_at IS NULL").
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}
