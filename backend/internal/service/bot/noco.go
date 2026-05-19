package bot

import (
	"errors"
	"tiny-forum/internal/infra/lua/nocode"
)

// ─── 零代码 ───────────────────────────────────────────────────────────────

func (s *service) GetNocodeMetadata() *nocode.NocodeMetadata {
	// 获取所有支持的 nocode 元数据
	return &nocode.NocodeMetadata{
		Triggers:   nocode.BuiltinTriggers,
		Actions:    nocode.BuiltinActions,
		Conditions: nocode.BuiltinConditions,
	}
}

// func (s *service) ValidateFlowRequest(req *request.ValidateFlowRequest) []error {
// 	// if flow == nil {
// 	// 	return []error{errors.New("flow is nil")}
// 	// }
// 	var errs []error
// 	if req.Trigger.Type == "" {
// 		errs = append(errs, errors.New("trigger.type is required"))
// 	}
// 	if len(req.Actions) == 0 {
// 		errs = append(errs, errors.New("at least one action is required"))
// 	}
// 	return errs
// }

func (s *service) ValidateFlow(req *nocode.Flow) []error {
	// if flow == nil {
	// 	return []error{errors.New("flow is nil")}
	// }
	var errs []error
	if req.Trigger.Type == "" {
		errs = append(errs, errors.New("trigger.type is required"))
	}
	if len(req.Actions) == 0 {
		errs = append(errs, errors.New("at least one action is required"))
	}
	return errs
}
