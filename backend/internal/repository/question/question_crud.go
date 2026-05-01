package question

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *questionRepository) Create(question *do.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) Update(question *do.Question) error {
	return r.db.Save(question).Error
}

func (r *questionRepository) FindByID(id uint) (*do.Question, error) {
	var question do.Question
	err := r.db.Preload("Post").Preload("Post.Tags").Where("id = ?", id).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) FindByPostID(postID uint) (*do.Question, error) {
	var question do.Question
	err := r.db.Where("post_id = ?", postID).
		Preload("Post").
		Preload("AcceptedAnswer").
		First(&question).Error
	return &question, err
}

func (r *questionRepository) IncrementAnswerCount(postID uint) error {
	return r.db.Model(&do.Question{}).Where("post_id = ?", postID).
		UpdateColumn("answer_count", gorm.Expr("answer_count + 1")).Error
}

func (r *questionRepository) SetAcceptedAnswer(postID, commentID uint) error {
	return r.db.Model(&do.Question{}).Where("post_id = ?", postID).
		Updates(map[string]interface{}{
			"accepted_answer_id": commentID,
		}).Error
}

func (r *questionRepository) UpdateCommentVoteCount(commentID uint, voteCount int) error {
	return r.db.Model(&do.Comment{}).Where("id = ?", commentID).
		UpdateColumn("vote_count", voteCount).Error
}

func (r *questionRepository) UpdateAnswerCount(questionID uint) error {
	return r.db.Model(&do.Question{}).
		Where("id = ?", questionID).
		Update("answer_count", gorm.Expr("answer_count + ?", 1)).Error
}

func (r *questionRepository) UpdateAcceptedAnswer(questionID uint, answerID uint) error {
	return r.db.Model(&do.Question{}).
		Where("id = ?", questionID).
		Update("accepted_answer_id", answerID).Error
}
