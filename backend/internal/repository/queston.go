package repository

import (
	"errors"
	"time"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

// func (r *QuestionRepository) Create(question *model.Question) error {
// 	return r.db.Create(question).Error
// }

func (r *QuestionRepository) Update(question *model.Question) error {
	return r.db.Save(question).Error
}

func (r *QuestionRepository) FindByPostID(postID uint) (*model.Question, error) {
	var question model.Question
	err := r.db.Where("post_id = ?", postID).
		Preload("Post").
		Preload("AcceptedAnswer").
		First(&question).Error
	return &question, err
}

func (r *QuestionRepository) IncrementAnswerCount(postID uint) error {
	return r.db.Model(&model.Question{}).Where("post_id = ?", postID).
		UpdateColumn("answer_count", gorm.Expr("answer_count + 1")).Error
}

func (r *QuestionRepository) SetAcceptedAnswer(postID, commentID uint) error {
	return r.db.Model(&model.Question{}).Where("post_id = ?", postID).
		Updates(map[string]interface{}{
			"accepted_answer_id": commentID,
		}).Error
}

// AnswerVote methods

func (r *QuestionRepository) CreateAnswerVote(vote *model.AnswerVote) error {
	return r.db.Create(vote).Error
}

// Bug Fix #2: 新增 UpdateAnswerVote，用 Save 更新已存在的投票记录，避免触发唯一索引冲突
func (r *QuestionRepository) UpdateAnswerVote(vote *model.AnswerVote) error {
	return r.db.Save(vote).Error
}

func (r *QuestionRepository) DeleteAnswerVote(userID, commentID uint) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&model.AnswerVote{}).Error
}

func (r *QuestionRepository) FindAnswerVote(userID, commentID uint) (*model.AnswerVote, error) {
	var vote model.AnswerVote
	err := r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).
		First(&vote).Error
	return &vote, err
}

func (r *QuestionRepository) GetAnswerVoteCount(commentID uint) (int, error) {
	var upCount, downCount int64
	r.db.Model(&model.AnswerVote{}).Where("comment_id = ? AND vote_type = ?", commentID, "up").Count(&upCount)
	r.db.Model(&model.AnswerVote{}).Where("comment_id = ? AND vote_type = ?", commentID, "down").Count(&downCount)
	return int(upCount - downCount), nil
}

func (r *QuestionRepository) UpdateCommentVoteCount(commentID uint, voteCount int) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", commentID).
		UpdateColumn("vote_count", voteCount).Error
}

// UpdateAnswerCount 更新回答数量
func (r *QuestionRepository) UpdateAnswerCount(questionID uint) error {
	return r.db.Model(&model.Question{}).
		Where("id = ?", questionID).
		Update("answer_count", gorm.Expr("answer_count + ?", 1)).Error
}

// UpdateAcceptedAnswer 更新采纳的答案
func (r *QuestionRepository) UpdateAcceptedAnswer(questionID uint, answerID uint) error {
	return r.db.Model(&model.Question{}).
		Where("id = ?", questionID).
		Update("accepted_answer_id", answerID).Error
}

