package bot

import (
	"context"
	"errors"
	"time"
	"tiny-forum/internal/botapi"
	"tiny-forum/internal/infra/lua/nocode"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

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
	api := botapi.NewForumAPI(do.SystemBotID, s.postRepo, s.commentRepo, s.userRepo, s.notifRepo, s.attachmentRepo)

	if flowRaw, ok := bot.ConfigValues["flow"]; ok {
		// ── 零代码流程（优先于 Lua） ──────────────────────────────────
		flow := parseFlowRequestRaw(flowRaw)
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

	} else if bot.ScriptCode != "" {
		// ── Lua 脚本 ──────────────────────────────────────────────────
		result := s.sandbox.Execute(ctx, bot, api, eventData)
		execErr = result.Err
		logs = result.Logs

	} else {
		execErr = errors.New("bot has neither nocode flow nor script_code")
	}

	duration := time.Since(start).Milliseconds()

	// 更新 bot 状态 + 执行日志
	bgCtx := context.Background()
	updates := map[string]interface{}{
		"last_exec_duration_ms": duration,
		"last_exec_logs":        sanitizeLogs(logs),
	}
	if execErr != nil {
		updates["status"] = do.BotStatusError
		updates["error_msg"] = sanitizeError(execErr)
	} else {
		updates["exec_count"] = gorm.Expr("exec_count + 1")
		updates["last_exec_at"] = time.Now()
		updates["status"] = do.BotStatusActive
		updates["error_msg"] = ""
	}
	_ = s.repo.Update(bgCtx, bot.ID, updates)
}
