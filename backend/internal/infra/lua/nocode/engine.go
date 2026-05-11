package nocode

// 零代码流程引擎

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"tiny-forum/internal/infra/lua/sdk"
)

// FlowEngine 执行一次零代码流程
type FlowEngine struct {
	api sdk.ForumAPI
}

func NewFlowEngine(api sdk.ForumAPI) *FlowEngine {
	return &FlowEngine{api: api}
}

// Run 执行 Flow，返回执行上下文（含日志）和错误
func (e *FlowEngine) Run(ctx context.Context, flow *Flow, event map[string]any) (*FlowContext, error) {
	fctx := NewFlowContext(event)
	fctx.Log(fmt.Sprintf("▶ 流程开始，触发器: %s", flow.Trigger.Type))

	// 1. 依次评估前置条件
	for i, cond := range flow.Conditions {
		ok, err := e.evalCondition(&cond, fctx)
		if err != nil {
			fctx.Log(fmt.Sprintf("✗ 条件[%d](%s) 评估出错: %v", i, cond.Type, err))
			return fctx, fmt.Errorf("condition[%d] %s: %w", i, cond.Type, err)
		}
		if cond.Negate {
			ok = !ok
		}
		if !ok {
			fctx.Log(fmt.Sprintf("⊘ 条件[%d](%s) 不满足，流程终止", i, cond.Type))
			return fctx, nil
		}
		fctx.Log(fmt.Sprintf("✓ 条件[%d](%s) 满足", i, cond.Type))
	}

	// 2. 顺序执行 Actions
	for i, action := range flow.Actions {
		fctx.Log(fmt.Sprintf("→ 动作[%d]: %s", i, action.Type))
		stop, err := e.execAction(ctx, &action, fctx)
		if err != nil {
			fctx.Log(fmt.Sprintf("✗ 动作[%d](%s) 失败: %v", i, action.Type, err))
			return fctx, fmt.Errorf("action[%d] %s: %w", i, action.Type, err)
		}
		fctx.Log(fmt.Sprintf("✓ 动作[%d](%s) 完成", i, action.Type))
		if stop {
			fctx.Log("⏹ stop_if 触发，流程提前结束")
			break
		}
	}

	fctx.Log("✅ 流程执行完毕")
	return fctx, nil
}

// ─── Condition 评估 ──────────────────────────────────────────────────────────

func (e *FlowEngine) evalCondition(cond *CondNode, fctx *FlowContext) (bool, error) {
	switch cond.Type {

	case CondPostTitleContains:
		title, _ := fctx.Get("post_title")
		return containsAny(fmt.Sprint(title), strSlice(cond.Params["keywords"])), nil

	case CondPostContentContains:
		content, _ := fctx.Get("post_content")
		return containsAny(fmt.Sprint(content), strSlice(cond.Params["keywords"])), nil

	case CondUserRoleIs:
		role, _ := fctx.Get("user_role")
		return fmt.Sprint(role) == fmt.Sprint(cond.Params["role"]), nil

	case CondUserPostCountGte:
		count, _ := fctx.Get("user_post_count")
		return toFloat64(count) >= toFloat64(cond.Params["count"]), nil

	case CondBoardIDIn:
		bid, _ := fctx.Get("board_id")
		for _, id := range float64Slice(cond.Params["ids"]) {
			if id == toFloat64(bid) {
				return true, nil
			}
		}
		return false, nil

	case CondTimeRange:
		start := strParam(cond.Params, "start") // "09:00"
		end := strParam(cond.Params, "end")
		tz := strParam(cond.Params, "tz")
		if tz == "" {
			tz = "Asia/Shanghai"
		}
		loc, _ := time.LoadLocation(tz)
		if loc == nil {
			loc = time.UTC
		}
		now := time.Now().In(loc).Format("15:04")
		return now >= start && now <= end, nil

	case CondCustomExpr:
		return evalExpr(strParam(cond.Params, "expr"), fctx)

	default:
		return false, fmt.Errorf("unknown condition type: %s", cond.Type)
	}
}

// ─── Action 执行 ─────────────────────────────────────────────────────────────

