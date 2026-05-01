package question

import "tiny-forum/internal/model/do"

func (r *questionRepository) CreateAnswerVote(vote *do.AnswerVote) error {
	return r.db.Create(vote).Error
}

func (r *questionRepository) UpdateAnswerVote(vote *do.AnswerVote) error {
	return r.db.Save(vote).Error
}

func (r *questionRepository) DeleteAnswerVote(userID, commentID uint) error {
	return r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&do.AnswerVote{}).Error
}

func (r *questionRepository) FindAnswerVote(userID, commentID uint) (*do.AnswerVote, error) {
	var vote do.AnswerVote
	err := r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).
		First(&vote).Error
	return &vote, err
}

func (r *questionRepository) GetAnswerVoteCount(commentID uint) (int, error) {
	var upCount, downCount int64
	r.db.Model(&do.AnswerVote{}).Where("comment_id = ? AND vote_type = ?", commentID, "up").Count(&upCount)
	r.db.Model(&do.AnswerVote{}).Where("comment_id = ? AND vote_type = ?", commentID, "down").Count(&downCount)
	return int(upCount - downCount), nil
}
