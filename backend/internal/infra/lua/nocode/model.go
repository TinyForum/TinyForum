// Package nocode 定义零代码机器人的流程节点模型。
//
// 用户通过前端拖拽，将节点组合成如下线性流程：
//
//	Trigger → [Condition…] → [Action…]
//
// Flow 序列化后存储在 bot.config_values["flow"] 中（JSON 字符串）。
package nocode

import "encoding/json"

// ─── Flow ────────────────────────────────────────────────────────────────────

// Flow 是一个零代码机器人的完整流程描述。
type Flow struct {
	Version    string       `json:"version"`              // 目前固定 "1"
	Trigger    TriggerNode  `json:"trigger"`              // 触发器（唯一）
	Conditions []CondNode   `json:"conditions,omitempty"` // 前置条件，全部满足才执行 Actions
	Actions    []ActionNode `json:"actions"`              // 顺序执行的动作列表
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

// TriggerNode 描述触发方式和参数
type TriggerNode struct {
	Type   TriggerType    `json:"type"`
	Params map[string]any `json:"params,omitempty"`
	// on_schedule  → { "cron": "0 9 * * 1" }
	// on_keyword   → { "keywords": ["广告"], "scope": "post|comment|both" }
	// on_new_post  → { "board_ids": [1,2] }  // 空=全部板块
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
)

// CondNode 单个条件
type CondNode struct {
	Type   CondType       `json:"type"`
	Negate bool           `json:"negate,omitempty"` // true → NOT 取反
	Params map[string]any `json:"params"`
}

// ─── Action ──────────────────────────────────────────────────────────────────

type ActionType string

const (
	// Post
	ActionReplyPost  ActionType = "reply_post"
	ActionDeletePost ActionType = "delete_post"
	ActionHidePost   ActionType = "hide_post"
	ActionPinPost    ActionType = "pin_post"
	ActionLockPost   ActionType = "lock_post"
	ActionCreatePost ActionType = "create_post"
	// Comment
	ActionDeleteComment ActionType = "delete_comment"
	// User
	ActionBanUser     ActionType = "ban_user"
	ActionSendMessage ActionType = "send_message"
	// Integration
	ActionWebhook     ActionType = "webhook"
	ActionNotifyAdmin ActionType = "notify_admin"
	// Control
	ActionWait        ActionType = "wait"
	ActionSetVariable ActionType = "set_variable"
	ActionStopIf      ActionType = "stop_if"
)

// ActionNode 单个动作
type ActionNode struct {
	Type   ActionType     `json:"type"`
	Params map[string]any `json:"params"`
	// reply_post   → { "content": "感谢 {{username}} 发帖！" }
	// ban_user     → { "reason": "违规", "duration_sec": 86400 }
	// webhook      → { "url": "https://...", "method": "POST", "body": "..." }
	// wait         → { "seconds": 3 }
	// set_variable → { "name": "score", "value": "{{event.score}}" }
}

// ─── FlowContext ──────────────────────────────────────────────────────────────

// FlowContext 在一次流程执行中传递，供条件和动作读写变量。
type FlowContext struct {
	Event     map[string]any // 触发事件原始数据
	Variables map[string]any // set_variable 写入的变量
	Logs      []string
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

// Get 从 Variables（优先）或 Event 中读取值
func (c *FlowContext) Get(key string) (any, bool) {
	if v, ok := c.Variables[key]; ok {
		return v, true
	}
	v, ok := c.Event[key]
	return v, ok
}
