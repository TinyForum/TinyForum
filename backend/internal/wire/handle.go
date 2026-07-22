// internal/wire/types.go
package wire

import (
	"log"
	adminHandler "tiny-forum/internal/handler/admin"
	announcementHandler "tiny-forum/internal/handler/announcement"
	answerHandler "tiny-forum/internal/handler/answer"
	articleHandler "tiny-forum/internal/handler/article"
	"tiny-forum/internal/handler/attachment"
	authHandler "tiny-forum/internal/handler/auth"
	boardHandler "tiny-forum/internal/handler/board"
	botHandler "tiny-forum/internal/handler/bot"
	commentHandler "tiny-forum/internal/handler/comment"
	configHandler "tiny-forum/internal/handler/config"
	notificationHandler "tiny-forum/internal/handler/notification"
	pluginHandler "tiny-forum/internal/handler/plugin"
	questionHandler "tiny-forum/internal/handler/questions"
	riskhandler "tiny-forum/internal/handler/risk"
	statsHandler "tiny-forum/internal/handler/stats"
	tagHandler "tiny-forum/internal/handler/tags"
	timelineHandler "tiny-forum/internal/handler/timelines"
	topicHandler "tiny-forum/internal/handler/topic"
	userHandler "tiny-forum/internal/handler/user"
	"tiny-forum/internal/infra/config"
	configService "tiny-forum/internal/service/config"
	"tiny-forum/pkg/timeutil"
)

// Handlers 聚合所有 Handler
type Handlers struct {
	Auth         *authHandler.AuthHandler
	User         *userHandler.UserHandler
	Tag          *tagHandler.TagHandler
	Notification *notificationHandler.NotificationHandler
	Article      *articleHandler.ArticleHandler
	Comment      *commentHandler.CommentHandler
	Board        *boardHandler.BoardHandler
	Timeline     *timelineHandler.TimelineHandler
	Topic        *topicHandler.TopicHandler
	Question     *questionHandler.QuestionHandler
	Answer       *answerHandler.AnswerHandler
	Announcement *announcementHandler.AnnouncementHandler
	Stats        *statsHandler.StatsHandler
	Risk         *riskhandler.RiskHandler
	Attachment   *attachment.AttachmentHandler
	Admin        *adminHandler.AdminHandler
	Plugin       *pluginHandler.Handler
	Config       *configHandler.ConfigHandler
	Bot          *botHandler.Handler
}

// NewHandlers 创建所有 Handler 实例
func NewHandlers(svc *Services, timeHelpers *timeutil.TimeHelpers, cfg *config.Config, configSvc *configService.ConfigService) *Handlers {

	auth := authHandler.NewAuthHandler(svc.Auth, cfg)
	user := userHandler.NewUserHandler(svc.User, svc.Notification, svc.Auth)
	tag := tagHandler.NewTagHandler(svc.Tag)
	notification := notificationHandler.NewNotificationHandler(svc.Notification)
	article := articleHandler.NewPostHandler(svc.Article)
	comment := commentHandler.NewCommentHandler(svc.Comment, svc.Question)
	board := boardHandler.NewBoardHandler(svc.Board)
	timeline := timelineHandler.NewTimelineHandler(svc.Timeline)
	topic := topicHandler.NewTopicHandler(svc.Topic)
	answer := answerHandler.NewAnswerHandler(svc.Question, svc.Comment, svc.Article)
	question := questionHandler.NewQuestionHandler(svc.Question, svc.Comment, svc.Article, answer)
	announcement := announcementHandler.NewAnnouncementHandler(svc.Announcement)
	stats := statsHandler.NewStatsHandler(svc.Stats, timeHelpers)
	risk := riskhandler.NewRiskHandler(svc.ContentCheck, svc.Risk)
	attachment := attachment.NewAttachmentHandler(svc.Attachment)
	admin := adminHandler.NewAdminHandler(svc.Admin)
	plugin := pluginHandler.NewHandler(svc.Plugin)
	bot := botHandler.NewHandler(svc.Bot)
	config := configHandler.NewConfigHandler(configSvc)
	return &Handlers{
		Auth:         auth,
		User:         user,
		Tag:          tag,
		Notification: notification,
		Article:      article,
		Comment:      comment,
		Board:        board,
		Timeline:     timeline,
		Topic:        topic,
		Question:     question,
		Answer:       answer,
		Announcement: announcement,
		Stats:        stats,
		Risk:         risk,
		Attachment:   attachment,
		Admin:        admin,
		Plugin:       plugin,
		Bot:          bot,
		Config:       config,
	}
}

// UpdateConfig 更新所有 Handler 的配置
// 当配置文件变更时，会调用此方法将新配置传播到各个 Handler
func (h *Handlers) UpdateConfig(cfg *config.Config) {
	log.Printf("[Handlers] Updating all handlers with new config")

	// 更新需要配置的 Handler
	// 注意：只有那些真正需要动态配置的 Handler 才需要更新

	// 1. AuthHandler - 可能需要 JWT 配置更新
	if h.Auth != nil {
		if updater, ok := interface{}(h.Auth).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 2. UserHandler - 可能需要用户相关配置更新
	if h.User != nil {
		if updater, ok := interface{}(h.User).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 3. RiskHandler - 风控配置更新（重要）
	if h.Risk != nil {
		if updater, ok := interface{}(h.Risk).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 4. PostHandler - 帖子相关配置更新
	if h.Article != nil {
		if updater, ok := interface{}(h.Article).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 5. CommentHandler - 评论相关配置更新
	if h.Comment != nil {
		if updater, ok := interface{}(h.Comment).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 6. BoardHandler - 版块配置更新
	if h.Board != nil {
		if updater, ok := interface{}(h.Board).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 7. StatsHandler - 统计配置更新
	if h.Stats != nil {
		if updater, ok := interface{}(h.Stats).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 8. AttachmentHandler - 附件配置更新
	if h.Attachment != nil {
		if updater, ok := interface{}(h.Attachment).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 9. PluginHandler - 插件配置更新
	if h.Plugin != nil {
		if updater, ok := interface{}(h.Plugin).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 10. BotHandler - 机器人配置更新
	if h.Bot != nil {
		if updater, ok := interface{}(h.Bot).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 11. AdminHandler - 管理配置更新
	if h.Admin != nil {
		if updater, ok := interface{}(h.Admin).(interface{ UpdateConfig(*config.Config) }); ok {
			updater.UpdateConfig(cfg)
		}
	}

	// 其他 Handler 根据需求添加...

	log.Printf("[Handlers] Config update completed")
}
