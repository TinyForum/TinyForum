package repository

import (
	"errors"
	"tiny-forum/internal/model"

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
		Title:      input.Title,
		Content:    input.Content,
		Summary:    input.Summary,
		Cover:      input.Cover,
		BoardID:    input.BoardID,
		AuthorID:   userID,
		IsQuestion: true,
		Status:     "published",
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
