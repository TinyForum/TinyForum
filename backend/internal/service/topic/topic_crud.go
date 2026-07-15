package topic

import (
	"errors"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

// Create 创建专题
func (s *topicService) Create(creatorID uint, input request.CreateTopicReqeust) (*do.Topic, error) {
	topic := &do.Topic{
		Title:       input.Title,
		Description: input.Description,
		CoverUrl:    input.CoverUrl,
		CreatorID:   creatorID,
		IsPublic:    input.IsPublic,
	}
	if err := s.topicRepo.Create(topic); err != nil {
		return nil, err
	}
	return s.topicRepo.FindByID(topic.ID)
}

// Update 更新专题
func (s *topicService) Update(id uint, input request.CreateTopicReqeust) (*do.Topic, error) {
	topic, err := s.topicRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("专题不存在")
	}
	topic.Title = input.Title
	topic.Description = input.Description
	topic.CoverUrl = input.CoverUrl
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
func (s *topicService) GetByID(id uint) (*do.Topic, error) {
	return s.topicRepo.FindByID(id)
}

// List 分页获取专题列表
func (s *topicService) List(page, pageSize int) ([]do.Topic, int64, error) {
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
func (s *topicService) GetByCreator(creatorID uint, page, pageSize int) ([]do.Topic, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetByCreator(creatorID, pageSize, offset)
}
