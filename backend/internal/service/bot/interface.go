package bot

import (
	"context"
	"errors"
	"fmt"
	"time"
	luaengine "tiny-forum/internal/infra/lua"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	"tiny-forum/internal/repository/bot"

	"github.com/robfig/cron/v3"
	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, creatorID uint, req *request.CreateBotRequest) (*do.Bot, error)
	Update(ctx context.Context, userID uint, botID uint, req *request.UpdateBotRequest) error
	Delete(ctx context.Context, userID uint, botID uint) error
	Get(ctx context.Context, id uint) (*vo.BotResponse, error)
	ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*vo.BotResponse, int64, error)
	List(ctx context.Context, page, pageSize int) ([]*vo.BotResponse, int64, error)
	RunNow(ctx context.Context, botID uint, eventData map[string]any) error
	StartScheduler() // 启动后台调度器
	StopScheduler()
}

type service struct {
	repo     bot.Repository
	sandbox  *luaengine.LuaSandbox
	cron     *cron.Cron
	eventBus *EventBus // 自定义事件总线，简单实现见下方
	// db       *gorm.DB
}

func NewService(repo bot.Repository) Service {
	s := &service{
		repo:     repo,
		sandbox:  luaengine.NewLuaSandbox(10),
		cron:     cron.New(),
		eventBus: NewEventBus(),
		// db:       db,
	}
	return s
}

// Create 创建机器人
func (s *service) Create(ctx context.Context, creatorID uint, req *request.CreateBotRequest) (*do.Bot, error) {
	// 获取创建者名称（略，假设通过 repo 获取 user）
	bot := &do.Bot{
		Name:          req.Name,
		Version:       req.Version,
		Description:   req.Description,
		Summary:       req.Summary,
		AvatarURL:     req.AvatarURL,
		Screenshots:   req.Screenshots,
		HomepageURL:   req.HomepageURL,
		Type:          req.Type,
		Tags:          req.Tags,
		CreatorID:     creatorID,
		CreatorName:   "user", // 实际需查询用户表
		ScriptCode:    req.ScriptCode,
		ScriptURL:     req.ScriptURL,
		TriggerType:   req.TriggerType,
		CronExpr:      req.CronExpr,
		EventFilter:   req.EventFilter,
		TimeoutSec:    req.TimeoutSec,
		RetryTimes:    req.RetryTimes,
		EnvVars:       req.EnvVars,
		ResourceLimit: req.ResourceLimit,
		Pricing:       req.Pricing,
		Permissions:   req.Permissions,
		Enabled:       false,
		Status:        do.BotStatusInactive,
		ConfigSchema:  req.ConfigSchema,
		ConfigValues:  req.ConfigValues,
	}
	if bot.TimeoutSec == 0 {
		bot.TimeoutSec = 10
	}
	if err := s.repo.Create(ctx, bot); err != nil {
		return nil, err
	}
	return bot, nil
}

// Update 更新机器人（仅允许自己创建的）
func (s *service) Update(ctx context.Context, userID uint, botID uint, req *request.UpdateBotRequest) error {
	bot, err := s.repo.GetByID(ctx, botID)
	if err != nil {
		return err
	}
	if bot.CreatorID != userID {
		return errors.New("permission denied")
	}
	updates := make(map[string]interface{})
	// 字段映射略（可根据 req 非空指针填充）
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	// ... 其他字段类似
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
		if *req.Enabled {
			updates["status"] = do.BotStatusActive
		} else {
			updates["status"] = do.BotStatusInactive
		}
	}
	return s.repo.Update(ctx, botID, updates)
}

func (s *service) Delete(ctx context.Context, userID uint, botID uint) error {
	bot, err := s.repo.GetByID(ctx, botID)
	if err != nil {
		return err
	}
	if bot.CreatorID != userID {
		return errors.New("permission denied")
	}
	return s.repo.Delete(ctx, botID)
}

func (s *service) Get(ctx context.Context, id uint) (*vo.BotResponse, error) {
	bot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(bot), nil
}

// func (s *service) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*vo.BotResponse, int64, error) {
// 	offset := (page - 1) * pageSize
// 	bots, total, err := s.repo.List(ctx, userID, offset, pageSize)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	var res []*vo.BotResponse
// 	for _, b := range bots {
// 		res = append(res, s.toResponse(b))
// 	}
// 	return res, total, nil
// }

// RunNow 手动执行机器人
func (s *service) RunNow(ctx context.Context, botID uint, eventData map[string]any) error {
	bot, err := s.repo.GetByID(ctx, botID)
	if err != nil {
		return err
	}
	if !bot.Enabled {
		return errors.New("bot not enabled")
	}
	go s.executeBot(bot, eventData)
	return nil
}

