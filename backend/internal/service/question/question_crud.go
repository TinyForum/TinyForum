package question

import (
	"errors"
	"fmt"

	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// CreateQuestion 创建问答帖
func (s *QuestionService) CreateQuestion(userID uint, input model.CreateQuestionInput) (*model.QuestionResponse, error) {
	if err := s.validateCreateQuestion(input); err != nil {
		return nil, err
	}
	question, err := s.questionRepo.CreateWithTransaction(userID, input)
	if err != nil {
		return nil, fmt.Errorf("创建问答失败: %w", err)
	}
	return question, nil
}

// validateCreateQuestion 验证创建问答的输入
func (s *QuestionService) validateCreateQuestion(input model.CreateQuestionInput) error {
	if input.Title == "" {
		return errors.New("标题不能为空")
	}
	if input.Content == "" {
		return errors.New("内容不能为空")
	}
	if len(input.Title) > 100 {
		return errors.New("标题长度不能超过100个字符")
	}
	if len(input.Summary) > 500 {
		return errors.New("摘要长度不能超过500个字符")
	}
	if input.RewardScore < 0 || input.RewardScore > 100 {
		return errors.New("悬赏积分必须在0-100之间")
	}
	return nil
}

// GetQuestionDetail 获取问答帖详情
func (s *QuestionService) GetQuestionDetail(questionID uint) (*model.QuestionResponse, error) {
	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("问答不存在")
		}
		return nil, fmt.Errorf("查询问答失败: %w", err)
	}
	return &model.QuestionResponse{
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
		Status:           string(question.Post.Status),
		CreatedAt:        question.CreatedAt,
		UpdatedAt:        question.UpdatedAt,
	}, nil
}

// GetQuestionsList 获取问答帖列表（支持只看未回答）
func (s *QuestionService) GetQuestionsList(page, pageSize int, unanswered bool) ([]model.Post, int64, error) {
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
func (s *QuestionService) GetQuestionByID(questionID uint) (*model.Question, error) {
	return s.questionRepo.FindByID(questionID)
}
