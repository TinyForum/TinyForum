package topic

import (
	topicService "tiny-forum/internal/service/topic"
)

type TopicHandler struct {
	topicSvc topicService.TopicService
}

func NewTopicHandler(topicSvc topicService.TopicService) *TopicHandler {
	return &TopicHandler{topicSvc: topicSvc}
}