func (s *service) executeBot(bot *do.Bot, eventData map[string]any) {
	start := time.Now()
	var execErr error
	var output interface{}

	// 准备 API 回调（需根据机器人权限包装）
	apiCallbacks := map[string]interface{}{
		"getPost":   s.makeGetPostCallback(bot),
		"replyPost": s.makeReplyPostCallback(bot),
	}

	result, err := s.sandbox.Execute(bot.ScriptCode, bot.ConfigValues, eventData, apiCallbacks)
	execErr = err
	output = result

	duration := time.Since(start).Milliseconds()
	// 记录执行日志
	s.recordLog(bot.ID, execErr == nil, duration, fmt.Sprintf("%v", output), execErr)
	if execErr != nil {
		s.repo.Update(context.Background(), bot.ID, map[string]interface{}{
			"status":    do.BotStatusError,
			"error_msg": execErr.Error(),
		})
	} else {
		s.repo.Update(context.Background(), bot.ID, map[string]interface{}{
			"exec_count":   gorm.Expr("exec_count + 1"),
			"last_exec_at": time.Now(),
			"status":       do.BotStatusActive,
			"error_msg":    "",
		})
	}
}

func (s *service) recordLog(botID uint, success bool, duration int64, output string, err error) {
	// 实际应写入数据库表 bot_execution_logs
	// 此处省略具体存储
}

func (s *service) toResponse(bot *do.Bot) *vo.BotResponse {
	return &vo.BotResponse{
		ID:            bot.ID,
		Name:          bot.Name,
		Version:       bot.Version,
		Description:   bot.Description,
		Summary:       bot.Summary,
		AvatarURL:     bot.AvatarURL,
		Screenshots:   bot.Screenshots,
		HomepageURL:   bot.HomepageURL,
		Type:          bot.Type,
		Tags:          bot.Tags,
		CreatorID:     bot.CreatorID,
		CreatorName:   bot.CreatorName,
		TriggerType:   bot.TriggerType,
		CronExpr:      bot.CronExpr,
		EventFilter:   bot.EventFilter,
		TimeoutSec:    bot.TimeoutSec,
		RetryTimes:    bot.RetryTimes,
		ResourceLimit: bot.ResourceLimit,
		Pricing:       bot.Pricing,
		Permissions:   bot.Permissions,
		Enabled:       bot.Enabled,
		Status:        bot.Status,
		ExecCount:     bot.ExecCount,
		LastExecAt:    bot.LastExecAt,
		ErrorMsg:      bot.ErrorMsg,
		ConfigSchema:  bot.ConfigSchema,
		ConfigValues:  bot.ConfigValues,
		CreatedAt:     bot.CreatedAt,
		UpdatedAt:     bot.UpdatedAt,
	}
}

// 启动调度器：从数据库加载所有 enabled bots 并注册 cron/event
func (s *service) StartScheduler() {
	ctx := context.Background()
	bots, err := s.repo.ListActive(ctx)
	if err != nil {
		return
	}
	for _, bot := range bots {
		s.registerBot(bot)
	}
	s.cron.Start()
}

func (s *service) registerBot(bot *do.Bot) {
	switch bot.TriggerType {
	case do.TriggerSchedule:
		if bot.CronExpr != "" {
			s.cron.AddFunc(bot.CronExpr, func() {
				s.executeBot(bot, nil)
			})
		}
	case do.TriggerEvent:
		if bot.EventFilter != "" {
			s.eventBus.Subscribe(bot.EventFilter, func(data map[string]any) {
				s.executeBot(bot, data)
			})
		}
	}
}

func (s *service) StopScheduler() {
	s.cron.Stop()
}

// 简单事件总线实现
type EventBus struct {
	subscribers map[string][]func(map[string]any)
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string][]func(map[string]any))}
}

func (e *EventBus) Subscribe(event string, fn func(map[string]any)) {
	e.subscribers[event] = append(e.subscribers[event], fn)
}

func (e *EventBus) Publish(event string, data map[string]any) {
	for _, fn := range e.subscribers[event] {
		go fn(data)
	}
}

// 回调函数示例（需权限检查）
func (s *service) makeGetPostCallback(bot *do.Bot) func(*lua.LState) int {
	return func(L *lua.LState) int {
		postID := L.CheckInt64(1)
		// 检查权限：是否包含 read:posts
		// 模拟返回一个表
		tbl := L.NewTable()
		tbl.RawSetString("id", lua.LNumber(postID))
		tbl.RawSetString("title", lua.LString("demo post"))
		L.Push(tbl)
		return 1
	}
}

func (s *service) makeReplyPostCallback(bot *do.Bot) func(*lua.LState) int {
	return func(L *lua.LState) int {
		postID := L.CheckInt64(1)
		content := L.CheckString(2)
		// 检查权限 write:posts
		_ = postID
		_ = content
		L.Push(lua.LBool(true))
		return 1
	}
}
