// Package nocode 定义零代码机器人的流程节点模型。
//
// 四类节点：
//
//	触发器 → 控制 → 变量 → 动作
//
// Flow 序列化后存储在 bot.config_values["flow"] 中（JSON 字符串）。
package nocode

import "encoding/json"

// ─── Flow ────────────────────────────────────────────────────────────────────

type Flow struct {
	Version string      `json:"version"`
	Trigger TriggerNode `json:"trigger"`
	Steps   []FlowStep  `json:"steps"` // 统一步骤（控制/变量/动作）
	Actions []FlowStep  `json:"actions,omitempty"`
}

func FlowToJSON(f *Flow) (string, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

func FlowFromJSON(s string) (*Flow, error) {
	var f Flow
	return &f, json.Unmarshal([]byte(s), &f)
}

// ─── Trigger ─────────────────────────────────────────────────────────────────

type TriggerType string

const (
	TriggerOnSchedule     TriggerType = "on_schedule"
	TriggerOnNewPost      TriggerType = "on_new_post"
	TriggerOnNewComment   TriggerType = "on_new_comment"
	TriggerOnUserRegister TriggerType = "on_user_register"
	TriggerOnKeyword      TriggerType = "on_keyword"
	TriggerOnManual       TriggerType = "on_manual"
)

type TriggerNode struct {
	Type   TriggerType    `json:"type"`
	Params map[string]any `json:"params,omitempty"`
}

// ─── Condition ───────────────────────────────────────────────────────────────

type CondType string

const (
	CondPostTitleContains   CondType = "post_title_contains"
	CondPostContentContains CondType = "post_content_contains"
	CondUserRoleIs          CondType = "user_role_is"
	CondUserPostCountGte    CondType = "user_post_count_gte"
	CondBoardIDIn           CondType = "board_id_in"
	CondTimeRange           CondType = "time_range"
	CondCustomExpr          CondType = "custom_expr"
	CondFieldEquals         CondType = "field_equals"
	CondFieldNotEquals      CondType = "field_not_equals"
	CondFieldContains       CondType = "field_contains"
	CondFieldNotContains    CondType = "field_not_contains"
	CondFieldGreaterThan    CondType = "field_greater_than"
	CondFieldLessThan       CondType = "field_less_than"
	CondFieldIsEmpty        CondType = "field_is_empty"
	CondFieldNotEmpty       CondType = "field_is_not_empty"
)

type CondNode struct {
	Type   CondType       `json:"type"`
	Negate bool           `json:"negate,omitempty"`
	Params map[string]any `json:"params"`
}

// ─── FlowStep ────────────────────────────────────────────────────────────────

type VarOutput struct {
	Name string `json:"name"`
	Type string `json:"type"` // string | number | boolean | object
	Desc string `json:"desc,omitempty"`
}

type FlowStep struct {
	ID      string         `json:"id,omitempty"`
	Type    string         `json:"type"`
	Label   string         `json:"label,omitempty"`
	Params  map[string]any `json:"params,omitempty"`
	Branch  *BranchConfig  `json:"branch,omitempty"`
	Loop    *LoopConfig    `json:"loop,omitempty"`
	Outputs []VarOutput    `json:"outputs,omitempty"` // 该节点产生的变量
}

type BranchConfig struct {
	Condition  CondNode   `json:"condition"`
	TrueSteps  []FlowStep `json:"true"`
	FalseSteps []FlowStep `json:"false,omitempty"`
}

type LoopConfig struct {
	Condition CondNode   `json:"condition"`
	Body      []FlowStep `json:"body"`
	MaxIter   int        `json:"max_iter,omitempty"`
}

// ─── ActionType（向后兼容）────────────────────────────────────────────────

type ActionType string

const (
	ActionReplyPost     ActionType = "reply_post"
	ActionDeletePost    ActionType = "delete_post"
	ActionHidePost      ActionType = "hide_post"
	ActionPinPost       ActionType = "pin_post"
	ActionLockPost      ActionType = "lock_post"
	ActionCreatePost    ActionType = "create_post"
	ActionDeleteComment ActionType = "delete_comment"
	ActionBanUser       ActionType = "ban_user"
	ActionSendMessage   ActionType = "send_message"
	ActionWebhook       ActionType = "webhook"
	ActionNotifyAdmin   ActionType = "notify_admin"
	ActionWait          ActionType = "wait"
	ActionStopIf        ActionType = "stop_if"
	// 变量/数据
	ActionSetVariable    ActionType = "set_variable"
	ActionGetPostInfo    ActionType = "get_post_info"
	ActionGetUserInfo    ActionType = "get_user_info"
	ActionGetCommentInfo ActionType = "get_comment_info"
	// 控制
	ActionBranch ActionType = "branch"
	ActionLoop   ActionType = "loop"
	// 算术
	ActionAdd      ActionType = "add"
	ActionSubtract ActionType = "subtract"
	ActionMultiply ActionType = "multiply"
	ActionDivide   ActionType = "divide"
	ActionModulo   ActionType = "modulo"
	ActionConcat   ActionType = "concat"
	ActionLength   ActionType = "length"
)

// ─── FlowContext ──────────────────────────────────────────────────────────────

type FlowContext struct {
	Event     map[string]any
	Variables map[string]any
	Logs      []string
	Depth     int // 当前嵌套深度
}

func NewFlowContext(event map[string]any) *FlowContext {
	if event == nil {
		event = map[string]any{}
	}
	return &FlowContext{
		Event:     event,
		Variables: make(map[string]any),
	}
}

func (c *FlowContext) Log(msg string) {
	c.Logs = append(c.Logs, msg)
}

func (c *FlowContext) Get(key string) (any, bool) {
	if v, ok := c.Variables[key]; ok {
		return v, true
	}
	v, ok := c.Event[key]
	return v, ok
}

const MaxNestingDepth = 5
