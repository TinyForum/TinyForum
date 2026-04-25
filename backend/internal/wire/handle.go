package wire

import (
	announcementHandler "tiny-forum/internal/handler/announcement"
	answerHandler "tiny-forum/internal/handler/answer"
	authHandler "tiny-forum/internal/handler/auth"
	boardHandler "tiny-forum/internal/handler/board"
	commentHandler "tiny-forum/internal/handler/comment"
	notificationHandler "tiny-forum/internal/handler/notification"
	postHandler "tiny-forum/internal/handler/post"
	questionHandler "tiny-forum/internal/handler/questions"
	riskhandler "tiny-forum/internal/handler/risk"
	statsHandler "tiny-forum/internal/handler/stats"
	tagHandler "tiny-forum/internal/handler/tags"
	timelineHandler "tiny-forum/internal/handler/timelines"
	topicHandler "tiny-forum/internal/handler/topic"
	userHandler "tiny-forum/internal/handler/user"
	"tiny-forum/pkg/timeutil"
)

// Handlers 聚合所有 Handler
type Handlers struct {
	Auth         *authHandler.AuthHandler
	User         *userHandler.UserHandler
	Tag          *tagHandler.TagHandler
	Notification *notificationHandler.NotificationHandler
	Post         *postHandler.PostHandler
	Comment      *commentHandler.CommentHandler
	Board        *boardHandler.BoardHandler
	Timeline     *timelineHandler.TimelineHandler
	Topic        *topicHandler.TopicHandler
	Question     *questionHandler.QuestionHandler
	Answer       *answerHandler.AnswerHandler
	Announcement *announcementHandler.AnnouncementHandler
	Stats        *statsHandler.StatsHandler
	Risk         *riskhandler.RiskHandler
}

// NewHandlers 创建所有 Handler 实例
func NewHandlers(svc *Services, timeHelpers *timeutil.TimeHelpers) *Handlers {

	return &Handlers{
		Auth:         authHandler.NewAuthHandler(svc.Auth),
		User:         userHandler.NewUserHandler(svc.User, svc.Notification, svc.Auth),
		Tag:          tagHandler.NewTagHandler(svc.Tag),
		Notification: notificationHandler.NewNotificationHandler(svc.Notification),
		Post:         postHandler.NewPostHandler(svc.Post),
		Comment:      commentHandler.NewCommentHandler(svc.Comment, svc.Question),
		Board:        boardHandler.NewBoardHandler(svc.Board),
		Timeline:     timelineHandler.NewTimelineHandler(svc.Timeline),
		Topic:        topicHandler.NewTopicHandler(svc.Topic),
		Question:     questionHandler.NewQuestionHandler(svc.Question, svc.Comment, svc.Post),
		Answer:       answerHandler.NewAnswerHandler(svc.Question, svc.Comment, svc.Post),
		Announcement: announcementHandler.NewAnnouncementHandler(svc.Announcement),
		Stats: statsHandler.NewStatsHandler(
			svc.Stats,
			timeHelpers,
		),
		Risk: riskhandler.NewRiskHandler(svc.ContentCheck, svc.Risk),
	}
}
