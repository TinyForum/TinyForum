package bot

import (
	"context"
	"fmt"
	"tiny-forum/internal/model/do"
)

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
