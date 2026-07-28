package topic

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
)

// AddPostToTopic 添加帖子到话题
func (s *topicService) AddPostToTopic(input request.AddPostToTopicRequest, userID uint) error {
	topic, err := s.topicRepo.FindByID(input.TopicID)
	if err != nil {
		return apperrors.ErrTopicNotFound
	}
	// if topic.CreatorID != userID {
	// 	return errors.New("只有话题创建者可以添加内容")
	// }
	post, err := s.postRepo.FindByArticleID(input.PostID)
	if err != nil {
		return apperrors.ErrPostNotFound
	}
	topicPost := &do.TopicCreation{

		TopicID:    input.TopicID,
		CreationID: input.PostID,
		// ID:         input.TopicID,
		// PostID:    input.PostID,
		SortOrder: input.SortOrder, // 排序
		CreatorID: userID,
	}
	if err := s.topicRepo.AddPost(topicPost); err != nil {
		return err
	}
	_ = s.topicRepo.IncrementPostCount(input.TopicID)
	if post.Creation.AuthorID != userID {
		s.notifSvc.Create(post.Creation.AuthorID, &userID, do.NotifySystem,
			"你的帖子被收录到专题《"+topic.Title+"》", &input.TopicID, "topic")
	}
	return nil
}

// RemovePostFromTopic 从专题移除帖子
func (s *topicService) RemovePostFromTopic(topicID, postID uint, userID uint) error {
	topic, err := s.topicRepo.FindByID(topicID)
	if err != nil {
		return apperrors.ErrTopicNotFound
	}
	// if topic.CreatorID != userID {
	// 	return errors.New("只有专题创建者可以移除内容")
	// }
	logger.Infof("移除话题： %s", topic)
	if err := s.topicRepo.RemovePost(topicID, postID); err != nil {
		return err
	}
	return s.topicRepo.DecrementPostCount(topicID)
}

// GetTopicPosts 获取专题下的帖子列表（分页）
func (s *topicService) GetTopicPosts(topicID uint, page, pageSize int) ([]do.TopicCreation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetTopicPosts(topicID, pageSize, offset)
}
