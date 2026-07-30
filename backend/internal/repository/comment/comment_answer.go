package comment

import "tiny-forum/internal/model/do"

// MarkAsAccepted 标记评论为已采纳答案
func (r *commentRepository) MarkAsAccepted(answerID uint) error {
	return r.db.Model(&do.Answer{}).Where("id = ?", answerID).
		Update("is_accepted", true).Error
}

// MarkAsAnswer 标记/取消标记评论为答案
func (r *commentRepository) MarkAsAnswer(commentID uint, isAnswer bool) error {
	return r.db.Model(&do.Comment{}).Where("id = ?", commentID).
		Update("is_answer", isAnswer).Error
}

// UnacceptAnswer 取消接受答案
func (r *commentRepository) UnacceptAnswer(answerID uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 更新 answers 表的 is_accepted 字段
	if err := tx.Model(&do.Answer{}).
		Where("id = ?", answerID).
		Update("is_accepted", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 更新 questions 表的 accepted_answer_id 字段
	if err := tx.Model(&do.Question{}).
		Where("accepted_answer_id = ?", answerID).
		Update("accepted_answer_id", nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetAcceptedAnswer 获取问题已接受的答案
func (r *commentRepository) GetAcceptedAnswer(postID uint) (*do.Answer, error) {
	var comment do.Answer
	err := r.db.Where("post_id = ? AND is_accepted = ?", postID, true).
		First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
