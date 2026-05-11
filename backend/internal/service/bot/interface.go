// Package bot 提供机器人业务逻辑服务。
// 两种执行路径：
//   - Lua 脚本：bot.ScriptCode != "" → LuaSandbox.ExecuteWithAPI
//   - 零代码：bot.ConfigValues["flow"] 存在 → FlowEngine.Run
package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"tiny-forum/internal/botapi"
	"tiny-forum/internal/infra/lua/engine"
	"tiny-forum/internal/infra/lua/nocode"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	botrepo "tiny-forum/internal/repository/bot"
	commentrepo "tiny-forum/internal/repository/comment"
	notificationrepo "tiny-forum/internal/repository/notification"
	postrepo "tiny-forum/internal/repository/post"
	userrepo "tiny-forum/internal/repository/user"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
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
	postRepo    postrepo.PostRepository
	commentRepo commentrepo.CommentRepository
	userRepo    userrepo.UserRepository
	notifRepo   notificationrepo.NotificationRepository
}

// NewService 创建 bot Service。依赖现有各 repository，由 wire/service.go 注入。
func NewService(
	repo botrepo.Repository,
	postRepo postrepo.PostRepository,
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

// ─── CRUD ─────────────────────────────────────────────────────────────────

func (s *service) Create(ctx context.Context, creatorID uint, req *request.CreateBotRequest) (*do.Bot, error) {
	// 零代码模式：预先校验 Flow JSON
	if req.ScriptCode == "" && req.ConfigValues != nil {
		if flowRaw, ok := req.ConfigValues["flow"]; ok {
			if errs := s.ValidateFlow(parseFlowRaw(flowRaw)); len(errs) > 0 {
				return nil, fmt.Errorf("invalid nocode flow: %v", errs[0])
			}
		}
	}

	// 查询创建者用户名
	creatorName := "user"
	if u, err := s.userRepo.FindByID(creatorID); err == nil {
		creatorName = u.Username
	}

	bot := &do.Bot{
		Name:          req.Name,
		Version:       req.Version,
		Description:   req.Description,
		Summary:       req.Summary,
		AvatarURL:     req.AvatarURL,
		Screenshots:   orStrSlice(req.Screenshots),
		HomepageURL:   req.HomepageURL,
		Type:          req.Type,
		Tags:          orStrSlice(req.Tags),
		CreatorID:     creatorID,
		CreatorName:   creatorName,
		ScriptCode:    req.ScriptCode,
		ScriptURL:     req.ScriptURL,
		TriggerType:   req.TriggerType,
		CronExpr:      req.CronExpr,
		EventFilter:   req.EventFilter,
		TimeoutSec:    req.TimeoutSec,
		RetryTimes:    req.RetryTimes,
		EnvVars:       orStrMap(req.EnvVars),
		ResourceLimit: req.ResourceLimit,
		Pricing:       req.Pricing,
		Permissions:   req.Permissions,
		Enabled:       false,
		Status:        do.BotStatusInactive,
		ConfigSchema:  req.ConfigSchema,
		ConfigValues:  orAnyMap(req.ConfigValues),
	}
	if bot.TimeoutSec == 0 {
		bot.TimeoutSec = 10
	}
	if err := s.repo.Create(ctx, bot); err != nil {
		return nil, err
	}
	return bot, nil
}

func (s *service) Update(ctx context.Context, userID uint, botID uint, req *request.UpdateBotRequest) error {
	bot, err := s.repo.GetByID(ctx, botID)
	if err != nil {
		return err
	}
	if bot.CreatorID != userID {
		return errors.New("permission denied")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}
	if req.ScriptCode != nil {
		updates["script_code"] = *req.ScriptCode
	}
	if req.ScriptURL != nil {
		updates["script_url"] = *req.ScriptURL
	}
	if req.TriggerType != nil {
		updates["trigger_type"] = *req.TriggerType
	}
	if req.CronExpr != nil {
		updates["cron_expr"] = *req.CronExpr
	}
	if req.EventFilter != nil {
		updates["event_filter"] = *req.EventFilter
	}
	if req.TimeoutSec != nil {
		updates["timeout_sec"] = *req.TimeoutSec
	}
	if req.RetryTimes != nil {
		updates["retry_times"] = *req.RetryTimes
	}
	if req.EnvVars != nil {
		updates["env_vars"] = req.EnvVars
	}
	if req.ResourceLimit != nil {
		updates["resource_limit"] = req.ResourceLimit
	}
	if req.Pricing != nil {
		updates["pricing"] = req.Pricing
	}
	if req.Permissions != nil {
		updates["permissions"] = req.Permissions
	}
	if req.ConfigSchema != nil {
		updates["config_schema"] = req.ConfigSchema
	}
	if req.ConfigValues != nil {
		updates["config_values"] = req.ConfigValues
	}
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
	return toResponse(bot), nil
}

func (s *service) ListByUser(ctx context.Context, userID uint, page, pageSize int) ([]*vo.BotResponse, int64, error) {
	offset := (page - 1) * pageSize
	bots, total, err := s.repo.ListByUser(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return mapToResponse(bots), total, nil
}

func (s *service) List(ctx context.Context, page, pageSize int) ([]*vo.BotResponse, int64, error) {
	offset := (page - 1) * pageSize
	bots, total, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return mapToResponse(bots), total, nil
}

// ─── 执行 ─────────────────────────────────────────────────────────────────

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(bot.TimeoutSec+10)*time.Second)
	defer cancel()

	start := time.Now()
	var execErr error
	var logs []string

	// 构造 ForumAPI（每次执行隔离，携带权限列表）
	perms := make([]string, len(bot.Permissions))
	for i, p := range bot.Permissions {
		perms[i] = string(p)
	}
	api := botapi.NewForumAPI(do.SystemBotID, s.postRepo, s.commentRepo, s.userRepo, s.notifRepo)

	if bot.ScriptCode != "" {
		// ── Lua 脚本 ──────────────────────────────────────────────────
		result := s.sandbox.Execute(ctx, bot, api, eventData)
		execErr = result.Err
		logs = result.Logs

	} else if flowRaw, ok := bot.ConfigValues["flow"]; ok {
		// ── 零代码流程 ────────────────────────────────────────────────
		flow := parseFlowRaw(flowRaw)
		if flow == nil {
			execErr = errors.New("invalid flow configuration")
		} else {
			engine := nocode.NewFlowEngine(api)
			fctx, err := engine.Run(ctx, flow, eventData)
			if fctx != nil {
				logs = fctx.Logs
			}
			execErr = err
		}
	} else {
		execErr = errors.New("bot has neither script_code nor nocode flow")
	}

	duration := time.Since(start).Milliseconds()
	s.recordLog(bot.ID, execErr == nil, duration, logs, execErr)

	// 更新 bot 状态
	bgCtx := context.Background()
	if execErr != nil {
		_ = s.repo.Update(bgCtx, bot.ID, map[string]interface{}{
			"status":    do.BotStatusError,
			"error_msg": execErr.Error(),
		})
	} else {
		_ = s.repo.Update(bgCtx, bot.ID, map[string]interface{}{
			"exec_count":   gorm.Expr("exec_count + 1"),
			"last_exec_at": time.Now(),
			"status":       do.BotStatusActive,
			"error_msg":    "",
		})
	}
}

