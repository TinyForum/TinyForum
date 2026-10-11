package bot

import (
	"context"
	"tiny-forum/internal/infra/lua/nocode"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/logger"
)

// ─── 调度器 ───────────────────────────────────────────────────────────────

func (s *service) StartScheduler() {
	ctx := context.Background()
	bots, err := s.repo.ListActive(ctx)
	if err != nil {
		logger.Errorf("调度器加载机器人失败: %v", err)
		return
	}
	for _, bot := range bots {
		s.registerBot(bot)
	}
	s.cron.Start()
	logger.Infof("[Scheduler] 启动，已注册 %d 个机器人", len(bots))
}

func (s *service) registerBot(bot *do.Bot) {
	b := bot

	// 优先检查 nocode 流程中的触发器类型
	if flowRaw, ok := bot.ConfigValues["flow"]; ok {
		flow := parseFlowRequestRaw(flowRaw)
		if flow != nil {
			switch flow.Trigger.Type {
			case nocode.TriggerOnSchedule:
				if cron, ok := flow.Trigger.Params["cron"].(string); ok && cron != "" {
					_, err := s.cron.AddFunc(cron, func() { s.executeBot(b, nil) })
					if err != nil {
						logger.Errorf("[Scheduler] bot=%d nocode cron='%s' 注册失败: %v", bot.ID, cron, err)
					}
				}
			case nocode.TriggerOnNewPost:
				s.eventBus.Subscribe("post.created", func(data map[string]any) { s.executeBot(b, data) })
			case nocode.TriggerOnNewComment:
				s.eventBus.Subscribe("comment.created", func(data map[string]any) { s.executeBot(b, data) })
			case nocode.TriggerOnUserRegister:
				s.eventBus.Subscribe("user.registered", func(data map[string]any) { s.executeBot(b, data) })
			case nocode.TriggerOnKeyword:
				s.eventBus.Subscribe("post.created", func(data map[string]any) { s.executeBot(b, data) })
				s.eventBus.Subscribe("comment.created", func(data map[string]any) { s.executeBot(b, data) })
			}
			return
		}
	}

	// 回退到顶层 trigger 字段
	switch bot.TriggerType {
	case do.TriggerSchedule:
		if bot.CronExpr != "" {
			_, err := s.cron.AddFunc(bot.CronExpr, func() { s.executeBot(b, nil) })
			if err != nil {
				logger.Errorf("[Scheduler] bot=%d cron='%s' 注册失败: %v", bot.ID, bot.CronExpr, err)
			}
		}
	case do.TriggerEvent:
		if bot.EventFilter != "" {
			s.eventBus.Subscribe(bot.EventFilter, func(data map[string]any) { s.executeBot(b, data) })
		}
	}
}

func (s *service) StopScheduler() {
	s.cron.Stop()
}
