package wire

import (
	"tiny-forum/config"
	adminHandler "tiny-forum/internal/handler/admin"
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
	uploadHandler "tiny-forum/internal/handler/upload"
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
	Upload       *uploadHandler.UploadHandler
	Admin        *adminHandler.AdminHandler
}

// NewHandlers 创建所有 Handler 实例
func NewHandlers(svc *Services, timeHelpers *timeutil.TimeHelpers, cfg *config.Config) *Handlers {
	auth := authHandler.NewAuthHandler(svc.Auth, cfg)
	user := userHandler.NewUserHandler(svc.User, svc.Notification, svc.Auth)
	tag := tagHandler.NewTagHandler(svc.Tag)
	notification := notificationHandler.NewNotificationHandler(svc.Notification)
	post := postHandler.NewPostHandler(svc.Post)
	comment := commentHandler.NewCommentHandler(svc.Comment, svc.Question)
	board := boardHandler.NewBoardHandler(svc.Board)
	timeline := timelineHandler.NewTimelineHandler(svc.Timeline)
	topic := topicHandler.NewTopicHandler(svc.Topic)
	answer := answerHandler.NewAnswerHandler(svc.Question, svc.Comment, svc.Post)
	question := questionHandler.NewQuestionHandler(svc.Question, svc.Comment, svc.Post, answer)
	announcement := announcementHandler.NewAnnouncementHandler(svc.Announcement)
	stats := statsHandler.NewStatsHandler(svc.Stats, timeHelpers)
	risk := riskhandler.NewRiskHandler(svc.ContentCheck, svc.Risk)
	upload := uploadHandler.NewUploadHandler(svc.Upload)
	admin := adminHandler.NewAdminHandler(svc.admin)

	return &Handlers{
		Auth:         auth,
		User:         user,
		Tag:          tag,
		Notification: notification,
		Post:         post,
		Comment:      comment,
		Board:        board,
		Timeline:     timeline,
		Topic:        topic,
		Question:     question,
		Answer:       answer,
		Announcement: announcement,
		Stats:        stats,
		Risk:         risk,
		Upload:       upload,
		Admin:        admin,
	}
}