func (s *service) recordLog(botID uint, success bool, durationMs int64, logs []string, err error) {
	// TODO: 写入 bot_execution_logs 表
	status := "success"
	errMsg := ""
	if !success {
		status = "fail"
		if err != nil {
			errMsg = err.Error()
		}
	}
	fmt.Printf("[BotLog] bot=%d status=%s duration=%dms logs=%d err=%s\n",
		botID, status, durationMs, len(logs), errMsg)
}

// ─── 零代码 ───────────────────────────────────────────────────────────────

func (s *service) GetNocodeMetadata() *nocode.NocodeMetadata {
	// 获取所有支持的 nocode 元数据
	return &nocode.NocodeMetadata{
		Triggers:   nocode.BuiltinActions,
		Actions:    nocode.BuiltinActions,
		Conditions: nocode.BuiltinConditions,
	}
}

func (s *service) ValidateFlow(flow *nocode.Flow) []error {
	if flow == nil {
		return []error{errors.New("flow is nil")}
	}
	var errs []error
	if flow.Trigger.Type == "" {
		errs = append(errs, errors.New("trigger.type is required"))
	}
	if len(flow.Actions) == 0 {
		errs = append(errs, errors.New("at least one action is required"))
	}
	return errs
}

// ─── 调度器 ───────────────────────────────────────────────────────────────

