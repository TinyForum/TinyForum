package question

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"

	"gorm.io/gorm"
)

// CreateQuestion 创建问答帖
func (s *questionService) CreateQuestion(userID uint, input dto.CreateQuestionRequest) (*do.QuestionResponse, error) {
	if err := s.validateCreateQuestion(input); err != nil {
		return nil, err
	}
	question, err := s.questionRepo.CreateWithTransaction(userID, input)
	if err != nil {
		return nil, fmt.Errorf("创建问答失败: %w", err)
	}
	return question, nil
}

// GetQuestionDetail 获取问答帖详情
func (s *questionService) GetQuestionDetail(questionID uint) (*do.QuestionResponse, error) {
	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("问答不存在")
		}
		return nil, fmt.Errorf("查询问答失败: %w", err)
	}
	return &do.QuestionResponse{
		ID:               question.ID,
		PostID:           question.PostID,
		Title:            question.Post.Title,
		Content:          question.Post.Content,
		Summary:          question.Post.Summary,
		Cover:            question.Post.Cover,
		BoardID:          question.Post.BoardID,
		AuthorID:         question.Post.AuthorID,
		RewardScore:      question.RewardScore,
		AnswerCount:      question.AnswerCount,
		AcceptedAnswerID: question.AcceptedAnswerID,
		Status:           string(question.Post.PostStatus),
		CreatedAt:        question.CreatedAt,
		UpdatedAt:        question.UpdatedAt,
	}, nil
}

// GetQuestionsList 获取问答帖列表（支持只看未回答）
func (s *questionService) GetQuestionsList(page, pageSize int, unanswered bool) ([]do.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	if unanswered {
		return s.postRepo.GetUnansweredQuestions(pageSize, offset)
	}
	return s.postRepo.GetQuestions(pageSize, offset)
}

// GetQuestionByID 根据 ID 获取 Question 模型（不含 Post 详情）
func (s *questionService) GetQuestionByID(questionID uint) (*do.Question, error) {
	return s.questionRepo.FindByID(questionID)
}
