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
	"sync"
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

	// 新格式：统一 Steps
	if len(flow.Steps) > 0 {
		if err := e.execSteps(ctx, flow.Steps, fctx, 0); err != nil {
			return fctx, err
		}
		fctx.Log("✅ 流程执行完毕")
		return fctx, nil
	}

	// 兼容旧格式：Conditions → Tools → Actions
	// for i, cond := range flow.Conditions {
	// 	ok, err := e.evalCondition(&cond, fctx)
	// 	if err != nil {
	// 		return fctx, fmt.Errorf("condition[%d]: %w", i, err)
	// 	}
	// 	if cond.Negate {
	// 		ok = !ok
	// 	}
	// 	if !ok {
	// 		fctx.Log(fmt.Sprintf("⊘ 条件[%d] 不满足，流程终止", i))
	// 		return fctx, nil
	// 	}
	// }
	// for i := range flow.Tools {
	// 	if err := e.execFlowStep(ctx, &flow.Tools[i], fctx, "工具", i); err != nil {
	// 		return fctx, err
	// 	}
	// }
	for i := range flow.Actions {
		if err := e.execFlowStep(ctx, &flow.Actions[i], fctx, "动作", i); err != nil {
			return fctx, err
		}
	}
	fctx.Log("✅ 流程执行完毕")
	return fctx, nil
}

// execSteps 执行步骤序列，depth 跟踪嵌套深度
func (e *FlowEngine) execSteps(ctx context.Context, steps []FlowStep, fctx *FlowContext, depth int) error {
	if depth > MaxNestingDepth {
		return fmt.Errorf("超过最大嵌套深度 %d", MaxNestingDepth)
	}
	for i := range steps {
		if err := e.execFlowStep(ctx, &steps[i], fctx, "步骤", i); err != nil {
			return err
		}
	}
	return nil
}

// execFlowStep 执行单个步骤，支持分支、循环、深度检查
func (e *FlowEngine) execFlowStep(ctx context.Context, step *FlowStep, fctx *FlowContext, label string, index int) error {
	fctx.Depth++

	// if 分支
	if step.Type == "if" || step.Branch != nil {
		return e.execIf(ctx, step, fctx)
	}

	// while 循环
	if step.Type == "while" || step.Loop != nil {
		return e.execWhile(ctx, step, fctx)
	}

	// 终止
	if step.Type == "stop" {
		fctx.Log("⏹ 终止流程")
		return fmt.Errorf("flow stopped")
	}

	// 数据获取
	if step.Type == "get_post_info" || step.Type == "get_user_info" || step.Type == "get_comment_info" {
		return e.execGetData(step, fctx)
	}

	fctx.Log(fmt.Sprintf("→ %s[%d]: %s", label, index, step.Type))
	stop, err := e.execStepAction(ctx, step, fctx)
	if err != nil {
		fctx.Log(fmt.Sprintf("✗ %s[%d](%s) 失败: %v", label, index, step.Type, err))
		return fmt.Errorf("%s[%d] %s: %w", label, index, step.Type, err)
	}
	fctx.Log(fmt.Sprintf("✓ %s[%d](%s) 完成", label, index, step.Type))
	if stop {
		return fmt.Errorf("flow stopped by stop_if")
	}
	return nil
}

// execIf 执行条件分支
func (e *FlowEngine) execIf(ctx context.Context, step *FlowStep, fctx *FlowContext) error {
	var cond CondNode
	if step.Branch != nil {
		cond = step.Branch.Condition
	} else {
		field := strParam(step.Params, "condition_field")
		op := strParam(step.Params, "condition_op")
		val := strParam(step.Params, "condition_value")
		cond = buildCondition(field, op, val)
	}

	ok, err := e.evalCondition(&cond, fctx)
	if err != nil {
		return err
	}

	if ok {
		fctx.Log("  → if: TRUE")
		if step.Branch != nil {
			return e.execSteps(ctx, step.Branch.TrueSteps, fctx, fctx.Depth)
		}
	} else {
		fctx.Log("  → if: FALSE")
		if step.Branch != nil && len(step.Branch.FalseSteps) > 0 {
			return e.execSteps(ctx, step.Branch.FalseSteps, fctx, fctx.Depth)
		}
	}
	return nil
}