// CreateWithTransaction 使用事务创建问答（包括帖子、标签、积分扣减）
func (r *QuestionRepository) CreateWithTransaction(userID uint, input model.CreateQuestionInput) (*model.QuestionResponse, error) {
	// 开启事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建帖子
	post := &model.Post{
		Title:    input.Title,
		Content:  input.Content,
		Summary:  input.Summary,
		Cover:    input.Cover,
		BoardID:  input.BoardID,
		AuthorID: userID,
		Type:     "question",
		Status:   "published",
	}

	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. 关联标签
	if len(input.TagIDs) > 0 {
		var tags []model.Tag
		if err := tx.Where("id IN ?", input.TagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(post).Association("Tags").Append(&tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 3. 扣减积分（如果有悬赏）
	if input.RewardScore > 0 {
		var user model.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if user.Score < input.RewardScore {
			tx.Rollback()
			return nil, errors.New("积分不足")
		}

		if err := tx.Model(&user).Update("score", gorm.Expr("score - ?", input.RewardScore)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 4. 创建问答记录
	question := &model.Question{
		PostID:      post.ID,
		RewardScore: input.RewardScore,
		AnswerCount: 0,
	}

	if err := tx.Create(question).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 返回响应
	return &model.QuestionResponse{
		ID:          question.ID,
		PostID:      post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Summary:     post.Summary,
		Cover:       post.Cover,
		BoardID:     post.BoardID,
		AuthorID:    post.AuthorID,
		RewardScore: question.RewardScore,
		AnswerCount: question.AnswerCount,
		Status:      string(post.Status),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}, nil
}

// Create 创建问答记录
func (r *QuestionRepository) Create(question *model.Question) error {
	return r.db.Create(question).Error
}

// FindByID 根据ID查询问答
func (r *QuestionRepository) FindByID(id uint) (*model.Question, error) {
	var question model.Question
	err := r.db.Preload("Post").Preload("Post.Tags").Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// FindSimple 获取问题精简列表
func (r *QuestionRepository) FindSimple(pageSize, offset int, boardID *uint) ([]model.QuestionListResponse, int64, error) {
	var questions []model.QuestionListResponse
	var total int64
	logger.Info("[Repository] FindSimple")

	// 构建查询
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
		Where("posts.status = ?", "published")

	// 按板块过滤
	if boardID != nil && *boardID > 0 {
		query = query.Where("posts.board_id = ?", *boardID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.
		Order("questions.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&questions).Error

	return questions, total, err
}

// FindSimpleQuestions 只查询问题基础数据，不加载关联
func (r *QuestionRepository) FindSimpleQuestions(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]QuestionSimpleData, int64, error) {
	var questions []QuestionSimpleData
	var total int64
	logger.Info("[Repository] FindSimpleQuestions")

	// 构建查询
	query := r.db.Table("questions").
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
			posts.view_count,
			posts.board_id,
			posts.author_id
		`).
		Joins("LEFT JOIN posts ON posts.id = questions.post_id").
		Where("posts.deleted_at IS NULL").
		Where("posts.status = ?", "published")

	// 按板块过滤
	if boardID != nil && *boardID > 0 {
		query = query.Where("posts.board_id = ?", *boardID)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("posts.title LIKE ? OR posts.summary LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 应用筛选
	switch filter {
	case "unanswered":
		query = query.Where("questions.answer_count = 0")
	case "answered":
		query = query.Where("questions.accepted_answer_id IS NOT NULL")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用排序
	switch sort {
	case "hot":
		query = query.Order("posts.view_count DESC, questions.answer_count DESC, questions.created_at DESC")
	case "score":
		query = query.Order("questions.reward_score DESC, questions.created_at DESC")
	default: // latest
		query = query.Order("questions.created_at DESC")
	}

	// 获取分页数据
	err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&questions).Error

	return questions, total, err
}

// FindQuestionSimpleByID 根据ID查询单个问题基础数据
func (r *QuestionRepository) FindQuestionSimpleByID(questionID uint) (*QuestionSimpleData, error) {
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

type QuestionSimpleData struct {
	ID               uint      `gorm:"column:id"`
	PostID           uint      `gorm:"column:post_id"`
	Title            string    `gorm:"column:title"`
	Summary          string    `gorm:"column:summary"`
	ViewCount        int       `gorm:"column:view_count"`
	BoardID          uint      `gorm:"column:board_id"`
	AuthorID         uint      `gorm:"column:author_id"`
	RewardScore      int       `gorm:"column:reward_score"`
	AnswerCount      int       `gorm:"column:answer_count"`
	AcceptedAnswerID *uint     `gorm:"column:accepted_answer_id"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}
