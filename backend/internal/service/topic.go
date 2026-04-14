package service

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
	"tiny-forum/pkg/logger"
)

type TopicService struct {
	topicRepo *repository.TopicRepository
	postRepo  repository.PostRepository
	userRepo  *repository.UserRepository
	notifSvc  *NotificationService
}

func NewTopicService(
	topicRepo *repository.TopicRepository,
	postRepo repository.PostRepository,
	userRepo *repository.UserRepository,
	notifSvc *NotificationService,
) *TopicService {
	return &TopicService{
		topicRepo: topicRepo,
		postRepo:  postRepo,
		userRepo:  userRepo,
		notifSvc:  notifSvc,
	}
}

type CreateTopicInput struct {
	Title       string `json:"title" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	Cover       string `json:"cover" binding:"max=500"`
	IsPublic    bool   `json:"is_public"`
}

func (s *TopicService) Create(creatorID uint, input CreateTopicInput) (*model.Topic, error) {
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

func (s *TopicService) Update(id uint, input CreateTopicInput) (*model.Topic, error) {
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

func (s *TopicService) Delete(id uint, userID uint, isAdmin bool) error {
	topic, err := s.topicRepo.FindByID(id)
	if err != nil {
		return errors.New("专题不存在")
	}
	if topic.CreatorID != userID && !isAdmin {
		return errors.New("无权限删除此专题")
	}
	return s.topicRepo.Delete(id)
}

func (s *TopicService) GetByID(id uint) (*model.Topic, error) {
	return s.topicRepo.FindByID(id)
}

func (s *TopicService) List(page, pageSize int) ([]model.Topic, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.List(pageSize, offset)
}

func (s *TopicService) GetByCreator(creatorID uint, page, pageSize int) ([]model.Topic, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetByCreator(creatorID, pageSize, offset)
}

type AddPostToTopicInput struct {
	TopicID   uint `json:"topic_id" binding:"required"`
	PostID    uint `json:"post_id" binding:"required"`
	SortOrder int  `json:"sort_order"`
}

func (s *TopicService) AddPostToTopic(input AddPostToTopicInput, userID uint) error {
	// Check if topic exists
	topic, err := s.topicRepo.FindByID(input.TopicID)
	if err != nil {
		return errors.New("专题不存在")
	}

	// Check permission (creator or admin)
	if topic.CreatorID != userID {
		return errors.New("只有专题创建者可以添加内容")
	}

	// Check if post exists
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return errors.New("帖子不存在")
	}

	topicPost := &model.TopicPost{
		TopicID:   input.TopicID,
		PostID:    input.PostID,
		SortOrder: input.SortOrder,
		AddedBy:   userID,
	}

	if err := s.topicRepo.AddPost(topicPost); err != nil {
		return err
	}

	// Increment post count
	s.topicRepo.IncrementPostCount(input.TopicID)

	// Notify post author
	if post.AuthorID != userID {
		s.notifSvc.Create(post.AuthorID, &userID, model.NotifySystem,
			"你的帖子被收录到专题《"+topic.Title+"》", &input.TopicID, "topic")
	}

	return nil
}

func (s *TopicService) RemovePostFromTopic(topicID, postID uint, userID uint) error {
	topic, err := s.topicRepo.FindByID(topicID)
	if err != nil {
		return errors.New("专题不存在")
	}
	if topic.CreatorID != userID {
		return errors.New("只有专题创建者可以移除内容")
	}

	if err := s.topicRepo.RemovePost(topicID, postID); err != nil {
		return err
	}

	return s.topicRepo.DecrementPostCount(topicID)
}

func (s *TopicService) GetTopicPosts(topicID uint, page, pageSize int) ([]model.TopicPost, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetTopicPosts(topicID, pageSize, offset)
}

func (s *TopicService) Follow(userID, topicID uint) error {
	topic, err := s.topicRepo.FindByID(topicID)
	logger.Info(fmt.Sprintf("Topic: %+v", topic))

	if err != nil {
		return errors.New("专题不存在")
	}

	isFollowing, _ := s.topicRepo.IsFollowing(userID, topicID)
	if isFollowing {
		return errors.New("已经关注过了")
	}

	follow := &model.TopicFollow{
		UserID:  userID,
		TopicID: topicID,
	}

	if err := s.topicRepo.Follow(follow); err != nil {
		return err
	}

	return s.topicRepo.IncrementFollowerCount(topicID)
}

func (s *TopicService) Unfollow(userID, topicID uint) error {
	topic, err := s.topicRepo.FindByID(topicID)
	logger.Info(fmt.Sprintf("Topic: %+v", topic))

	if err != nil {
		return errors.New("专题不存在")
	}

	if err := s.topicRepo.Unfollow(userID, topicID); err != nil {
		return err
	}

	return s.topicRepo.DecrementFollowerCount(topicID)
}

func (s *TopicService) IsFollowing(userID, topicID uint) (bool, error) {
	return s.topicRepo.IsFollowing(userID, topicID)
}

func (s *TopicService) GetFollowers(topicID uint, page, pageSize int) ([]model.TopicFollow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetFollowers(topicID, pageSize, offset)
}