// execWhile 执行循环
func (e *FlowEngine) execWhile(ctx context.Context, step *FlowStep, fctx *FlowContext) error {
	var cond CondNode
	if step.Loop != nil {
		cond = step.Loop.Condition
		body := step.Loop.Body
		maxIter := step.Loop.MaxIter
		if maxIter <= 0 {
			maxIter = 100
		}
		for i := 0; i < maxIter; i++ {
			ok, err := e.evalCondition(&cond, fctx)
			if err != nil {
				return err
			}
			if !ok {
				fctx.Log(fmt.Sprintf("  → 循环退出（第 %d 次）", i+1))
				return nil
			}
			fctx.Log(fmt.Sprintf("  → 循环 #%d", i+1))
			if err := e.execSteps(ctx, body, fctx, fctx.Depth); err != nil {
				return err
			}
		}
		fctx.Log(fmt.Sprintf("  → 达到最大迭代 %d", maxIter))
		return nil
	}

	// 简化的 while：从 params 构建
	field := strParam(step.Params, "condition_field")
	op := strParam(step.Params, "condition_op")
	val := strParam(step.Params, "condition_value")
	cond = buildCondition(field, op, val)
	maxIter := int(toFloat64(step.Params["max_iter"]))
	if maxIter <= 0 {
		maxIter = 100
	}

	for i := 0; i < maxIter; i++ {
		ok, err := e.evalCondition(&cond, fctx)
		if err != nil {
			return err
		}
		if !ok {
			fctx.Log(fmt.Sprintf("  → 循环退出（第 %d 次）", i+1))
			return nil
		}
		fctx.Log(fmt.Sprintf("  → 循环 #%d", i+1))
		// body 通过画布上的 body 端口连接
	}
	return nil
}

// execGetData 从事件中提取数据到变量
func (e *FlowEngine) execGetData(step *FlowStep, fctx *FlowContext) error {
	switch step.Type {
	case "get_post_info":
		for k, v := range fctx.Event {
			if strings.HasPrefix(k, "post_") || k == "board_id" {
				fctx.Variables[k] = v
			}
		}
		fctx.Log("  → 已提取帖子信息到变量")
	case "get_user_info":
		for k, v := range fctx.Event {
			if strings.HasPrefix(k, "user_") || k == "username" || k == "author_id" || k == "user_role" {
				fctx.Variables[k] = v
			}
		}
		if uname, ok := fctx.Event["username"]; ok {
			fctx.Variables["username"] = uname
		}
		fctx.Log("  → 已提取用户信息到变量")
	case "get_comment_info":
		for k, v := range fctx.Event {
			if strings.HasPrefix(k, "comment_") {
				fctx.Variables[k] = v
			}
		}
		fctx.Log("  → 已提取评论信息到变量")
	}
	return nil
}

// buildCondition 从 field/op/value 构建 CondNode
func buildCondition(field, op, value string) CondNode {
	switch op {
	case "equals":
		return CondNode{Type: CondFieldEquals, Params: map[string]any{"field": field, "value": value}}
	case "not_equals":
		return CondNode{Type: CondFieldNotEquals, Params: map[string]any{"field": field, "value": value}}
	case "contains":
		return CondNode{Type: CondFieldContains, Params: map[string]any{"field": field, "value": value}}
	case "greater_than":
		return CondNode{Type: CondFieldGreaterThan, Params: map[string]any{"field": field, "value": value}}
	case "less_than":
		return CondNode{Type: CondFieldLessThan, Params: map[string]any{"field": field, "value": value}}
	case "is_empty":
		return CondNode{Type: CondFieldIsEmpty, Params: map[string]any{"field": field}}
	case "is_not_empty":
		return CondNode{Type: CondFieldNotEmpty, Params: map[string]any{"field": field}}
	default:
		return CondNode{Type: CondCustomExpr, Params: map[string]any{"expr": value}}
	}
}

