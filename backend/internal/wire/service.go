package wire

import (
	"tiny-forum/internal/infra/config"
	luasdk "tiny-forum/internal/infra/lua/sdk"
	"tiny-forum/internal/service/admin"
	"tiny-forum/internal/service/announcement"
	"tiny-forum/internal/service/article"
	attachment "tiny-forum/internal/service/attachment"
	"tiny-forum/internal/service/auth"
	"tiny-forum/internal/service/board"
	"tiny-forum/internal/service/bot"
	"tiny-forum/internal/service/check"
	"tiny-forum/internal/service/comment"
	"tiny-forum/internal/service/email"
	"tiny-forum/internal/service/notification"
	"tiny-forum/internal/service/plugin"
	"tiny-forum/internal/service/question"
	"tiny-forum/internal/service/risk"
	"tiny-forum/internal/service/stats"
	"tiny-forum/internal/service/tag"
	"tiny-forum/internal/service/timeline"
	"tiny-forum/internal/service/topic"
	"tiny-forum/internal/service/upload"
	"tiny-forum/internal/service/user"
	"tiny-forum/internal/service/violation"
	"tiny-forum/internal/storage"
	"tiny-forum/internal/strategy"
	jwtpkg "tiny-forum/pkg/jwt"
)

// Services 聚合所有 Service
type Services struct {
	User         user.UserService
	Auth         auth.AuthService
	Tag          tag.TagService
	Notification notification.NotificationService
	Article      article.ArticleService
	Comment      comment.CommentService
	Board        board.BoardService
	Timeline     timeline.TimelineService
	Topic        topic.TopicService
	Question     question.QuestionService
	Announcement announcement.AnnouncementService
	Stats        stats.StatsService
	Risk         risk.RiskService
	ContentCheck check.ContentCheckService
	Attachment   attachment.AttachmentService
	Admin        admin.AdminService
	Plugin       plugin.PluginService
	Bot          bot.Service
}

// NewServices 创建所有 Service 实例
func NewServices(
	cfg *config.Config,

	jwtMgr *jwtpkg.JWTManager,
	repos *Repositories,
	infra *Infra,
	userStorage storage.StorageDriver,
	publicStorage storage.StorageDriver,
	registry *strategy.HandlerRegistry,
	forumAPI luasdk.ForumAPI,

) *Services {
	// FIXME: 优先考虑横向调用（service），再考虑纵向调用（repo）
	// registry := strategy.NewHandlerRegistry()
	// userStorage := storage.NewLocalStorage("./uploads")
	// publicStorage := storage.NewLocalStorage("./public")
	riskSvc := risk.NewRiskService(repos.Risk, infra.RateLimiter)
	checkSvc := check.NewContentCheckService(repos.Risk, *infra.sensitiveChecker)
	// 基础服务
	notifSvc := notification.NewNotificationService(repos.Notification)
	violation := violation.NewViolationService(repos.Violation)
	userSvc := user.NewUserService(repos.User, jwtMgr, notifSvc, repos.Post, repos.Comment, violation)
	tagSvc := tag.NewTagService(repos.Tag)
	boardSvc := board.NewBoardService(repos.Board, repos.User, repos.Post, notifSvc)
	timelineSvc := timeline.NewTimelineService(repos.Timeline, repos.User, repos.Post, repos.Comment)
	topicSvc := topic.NewTopicService(repos.Topic, repos.Post, repos.User, notifSvc)
	questionSvc := question.NewQuestionService(repos.Question, repos.Post, repos.Comment, repos.User, notifSvc, repos.Tag)
	articleSvc := article.NewPostService(repos.Post, repos.Tag, repos.User, repos.Board, notifSvc, checkSvc)
	commentSvc := comment.NewCommentService(repos.Comment, repos.Post, repos.User, notifSvc, repos.Vote)
	announcementSvc := announcement.NewAnnouncementService(repos.Announcement)
	statsSvc := stats.NewStatsService(repos.Stats, repos.Post, repos.Tag, repos.Board, repos.User, repos.Comment)
	emailSvc := email.NewEmailService(&cfg.Private.Email)
	authSvc := auth.NewAuthService(repos.Auth, repos.User, jwtMgr, notifSvc, emailSvc, cfg, repos.Token, repos.Transaction, infra.RedisClient)
	adminSvc := admin.NewAdminService(announcementSvc, userSvc, articleSvc, boardSvc)
	pluginSvc := plugin.NewPluginService(repos.Plugin, publicStorage, &cfg.Basic.Plugins)
	engine := upload.NewEngine(userStorage, registry)
	attachmentSvc := attachment.NewAttachmentService(repos.Attachment, cfg.Basic.Attachment, engine)
	botSvc := bot.NewService(repos.Bot, repos.Post, repos.Comment, repos.User, repos.Notification)

	return &Services{
		User:         userSvc,
		Auth:         authSvc,
		Tag:          tagSvc,
		Notification: notifSvc,
		Article:      articleSvc,
		Comment:      commentSvc,
		Board:        boardSvc,
		Timeline:     timelineSvc,
		Topic:        topicSvc,
		Question:     questionSvc,
		Announcement: announcementSvc,
		Stats:        statsSvc,
		Risk:         riskSvc,
		ContentCheck: checkSvc,
		Attachment:   attachmentSvc,
		Admin:        adminSvc,
		Plugin:       pluginSvc,
		Bot:          botSvc,
	}
}
