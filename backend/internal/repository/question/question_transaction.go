package question

import (
	"errors"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"

	"gorm.io/gorm"
)

// CreateWithTransaction 使用事务创建问答（包括帖子、标签、积分扣减）
func (r *questionRepository) CreateWithTransaction(userID uint, input dto.CreateQuestionRequest) (*do.QuestionResponse, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 创建帖子
	post := &do.Post{
		Title:      input.Title,
		Content:    input.Content,
		Summary:    input.Summary,
		Cover:      input.Cover,
		BoardID:    input.BoardID,
		AuthorID:   userID,
		Type:       do.PostTypeQuestion,
		PostStatus: input.Status,
	}
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. 关联标签
	if len(input.TagIDs) > 0 {
		var tags []do.Tag
		if err := tx.Where("id IN ?", input.TagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(post).Association("Tags").Append(&tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 3. 扣减积分
	if input.RewardScore > 0 {
		var user do.User
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
	question := &do.Question{
		PostID:      post.ID,
		RewardScore: input.RewardScore,
		AnswerCount: 0,
	}
	if err := tx.Create(question).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &do.QuestionResponse{
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
		Status:      string(post.PostStatus),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}, nil
}
