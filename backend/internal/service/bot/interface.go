// Package bot 提供机器人业务逻辑服务。
// 两种执行路径：
//   - Lua 脚本：bot.ScriptCode != "" → LuaSandbox.ExecuteWithAPI
//   - 零代码：bot.ConfigValues["flow"] 存在 → FlowEngine.Run
package bot

import (
	"context"

	"tiny-forum/internal/infra/lua/engine"
	"tiny-forum/internal/infra/lua/nocode"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	postrepo "tiny-forum/internal/repository/article"
	botrepo "tiny-forum/internal/repository/bot"
	commentrepo "tiny-forum/internal/repository/comment"
	notificationrepo "tiny-forum/internal/repository/notification"
	userrepo "tiny-forum/internal/repository/user"

	"github.com/robfig/cron/v3"
)

// ─── Service 接口 ─────────────────────────────────────────────────────────

type Service interface {
	Create(ctx context.Context, creatorID uint, req *request.CreateBotRequest) (*do.Bot, error)
	Update(ctx context.Context, userID uint, botID uint, req *request.UpdateBotRequest) error
	Delete(ctx context.Context, userID uint, botID uint) error
	Get(ctx context.Context, id uint) (*vo.BotResponse, error)
	ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*vo.BotResponse, int64, error)
	List(ctx context.Context, page, pageSize int) ([]*vo.BotResponse, int64, error)
	RunNow(ctx context.Context, botID uint, eventData map[string]any) error

	// 零代码支持
	GetNocodeMetadata() *nocode.NocodeMetadata
	ValidateFlow(flow *nocode.Flow) []error

	// 事件总线（其他 service 发布事件触发机器人）
	PublishEvent(eventName string, data map[string]any)

	// 调度器生命周期
	StartScheduler()
	StopScheduler()
}

// ─── service 实现 ─────────────────────────────────────────────────────────

type service struct {
	repo        botrepo.Repository
	sandbox     *engine.LuaSandbox
	cron        *cron.Cron
	eventBus    *EventBus
	postRepo    postrepo.ArticleRepository
	commentRepo commentrepo.CommentRepository
	userRepo    userrepo.UserRepository
	notifRepo   notificationrepo.NotificationRepository
}

// NewService 创建 bot Service。依赖现有各 repository，由 wire/service.go 注入。
func NewService(
	repo botrepo.Repository,
	postRepo postrepo.ArticleRepository,
	commentRepo commentrepo.CommentRepository,
	userRepo userrepo.UserRepository,
	notifRepo notificationrepo.NotificationRepository,
) Service {
	return &service{
		repo:        repo,
		sandbox:     engine.NewLuaSandbox(30, []string{"0.0.0.0", "127.0.0.1"}),
		cron:        cron.New(),
		eventBus:    NewEventBus(),
		postRepo:    postRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		notifRepo:   notifRepo,
	}
}