// execParallel 并发执行多个步骤
func (e *FlowEngine) execParallel(ctx context.Context, steps []FlowStep, fctx *FlowContext) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i := range steps {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			fctx.Log(fmt.Sprintf("  ∥ 并行[%d]: %s 开始", idx, steps[idx].Type))
			_, err := e.execStepAction(ctx, &steps[idx], fctx)
			mu.Lock()
			if err != nil && firstErr == nil {
				firstErr = fmt.Errorf("parallel[%d] %s: %w", idx, steps[idx].Type, err)
			}
			mu.Unlock()
			if err != nil {
				fctx.Log(fmt.Sprintf("  ∥ 并行[%d]: %s 失败: %v", idx, steps[idx].Type, err))
			} else {
				fctx.Log(fmt.Sprintf("  ∥ 并行[%d]: %s 完成", idx, steps[idx].Type))
			}
		}(i)
	}
	wg.Wait()
	return firstErr
}

// ─── 通用字段比较辅助方法 ──────────────────────────────────────────────────────────

func (e *FlowEngine) getField(cond *CondNode, fctx *FlowContext) string {
	field := strParam(cond.Params, "field")
	val, _ := fctx.Get(field)
	return fmt.Sprint(val)
}

func (e *FlowEngine) fieldCmp(cond *CondNode, fctx *FlowContext, cmp func(string, string) bool) (bool, error) {
	fv := e.getField(cond, fctx)
	expected := strParam(cond.Params, "value")
	rendered, err := render(expected, fctx)
	if err != nil {
		return false, err
	}
	return cmp(fv, rendered), nil
}

func (e *FlowEngine) fieldCmpNum(cond *CondNode, fctx *FlowContext, cmp func(float64, float64) bool) (bool, error) {
	fv := e.getField(cond, fctx)
	expected := strParam(cond.Params, "value")
	rendered, err := render(expected, fctx)
	if err != nil {
		return false, err
	}
	return cmp(toFloat64(fv), toFloat64(rendered)), nil
}

func (e *FlowEngine) fieldContains(cond *CondNode, fctx *FlowContext, negate bool) (bool, error) {
	fv := strings.ToLower(e.getField(cond, fctx))
	keywords := strSlice(cond.Params["value"])
	match := false
	for _, kw := range keywords {
		if strings.Contains(fv, strings.ToLower(kw)) {
			match = true
			break
		}
	}
	if negate {
		return !match, nil
	}
	return match, nil
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

	// ── 通用字段比较 ──────────────────────────────────────────────────
	case CondFieldEquals:
		return e.fieldCmp(cond, fctx, func(a, b string) bool { return strings.EqualFold(a, b) })
	case CondFieldNotEquals:
		return e.fieldCmp(cond, fctx, func(a, b string) bool { return !strings.EqualFold(a, b) })
	case CondFieldContains:
		return e.fieldContains(cond, fctx, false)
	case CondFieldNotContains:
		return e.fieldContains(cond, fctx, true)
	case CondFieldGreaterThan:
		return e.fieldCmpNum(cond, fctx, func(a, b float64) bool { return a > b })
	case CondFieldLessThan:
		return e.fieldCmpNum(cond, fctx, func(a, b float64) bool { return a < b })
	case CondFieldIsEmpty:
		fv := e.getField(cond, fctx)
		return fv == "" || fv == "0", nil
	case CondFieldNotEmpty:
		fv := e.getField(cond, fctx)
		return fv != "" && fv != "0", nil

	default:
		return false, fmt.Errorf("unknown condition type: %s", cond.Type)
	}
}

// ─── Action 执行 ─────────────────────────────────────────────────────────────

