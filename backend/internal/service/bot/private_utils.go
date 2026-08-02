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
		ErrorMsg:     bot.ErrorMsg,
		LastExecLogs: bot.LastExecLogs, LastExecDurMs: bot.LastExecDurationMs,
		ConfigSchema: bot.ConfigSchema, ConfigValues: bot.ConfigValues,
		CreatedAt: bot.CreatedAt, UpdatedAt: bot.UpdatedAt,
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
		Version: req.Version,
		Trigger: req.Trigger,
		Steps:   req.Steps,
		Actions: req.Actions,
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

// sanitizeError 脱敏错误信息：截断 SQL 语句、去除路径前缀、限长
func sanitizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()

	// 截断 GORM SQL 错误：取最外层描述
	if i := len(msg); i > 300 {
		msg = msg[:300] + "..."
	}

	// 替换文件系统路径为 <path>
	msg = replacePathPrefix(msg)

	return msg
}

// sanitizeLogs 确保日志非 nil，并限长
func sanitizeLogs(logs []string) []string {
	if logs == nil {
		return []string{}
	}
	// 截断过长日志
	for i, l := range logs {
		if len(l) > 500 {
			logs[i] = l[:500] + "..."
		}
	}
	return logs
}

func replacePathPrefix(s string) string {
	// 替换常见的路径模式
	importPath := "tiny-forum/"
	if idx := indexAfter(s, importPath); idx > 0 {
		s = s[idx:]
	}
	// 替换文件系统绝对路径
	if len(s) > 80 && s[0] == '/' {
		s = "<path>/" + lastSegment(s)
	}
	return s
}

func indexAfter(s, prefix string) int {
	for i := 0; i < len(s); i++ {
		if i+len(prefix) <= len(s) && s[i:i+len(prefix)] == prefix {
			return i
		}
	}
	return -1
}

func lastSegment(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return s[i+1:]
		}
	}
	return s
}
