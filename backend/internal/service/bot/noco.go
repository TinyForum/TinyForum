package bot

import (
	"errors"
	"tiny-forum/internal/infra/lua/nocode"
)

// ─── 零代码 ───────────────────────────────────────────────────────────────

func (s *service) GetNocodeMetadata() *nocode.NocodeMetadata {
	return &nocode.NocodeMetadata{
		Triggers:  nocode.BuiltinTriggers,
		Control:   nocode.BuiltinControl,
		Variables: nocode.BuiltinVariables,
		Actions:   nocode.BuiltinActions,
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
	var errs []error
	if req.Trigger.Type == "" {
		errs = append(errs, errors.New("trigger.type is required"))
	}
	if len(req.Steps) == 0 && len(req.Actions) == 0 {
		errs = append(errs, errors.New("at least one step or action is required"))
	}
	return errs
}
