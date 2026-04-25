package topic

import (
	"errors"

	"tiny-forum/internal/model"
)

type CreateTopicInput struct {
	Title       string `json:"title" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	Cover       string `json:"cover" binding:"max=500"`
	IsPublic    bool   `json:"is_public"`
}

// Create 创建专题
func (s *topicService) Create(creatorID uint, input CreateTopicInput) (*model.Topic, error) {
	topic := &model.Topic{
		Title:       input.Title,
		Description: input.Description,
		Cover:       input.Cover,
		CreatorID:   creatorID,
		IsPublic:    input.IsPublic,
	}
	if err := s.topicRepo.Create(topic); err != nil {
		return nil, err
	}
	return s.topicRepo.FindByID(topic.ID)
}

// Update 更新专题
func (s *topicService) Update(id uint, input CreateTopicInput) (*model.Topic, error) {
	topic, err := s.topicRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("专题不存在")
	}
	topic.Title = input.Title
	topic.Description = input.Description
	topic.Cover = input.Cover
	topic.IsPublic = input.IsPublic
	if err := s.topicRepo.Update(topic); err != nil {
		return nil, err
	}
	return topic, nil
}

// Delete 删除专题
func (s *topicService) Delete(id uint, userID uint, isAdmin bool) error {
	topic, err := s.topicRepo.FindByID(id)
	if err != nil {
		return errors.New("专题不存在")
	}
	if topic.CreatorID != userID && !isAdmin {
		return errors.New("无权限删除此专题")
	}
	return s.topicRepo.Delete(id)
}

// GetByID 获取专题详情
func (s *topicService) GetByID(id uint) (*model.Topic, error) {
	return s.topicRepo.FindByID(id)
}

// List 分页获取专题列表
func (s *topicService) List(page, pageSize int) ([]model.Topic, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.List(pageSize, offset)
}

// GetByCreator 获取用户创建的专题列表
func (s *topicService) GetByCreator(creatorID uint, page, pageSize int) ([]model.Topic, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetByCreator(creatorID, pageSize, offset)
}
