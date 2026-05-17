package bot

import (
	"encoding/json"
	"tiny-forum/internal/infra/lua/nocode"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
)

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

func parseFlowRequestRaw(raw any) *nocode.Flow {
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

func toFlow(req *request.ValidateFlowRequest) *nocode.Flow {
	return &nocode.Flow{
		Version:    req.Version,
		Trigger:    req.Trigger,
		Conditions: req.Conditions,
		Actions:    req.Actions,
	}
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
