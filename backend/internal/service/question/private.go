package question

import (
	"errors"
	"tiny-forum/internal/model/dto"
)

// validateCreateQuestion 验证创建问答的输入
func (s *questionService) validateCreateQuestion(input dto.CreateQuestionRequest) error {
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
