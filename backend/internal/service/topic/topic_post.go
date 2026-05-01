package topic

import (
	"errors"
	"tiny-forum/internal/model/do"
)

type AddPostToTopicInput struct {
	TopicID   uint `json:"topic_id" binding:"required"`
	PostID    uint `json:"post_id" binding:"required"`
	SortOrder int  `json:"sort_order"`
}

// AddPostToTopic 添加帖子到专题
func (s *topicService) AddPostToTopic(input AddPostToTopicInput, userID uint) error {
	topic, err := s.topicRepo.FindByID(input.TopicID)
	if err != nil {
		return errors.New("专题不存在")
	}
	if topic.CreatorID != userID {
		return errors.New("只有专题创建者可以添加内容")
	}
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	topicPost := &do.TopicPost{
		TopicID:   input.TopicID,
		PostID:    input.PostID,
		SortOrder: input.SortOrder,
		AddedBy:   userID,
	}
	if err := s.topicRepo.AddPost(topicPost); err != nil {
		return err
	}
	_ = s.topicRepo.IncrementPostCount(input.TopicID)
	if post.AuthorID != userID {
		s.notifSvc.Create(post.AuthorID, &userID, do.NotifySystem,
			"你的帖子被收录到专题《"+topic.Title+"》", &input.TopicID, "topic")
	}
	return nil
}

// RemovePostFromTopic 从专题移除帖子
func (s *topicService) RemovePostFromTopic(topicID, postID uint, userID uint) error {
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

// GetTopicPosts 获取专题下的帖子列表（分页）
func (s *topicService) GetTopicPosts(topicID uint, page, pageSize int) ([]do.TopicPost, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetTopicPosts(topicID, pageSize, offset)
}
