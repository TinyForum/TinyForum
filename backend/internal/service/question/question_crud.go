package question

import (
	"errors"

	// "fmt"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// CreateQuestion 创建问答帖
func (s *questionService) CreateQuestion(userID uint, input dto.CreateQuestionRequest) (*do.QuestionResponse, error) {
	if err := s.validateCreateQuestion(input); err != nil {
		return nil, err
	}
	question, err := s.questionRepo.CreateWithTransaction(userID, input)
	if err != nil {
		return nil, apperrors.ErrCreateQuestionFailed
	}
	return question, nil
}

// GetQuestionDetail 获取问答帖详情
func (s *questionService) GetQuestionDetail(questionID uint) (*do.QuestionResponse, error) {
	question, err := s.questionRepo.FindByQuestionID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrQuestionNotFound
		}
		return nil, apperrors.ErrQueryQuestionFailed
	}
	return &do.QuestionResponse{
		ID:               question.ID,
		CreationsID:      question.CreationID,
		Title:            question.Creation.Title,
		Content:          question.Creation.Content,
		Summary:          question.Creation.Summary,
		Cover:            question.Creation.CoverUrl,
		BoardID:          question.Creation.BoardID,
		AuthorID:         question.Creation.AuthorID,
		RewardScore:      question.RewardScore,
		AnswerCount:      question.AnswerCount,
		AcceptedAnswerID: question.AcceptedAnswerID,
		Status:           string(question.Creation.CreationStatus),
		CreatedAt:        question.CreatedAt,
		UpdatedAt:        question.UpdatedAt,
	}, nil
}

// GetQuestionsList 获取问答帖列表（支持只看未回答）
func (s *questionService) GetQuestionsList(page, pageSize int, unanswered bool) ([]do.Question, int64, error) {
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
	return s.questionRepo.FindByQuestionID(questionID)
}
