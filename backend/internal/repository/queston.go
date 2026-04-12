package repository

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type QuestionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) Create(question *model.Question) error {
	return r.db.Create(question).Error
}

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
