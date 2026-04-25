package wire

import (
	"tiny-forum/config"
	"tiny-forum/internal/service/announcement"
	"tiny-forum/internal/service/auth"
	"tiny-forum/internal/service/board"
	"tiny-forum/internal/service/comment"
	"tiny-forum/internal/service/email"
	"tiny-forum/internal/service/notification"
	"tiny-forum/internal/service/post"
	"tiny-forum/internal/service/question"
	"tiny-forum/internal/service/risk"
	"tiny-forum/internal/service/stats"
	"tiny-forum/internal/service/tag"
	"tiny-forum/internal/service/timeline"
	"tiny-forum/internal/service/topic"
	"tiny-forum/internal/service/user"
	jwtpkg "tiny-forum/pkg/jwt"
)

// Services 聚合所有 Service
type Services struct {
	User         user.UserService
	Auth         auth.AuthService
	Tag          tag.TagService
	Notification notification.NotificationService
	Post         post.PostService
	Comment      comment.CommentService
	Board        *board.BoardService
	Timeline     *timeline.TimelineService
	Topic        *topic.TopicService
	Question     *question.QuestionService
	Announcement announcement.AnnouncementService
	Stats        *stats.StatsService
	Risk         *risk.RiskService
	ContentCheck *risk.ContentCheckService
}

// NewServices 创建所有 Service 实例
func NewServices(
	cfg *config.Config,
	jwtMgr *jwtpkg.Manager,
	repos *Repositories,
	infra *Infra,
) *Services {
	// 风险服务（先创建，因为其他服务可能依赖）
	riskSvc := risk.NewRiskService(repos.Risk, infra.RateLimiter)
	checkSvc := risk.NewContentCheckService(repos.Risk, infra.SensitiveFilter)

	// 基础服务
	notifSvc := notification.NewNotificationService(repos.Notification)
	userSvc := user.NewUserService(repos.User, jwtMgr, notifSvc)
	tagSvc := tag.NewTagService(repos.Tag)
	boardSvc := board.NewBoardService(repos.Board, repos.User, repos.Post, notifSvc)
	timelineSvc := timeline.NewTimelineService(repos.Timeline, repos.User, repos.Post, repos.Comment)
	topicSvc := topic.NewTopicService(repos.Topic, repos.Post, repos.User, notifSvc)
	questionSvc := question.NewQuestionService(repos.Question, repos.Post, repos.Comment, repos.User, notifSvc, repos.Tag)
	postSvc := post.NewPostService(repos.Post, repos.Tag, repos.User, repos.Board, notifSvc, checkSvc)
	commentSvc := comment.NewCommentService(repos.Comment, repos.Post, repos.User, notifSvc, repos.Vote)
	announcementSvc := announcement.NewAnnouncementService(repos.Announcement)
	statsSvc := stats.NewStatsService(repos.Stats, repos.Post, repos.Tag, repos.Board, repos.User, repos.Comment)
	emailSvc := email.NewEmailService(&cfg.Private.Email)
	authSvc := auth.NewAuthService(repos.Auth, repos.User, jwtMgr, notifSvc, emailSvc, cfg, repos.Token, repos.Transaction)

	// 辅助

	return &Services{
		User:         userSvc,
		Auth:         authSvc,
		Tag:          tagSvc,
		Notification: notifSvc,
		Post:         postSvc,
		Comment:      commentSvc,
		Board:        boardSvc,
		Timeline:     timelineSvc,
		Topic:        topicSvc,
		Question:     questionSvc,
		Announcement: announcementSvc,
		Stats:        statsSvc,
		Risk:         riskSvc,
		ContentCheck: checkSvc,
	}
}