// execAction 返回 (stop, err)；stop=true 时终止后续动作
func (e *FlowEngine) execAction(ctx context.Context, action *ActionNode, fctx *FlowContext) (bool, error) {
	switch action.Type {

	// ── Post ──────────────────────────────────────────────────────────────
	case ActionReplyPost:
		postID := uintFromCtx(fctx, "post_id")
		content, err := render(strParam(action.Params, "content"), fctx)
		if err != nil {
			return false, err
		}
		_, err = e.api.ReplyPost(ctx, postID, content)
		return false, err

	case ActionDeletePost:
		return false, e.api.DeletePost(ctx, uintFromCtx(fctx, "post_id"))

	case ActionHidePost:
		return false, e.api.ModeratePost(ctx, uintFromCtx(fctx, "post_id"), "hide", "")

	case ActionPinPost:
		return false, e.api.ModeratePost(ctx, uintFromCtx(fctx, "post_id"), "pin", "")

	case ActionLockPost:
		return false, e.api.ModeratePost(ctx, uintFromCtx(fctx, "post_id"), "lock", "")

	case ActionCreatePost:
		boardID := uint(toFloat64(action.Params["board_id"]))
		title, _ := render(strParam(action.Params, "title"), fctx)
		content, _ := render(strParam(action.Params, "content"), fctx)
		_, err := e.api.CreatePost(ctx, sdk.CreatePostReq{
			Title:   title,
			Content: content,
			BoardID: boardID,
		})
		return false, err

	// ── Comment ───────────────────────────────────────────────────────────
	case ActionDeleteComment:
		return false, e.api.DeleteComment(ctx, uintFromCtx(fctx, "comment_id"))

	// ── User ──────────────────────────────────────────────────────────────
	case ActionBanUser:
		uid := uintFromCtx(fctx, "user_id")
		reason, _ := render(strParam(action.Params, "reason"), fctx)
		dur := int(toFloat64(action.Params["duration_sec"]))
		if dur <= 0 {
			dur = 86400
		}
		return false, e.api.BanUser(ctx, uid, reason, dur)

	case ActionSendMessage:
		// to_user_id 未填则默认发给触发者
		toUID := uint(toFloat64(action.Params["to_user_id"]))
		if toUID == 0 {
			toUID = uintFromCtx(fctx, "user_id")
		}
		content, err := render(strParam(action.Params, "content"), fctx)
		if err != nil {
			return false, err
		}
		return false, e.api.SendMessage(ctx, toUID, content)

	// ── Integration ───────────────────────────────────────────────────────
	case ActionWebhook:
		return false, e.execWebhook(ctx, action, fctx)

	case ActionNotifyAdmin:
		msg, _ := render(strParam(action.Params, "message"), fctx)
		// 实际可写入通知队列；此处记录日志
		fctx.Log("[NotifyAdmin] " + msg)
		return false, nil

	// ── Control ───────────────────────────────────────────────────────────
	case ActionWait:
		sec := int(toFloat64(action.Params["seconds"]))
		if sec > 30 {
			sec = 30
		}
		time.Sleep(time.Duration(sec) * time.Second)
		return false, nil

	case ActionSetVariable:
		name := strParam(action.Params, "name")
		val, err := render(strParam(action.Params, "value"), fctx)
		if err != nil {
			return false, err
		}
		fctx.Variables[name] = val
		return false, nil

	case ActionStopIf:
		ok, err := evalExpr(strParam(action.Params, "expr"), fctx)
		if err != nil {
			return false, err
		}
		return ok, nil // stop=true 时由调用方终止循环

	default:
		return false, fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// ─── Webhook 执行 ─────────────────────────────────────────────────────────────

func (e *FlowEngine) execWebhook(ctx context.Context, action *ActionNode, fctx *FlowContext) error {
	rawURL := strParam(action.Params, "url")
	method := strParam(action.Params, "method")
	if method == "" {
		method = http.MethodPost
	}
	bodyTpl := strParam(action.Params, "body")
	body, _ := render(bodyTpl, fctx)

	headers := map[string]string{"Content-Type": "application/json"}
	if hRaw := strParam(action.Params, "headers"); hRaw != "" {
		_ = json.Unmarshal([]byte(hRaw), &headers)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	req, err := http.NewRequestWithContext(reqCtx, method, rawURL, reqBody)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook %s returned %d: %s", rawURL, resp.StatusCode, string(b))
	}
	return nil
}

// ─── 模板渲染 ─────────────────────────────────────────────────────────────────
// 支持 Go text/template 语法，数据来自 event + variables 合并

func render(tpl string, fctx *FlowContext) (string, error) {
	if !strings.Contains(tpl, "{{") {
		return tpl, nil
	}
	data := make(map[string]any)
	for k, v := range fctx.Event {
		data[k] = v
	}
	for k, v := range fctx.Variables {
		data[k] = v
	}
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return tpl, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return tpl, err
	}
	return buf.String(), nil
}

// ─── 极简表达式解析 ───────────────────────────────────────────────────────────
// 仅支持 "<key> <op> <value>"，op: > < >= <= == !=

func evalExpr(expr string, fctx *FlowContext) (bool, error) {
	expr = strings.TrimSpace(expr)
	for _, op := range []string{">=", "<=", "!=", "==", ">", "<"} {
		idx := strings.Index(expr, op)
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(expr[:idx])
		rhs := strings.TrimSpace(expr[idx+len(op):])
		lv, _ := fctx.Get(key)
		lf := toFloat64(lv)
		rf, _ := strconv.ParseFloat(rhs, 64)
		switch op {
		case ">":
			return lf > rf, nil
		case "<":
			return lf < rf, nil
		case ">=":
			return lf >= rf, nil
		case "<=":
			return lf <= rf, nil
		case "==":
			return fmt.Sprint(lv) == rhs, nil
		case "!=":
			return fmt.Sprint(lv) != rhs, nil
		}
	}
	return false, fmt.Errorf("cannot parse expression: %q", expr)
}

// ─── 工具函数 ─────────────────────────────────────────────────────────────────

func containsAny(s string, keywords []string) bool {
	s = strings.ToLower(s)
	for _, kw := range keywords {
		if strings.Contains(s, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func strParam(params map[string]any, key string) string {
	if v, ok := params[key]; ok {
		return fmt.Sprint(v)
	}
	return ""
}

func uintFromCtx(fctx *FlowContext, key string) uint {
	v, _ := fctx.Get(key)
	return uint(toFloat64(v))
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}

func strSlice(v any) []string {
	switch val := v.(type) {
	case []string:
		return val
	case []interface{}:
		out := make([]string, 0, len(val))
		for _, item := range val {
			out = append(out, fmt.Sprint(item))
		}
		return out
	case string:
		if val == "" {
			return nil
		}
		return strings.Split(val, ",")
	}
	return nil
}

func float64Slice(v any) []float64 {
	if arr, ok := v.([]interface{}); ok {
		out := make([]float64, 0, len(arr))
		for _, item := range arr {
			out = append(out, toFloat64(item))
		}
		return out
	}
	return nil
}
