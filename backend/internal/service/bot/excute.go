package bot

import (
	"context"
	"errors"
	"fmt"
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