// execStepAction 返回 (stop, err)；stop=true 时终止后续步骤
func (e *FlowEngine) execStepAction(ctx context.Context, step *FlowStep, fctx *FlowContext) (bool, error) {
	switch ActionType(step.Type) {

	// ── Post ──────────────────────────────────────────────────────────────
	case ActionReplyPost:
		postID := uintFromCtx(fctx, "post_id")
		content, err := render(strParam(step.Params, "content"), fctx)
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
		boardID := uint(toFloat64(step.Params["board_id"]))
		title, _ := render(strParam(step.Params, "title"), fctx)
		content, _ := render(strParam(step.Params, "content"), fctx)
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
		reason, _ := render(strParam(step.Params, "reason"), fctx)
		dur := int(toFloat64(step.Params["duration_sec"]))
		if dur <= 0 {
			dur = 86400
		}
		return false, e.api.BanUser(ctx, uid, reason, dur)

	case ActionSendMessage:
		// to_user_id 未填则默认发给触发者
		toUID := uint(toFloat64(step.Params["to_user_id"]))
		if toUID == 0 {
			toUID = uintFromCtx(fctx, "user_id")
		}
		content, err := render(strParam(step.Params, "content"), fctx)
		if err != nil {
			return false, err
		}
		return false, e.api.SendMessage(ctx, toUID, content)

	// ── Integration ───────────────────────────────────────────────────────
	case ActionWebhook:
		return false, e.execWebhook(ctx, step, fctx)

	case ActionNotifyAdmin:
		msg, _ := render(strParam(step.Params, "message"), fctx)
		// 实际可写入通知队列；此处记录日志
		fctx.Log("[NotifyAdmin] " + msg)
		return false, nil

	// ── Control ───────────────────────────────────────────────────────────
	case ActionWait:
		sec := int(toFloat64(step.Params["seconds"]))
		if sec > 30 {
			sec = 30
		}
		time.Sleep(time.Duration(sec) * time.Second)
		return false, nil

	case ActionSetVariable:
		name := strParam(step.Params, "name")
		val, err := render(strParam(step.Params, "value"), fctx)
		if err != nil {
			return false, err
		}
		fctx.Variables[name] = val
		return false, nil

	case ActionStopIf:
		ok, err := evalExpr(strParam(step.Params, "expr"), fctx)
		if err != nil {
			return false, err
		}
		return ok, nil

	// ── Math ────────────────────────────────────────────────────────────────
	case ActionType(ActionAdd):
		e.mathOp(step.Params, fctx, func(a, b float64) float64 { return a + b })
		return false, nil
	case ActionType(ActionSubtract):
		e.mathOp(step.Params, fctx, func(a, b float64) float64 { return a - b })
		return false, nil
	case ActionType(ActionMultiply):
		e.mathOp(step.Params, fctx, func(a, b float64) float64 { return a * b })
		return false, nil
	case ActionType(ActionDivide):
		e.mathOp(step.Params, fctx, func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		})
		return false, nil
	case ActionType(ActionModulo):
		e.mathOp(step.Params, fctx, func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return float64(int64(a) % int64(b))
		})
		return false, nil

	// ── String ──────────────────────────────────────────────────────────────
	case ActionType(ActionConcat):
		target := strParam(step.Params, "target")
		a, _ := render(strParam(step.Params, "a"), fctx)
		b, _ := render(strParam(step.Params, "b"), fctx)
		fctx.Variables[target] = a + b
		return false, nil
	case ActionType(ActionLength):
		target := strParam(step.Params, "target")
		src := strParam(step.Params, "source")
		val, _ := render(src, fctx)
		fctx.Variables[target] = len(val)
		return false, nil

	default:
		// 未知类型（如遗留的 branch 节点）：静默跳过
		fctx.Log(fmt.Sprintf("⊘ 跳过未知类型: %s", step.Type))
		return false, nil
	}
}

// ─── Math helper ──────────────────────────────────────────────────────────────

func (e *FlowEngine) mathOp(params map[string]any, fctx *FlowContext, op func(float64, float64) float64) {
	target := strParam(params, "target")
	aStr, _ := render(strParam(params, "a"), fctx)
	bStr, _ := render(strParam(params, "b"), fctx)
	a := parseNumber(aStr)
	if v, ok := fctx.Variables[target]; ok {
		a = toFloat64(v)
	}
	b := parseNumber(bStr)
	result := op(a, b)
	fctx.Variables[target] = result
	fctx.Log(fmt.Sprintf("  math: %s = %.2f", target, result))
}

func parseNumber(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

// ─── Webhook 执行 ─────────────────────────────────────────────────────────────

func (e *FlowEngine) execWebhook(ctx context.Context, step *FlowStep, fctx *FlowContext) error {
	rawURL := strParam(step.Params, "url")
	method := strParam(step.Params, "method")
	if method == "" {
		method = http.MethodPost
	}
	bodyTpl := strParam(step.Params, "body")
	body, _ := render(bodyTpl, fctx)

	headers := map[string]string{"Content-Type": "application/json"}
	if hRaw := strParam(step.Params, "headers"); hRaw != "" {
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
