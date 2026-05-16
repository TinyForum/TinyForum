package bot

import (
	"errors"
	"tiny-forum/internal/infra/lua/nocode"
)

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
