package service

import (
	"encoding/json"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

type TimelineService struct {
	timelineRepo *repository.TimelineRepository
	userRepo     *repository.UserRepository
	postRepo     repository.PostRepository
	commentRepo  *repository.CommentRepository
}

func NewTimelineService(
	timelineRepo *repository.TimelineRepository,
	userRepo *repository.UserRepository,
	postRepo repository.PostRepository,
	commentRepo *repository.CommentRepository,
) *TimelineService {
	return &TimelineService{
		timelineRepo: timelineRepo,
		userRepo:     userRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
	}
}

type CreateEventInput struct {
	UserID     uint
	ActorID    uint
	Action     model.ActionType
	TargetID   uint
	TargetType string
	Payload    interface{}
	Score      int
}

func (s *TimelineService) CreateEvent(input CreateEventInput) error {
	payloadJSON, _ := json.Marshal(input.Payload)

	event := &model.TimelineEvent{
		UserID:     input.UserID,
		ActorID:    input.ActorID,
		Action:     input.Action,
		TargetID:   input.TargetID,
		TargetType: input.TargetType,
		Payload:    string(payloadJSON),
		Score:      input.Score,
	}

	return s.timelineRepo.CreateEvent(event)
}

func (s *TimelineService) GetHomeTimeline(userID uint, page, pageSize int) ([]model.TimelineEvent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Update last read
	s.timelineRepo.UpdateLastRead(userID, "home")

	return s.timelineRepo.GetUserTimeline(userID, pageSize, offset)
}

func (s *TimelineService) GetFollowingTimeline(userID uint, page, pageSize int) ([]model.TimelineEvent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Update last read
	s.timelineRepo.UpdateLastRead(userID, "following")

	return s.timelineRepo.GetFollowingTimeline(userID, pageSize, offset)
}

func (s *TimelineService) Subscribe(subscriberID, targetUserID uint) error {
	// Check if target user exists
	_, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return err
	}

	sub := &model.TimelineSubscription{
		SubscriberID: subscriberID,
		TargetUserID: targetUserID,
		TargetType:   "user",
		IsActive:     true,
	}

	return s.timelineRepo.Subscribe(sub)
}

func (s *TimelineService) Unsubscribe(subscriberID, targetUserID uint) error {
	return s.timelineRepo.Unsubscribe(subscriberID, targetUserID)
}

func (s *TimelineService) GetSubscriptions(subscriberID uint) ([]model.TimelineSubscription, error) {
	return s.timelineRepo.GetSubscriptions(subscriberID)
}

func (s *TimelineService) IsSubscribed(subscriberID, targetUserID uint) (bool, error) {
	return s.timelineRepo.IsSubscribed(subscriberID, targetUserID)
}