func (s *service) StartScheduler() {
	ctx := context.Background()
	bots, err := s.repo.ListActive(ctx)
	if err != nil {
		fmt.Println("[Scheduler] ListActive error:", err)
		return
	}
	for _, bot := range bots {
		s.registerBot(bot)
	}
	s.cron.Start()
	fmt.Printf("[Scheduler] 启动，已注册 %d 个机器人\n", len(bots))
}

func (s *service) registerBot(bot *do.Bot) {
	switch bot.TriggerType {
	case do.TriggerSchedule:
		if bot.CronExpr != "" {
			b := bot
			_, err := s.cron.AddFunc(bot.CronExpr, func() { s.executeBot(b, nil) })
			if err != nil {
				fmt.Printf("[Scheduler] bot=%d cron='%s' 注册失败: %v\n", bot.ID, bot.CronExpr, err)
			}
		}
	case do.TriggerEvent:
		if bot.EventFilter != "" {
			b := bot
			s.eventBus.Subscribe(bot.EventFilter, func(data map[string]any) { s.executeBot(b, data) })
		}
	}
}

func (s *service) StopScheduler() {
	s.cron.Stop()
}

// ─── EventBus ─────────────────────────────────────────────────────────────

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

func (s *service) PublishEvent(eventName string, data map[string]any) {
	s.eventBus.Publish(eventName, data)
}

// ─── 工具函数 ─────────────────────────────────────────────────────────────

func toResponse(bot *do.Bot) *vo.BotResponse {
	return &vo.BotResponse{
		ID: bot.ID, Name: bot.Name, Version: bot.Version,
		Description: bot.Description, Summary: bot.Summary,
		AvatarURL: bot.AvatarURL, Screenshots: bot.Screenshots,
		HomepageURL: bot.HomepageURL, Type: bot.Type, Tags: bot.Tags,
		CreatorID: bot.CreatorID, CreatorName: bot.CreatorName,
		TriggerType: bot.TriggerType, CronExpr: bot.CronExpr,
		EventFilter: bot.EventFilter, TimeoutSec: bot.TimeoutSec,
		RetryTimes: bot.RetryTimes, ResourceLimit: bot.ResourceLimit,
		Pricing: bot.Pricing, Permissions: bot.Permissions,
		Enabled: bot.Enabled, Status: bot.Status,
		ExecCount: bot.ExecCount, LastExecAt: bot.LastExecAt,
		ErrorMsg: bot.ErrorMsg, ConfigSchema: bot.ConfigSchema,
		ConfigValues: bot.ConfigValues,
		CreatedAt:    bot.CreatedAt, UpdatedAt: bot.UpdatedAt,
	}
}

func mapToResponse(bots []*do.Bot) []*vo.BotResponse {
	res := make([]*vo.BotResponse, 0, len(bots))
	for _, b := range bots {
		res = append(res, toResponse(b))
	}
	return res
}

func parseFlowRaw(raw any) *nocode.Flow {
	var s string
	switch v := raw.(type) {
	case string:
		s = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		s = string(b)
	}
	f, err := nocode.FlowFromJSON(s)
	if err != nil {
		return nil
	}
	return f
}

func orStrSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func orStrMap(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}

func orAnyMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	return m
}
