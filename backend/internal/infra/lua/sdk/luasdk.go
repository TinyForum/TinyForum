// Package sdk 向 Lua 沙箱注入论坛 SDK，供机器人脚本调用。
// 调用顺序：NewBotSDK → InjectAll → sandbox.Execute
package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"tiny-forum/internal/model/do"

	lua "github.com/yuin/gopher-lua"
)

// ─── 论坛 API 接口（由 service/bot 层实现，注入到 SDK）────────────────────

// ForumAPI 是业务能力的最小接口集合，SDK 通过它调用论坛功能。
// 实现者需自行进行必要的鉴权（操作者固定为机器人身份）。
type ForumAPI interface {
	// Post
	GetPost(ctx context.Context, postID uint) (*PostVO, error)
	CreatePost(ctx context.Context, req CreatePostReq) (*PostVO, error)
	ReplyPost(ctx context.Context, postID uint, content string) (*CommentVO, error)
	DeletePost(ctx context.Context, postID uint) error
	ModeratePost(ctx context.Context, postID uint, action string, reason string) error

	// Comment
	GetComment(ctx context.Context, commentID uint) (*CommentVO, error)
	DeleteComment(ctx context.Context, commentID uint) error

	// User
	GetUser(ctx context.Context, userID uint) (*UserVO, error)
	BanUser(ctx context.Context, userID uint, reason string, durationSec int) error
	SendMessage(ctx context.Context, toUserID uint, content string) error

	// Stats
	GetForumStats(ctx context.Context) (*StatsVO, error)
}

// ─── VO（SDK 内部用，与 do 层解耦）──────────────────────────────────────────

type PostVO struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorID  uint   `json:"author_id"`
	BoardID   uint   `json:"board_id"`
	CreatedAt int64  `json:"created_at"`
}

type CommentVO struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	AuthorID  uint   `json:"author_id"`
	PostID    uint   `json:"post_id"`
	CreatedAt int64  `json:"created_at"`
}

type UserVO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type StatsVO struct {
	PostCount    int64 `json:"post_count"`
	UserCount    int64 `json:"user_count"`
	CommentCount int64 `json:"comment_count"`
	ActiveToday  int64 `json:"active_today"`
}

type CreatePostReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	BoardID uint   `json:"board_id"`
}

// ─── BotSDK ──────────────────────────────────────────────────────────────────

// BotSDK 持有 bot 实例（用于权限判断）和 ForumAPI 实现，向 Lua 注入所有 API。
type BotSDK struct {
	bot            *do.Bot
	api            ForumAPI
	allowedDomains []string
	logs           []string
}

func NewBotSDK(bot *do.Bot, api ForumAPI, allowedDomains []string) *BotSDK {
	return &BotSDK{bot: bot, api: api, allowedDomains: allowedDomains}
}

// Logs 返回本次执行收集到的日志行
func (s *BotSDK) Logs() []string { return s.logs }

func (s *BotSDK) addLog(msg string) {
	line := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg)
	s.logs = append(s.logs, line)
	fmt.Println("[BotLua]", line)
}

// InjectAll 将全部 SDK 注入 Lua 状态机，并注入 event 全局表
func (s *BotSDK) InjectAll(L *lua.LState, ctx context.Context, event map[string]any) {
	s.injectLog(L)
	s.injectConfig(L)
	s.injectEvent(L, event)
	s.injectForum(L, ctx)
	s.injectHTTP(L)
	s.injectUtil(L)
	s.injectJSON(L)
}

// ─── log ─────────────────────────────────────────────────────────────────────

func (s *BotSDK) injectLog(L *lua.LState) {
	L.SetGlobal("log", L.NewFunction(func(L *lua.LState) int {
		s.addLog(L.CheckString(1))
		return 0
	}))
	L.SetGlobal("logf", L.NewFunction(func(L *lua.LState) int {
		format := L.CheckString(1)
		args := make([]interface{}, 0, L.GetTop()-1)
		for i := 2; i <= L.GetTop(); i++ {
			args = append(args, luaToGo(L.Get(i)))
		}
		s.addLog(fmt.Sprintf(format, args...))
		return 0
	}))
}

// ─── config / event ──────────────────────────────────────────────────────────

func (s *BotSDK) injectConfig(L *lua.LState) {
	L.SetGlobal("config", goMapToLuaTable(L, s.bot.ConfigValues))
}

func (s *BotSDK) injectEvent(L *lua.LState, event map[string]any) {
	L.SetGlobal("event", goMapToLuaTable(L, event))
}

// ─── forum.* ─────────────────────────────────────────────────────────────────

func (s *BotSDK) injectForum(L *lua.LState, ctx context.Context) {
	tbl := L.NewTable()

	// forum.getPost(post_id) → post_table | nil, err_str
	tbl.RawSetString("getPost", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.BotPermReadPosts) {
			return permDenied(L, "read:posts")
		}
		id := uint(L.CheckInt64(1))
		post, err := s.api.GetPost(ctx, id)
		if err != nil {
			return nilErr(L, err)
		}
		L.Push(structToTable(L, post))
		return 1
	}))

	// forum.createPost(title, content, board_id) → post_table | nil, err_str
	tbl.RawSetString("createPost", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.BotPermWritePosts) {
			return permDenied(L, "write:posts")
		}
		post, err := s.api.CreatePost(ctx, CreatePostReq{
			Title:   L.CheckString(1),
			Content: L.CheckString(2),
			BoardID: uint(L.CheckInt64(3)),
		})
		if err != nil {
			return nilErr(L, err)
		}
		L.Push(structToTable(L, post))
		return 1
	}))

	// forum.replyPost(post_id, content) → comment_table | nil, err_str
	tbl.RawSetString("replyPost", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.BotPermWriteComments) {
			return permDenied(L, "write:comments")
		}
		comment, err := s.api.ReplyPost(ctx, uint(L.CheckInt64(1)), L.CheckString(2))
		if err != nil {
			return nilErr(L, err)
		}
		L.Push(structToTable(L, comment))
		return 1
	}))

	// forum.deletePost(post_id) → true | false, err_str
	tbl.RawSetString("deletePost", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.PermManageContent) {
			return permDenied(L, "manage:content")
		}
		return boolResult(L, s.api.DeletePost(ctx, uint(L.CheckInt64(1))))
	}))

	// forum.moderatePost(post_id, action, reason) → true | false, err_str
	// action: "hide" | "pin" | "lock" | "delete"
	tbl.RawSetString("moderatePost", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.PermManageContent) {
			return permDenied(L, "manage:content")
		}
		return boolResult(L, s.api.ModeratePost(ctx,
			uint(L.CheckInt64(1)), L.CheckString(2), L.OptString(3, "")))
	}))

	// forum.getComment(comment_id) → comment_table | nil, err_str
	tbl.RawSetString("getComment", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.BotPermReadComments) {
			return permDenied(L, "read:comments")
		}
		c, err := s.api.GetComment(ctx, uint(L.CheckInt64(1)))
		if err != nil {
			return nilErr(L, err)
		}
		L.Push(structToTable(L, c))
		return 1
	}))

	// forum.deleteComment(comment_id) → true | false, err_str
	tbl.RawSetString("deleteComment", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.PermManageContent) {
			return permDenied(L, "manage:content")
		}
		return boolResult(L, s.api.DeleteComment(ctx, uint(L.CheckInt64(1))))
	}))

	// forum.getUser(user_id) → user_table | nil, err_str
	tbl.RawSetString("getUser", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.BotPermReadUser) {
			return permDenied(L, "read:user")
		}
		u, err := s.api.GetUser(ctx, uint(L.CheckInt64(1)))
		if err != nil {
			return nilErr(L, err)
		}
		L.Push(structToTable(L, u))
		return 1
	}))

	// forum.banUser(user_id, reason, duration_sec) → true | false, err_str
	tbl.RawSetString("banUser", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.PermManageContent) {
			return permDenied(L, "manage:content")
		}
		dur := L.OptInt(3, 86400)
		return boolResult(L, s.api.BanUser(ctx, uint(L.CheckInt64(1)), L.CheckString(2), dur))
	}))

	// forum.sendMessage(to_user_id, content) → true | false, err_str
	tbl.RawSetString("sendMessage", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.PermSendMessage) {
			return permDenied(L, "send:message")
		}
		return boolResult(L, s.api.SendMessage(ctx, uint(L.CheckInt64(1)), L.CheckString(2)))
	}))

	// forum.getStats() → stats_table | nil, err_str
	tbl.RawSetString("getStats", L.NewFunction(func(L *lua.LState) int {
		if !s.hasPerm(do.PermReadStats) {
			return permDenied(L, "read:stats")
		}
		stats, err := s.api.GetForumStats(ctx)
		if err != nil {
			return nilErr(L, err)
		}
		L.Push(structToTable(L, stats))
		return 1
	}))

	L.SetGlobal("forum", tbl)
}

// ─── http.* ──────────────────────────────────────────────────────────────────

func (s *BotSDK) injectHTTP(L *lua.LState) {
	tbl := L.NewTable()

	doReq := func(L *lua.LState, method string) int {
		rawURL := L.CheckString(1)
		if !s.domainAllowed(rawURL) {
			L.Push(lua.LNil)
			L.Push(lua.LString("domain not allowed: " + rawURL))
			return 2
		}
		body := L.OptString(2, "")
		headers := map[string]string{}
		if L.GetTop() >= 3 {
			if t, ok := L.Get(3).(*lua.LTable); ok {
				t.ForEach(func(k, v lua.LValue) { headers[k.String()] = v.String() })
			}
		}
		respBody, status, err := execHTTP(method, rawURL, headers, body)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LString(respBody))
		L.Push(lua.LNumber(status))
		return 2
	}

	// http.get(url, headers?) → body, status | nil, err
	tbl.RawSetString("get", L.NewFunction(func(L *lua.LState) int { return doReq(L, http.MethodGet) }))
	// http.post(url, body, headers?) → body, status | nil, err
	tbl.RawSetString("post", L.NewFunction(func(L *lua.LState) int { return doReq(L, http.MethodPost) }))

	L.SetGlobal("http", tbl)
}

// ─── util.* ──────────────────────────────────────────────────────────────────

func (s *BotSDK) injectUtil(L *lua.LState) {
	tbl := L.NewTable()

	tbl.RawSetString("now", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(time.Now().Unix()))
		return 1
	}))
	tbl.RawSetString("format_time", L.NewFunction(func(L *lua.LState) int {
		ts := L.CheckInt64(1)
		layout := L.OptString(2, "2006-01-02 15:04:05")
		L.Push(lua.LString(time.Unix(ts, 0).Format(layout)))
		return 1
	}))
	tbl.RawSetString("sleep", L.NewFunction(func(L *lua.LState) int {
		ms := L.CheckInt(1)
		if ms > 5000 {
			ms = 5000
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return 0
	}))
	tbl.RawSetString("contains", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LBool(strings.Contains(L.CheckString(1), L.CheckString(2))))
		return 1
	}))
	tbl.RawSetString("split", L.NewFunction(func(L *lua.LState) int {
		parts := strings.Split(L.CheckString(1), L.CheckString(2))
		arr := L.CreateTable(len(parts), 0)
		for i, p := range parts {
			arr.RawSetInt(i+1, lua.LString(p))
		}
		L.Push(arr)
		return 1
	}))
	tbl.RawSetString("trim", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString(strings.TrimSpace(L.CheckString(1))))
		return 1
	}))
	tbl.RawSetString("lower", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString(strings.ToLower(L.CheckString(1))))
		return 1
	}))
	tbl.RawSetString("upper", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString(strings.ToUpper(L.CheckString(1))))
		return 1
	}))

	L.SetGlobal("util", tbl)
}

// ─── json.* ──────────────────────────────────────────────────────────────────

func (s *BotSDK) injectJSON(L *lua.LState) {
	tbl := L.NewTable()

	tbl.RawSetString("encode", L.NewFunction(func(L *lua.LState) int {
		b, err := json.Marshal(luaToGo(L.Get(1)))
		if err != nil {
			return nilErr(L, err)
		}
		L.Push(lua.LString(string(b)))
		return 1
	}))
	tbl.RawSetString("decode", L.NewFunction(func(L *lua.LState) int {
		var out interface{}
		if err := json.Unmarshal([]byte(L.CheckString(1)), &out); err != nil {
			return nilErr(L, err)
		}
		L.Push(goToLua(L, out))
		return 1
	}))

	L.SetGlobal("json", tbl)
}

// ─── 权限/辅助 ────────────────────────────────────────────────────────────────

func (s *BotSDK) hasPerm(p do.BotPermission) bool {
	for _, perm := range s.bot.Permissions {
		if perm == p {
			return true
		}
	}
	return false
}

func (s *BotSDK) domainAllowed(rawURL string) bool {
	for _, d := range s.allowedDomains {
		if strings.Contains(rawURL, d) {
			return true
		}
	}
	return len(s.allowedDomains) == 0 // 未配置白名单则全放行（开发环境）
}

// ─── 共享类型转换工具 ─────────────────────────────────────────────────────────

func permDenied(L *lua.LState, perm string) int {
	L.Push(lua.LNil)
	L.Push(lua.LString("permission denied: " + perm))
	return 2
}

func nilErr(L *lua.LState, err error) int {
	L.Push(lua.LNil)
	L.Push(lua.LString(err.Error()))
	return 2
}

func boolResult(L *lua.LState, err error) int {
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LTrue)
	return 1
}

func structToTable(L *lua.LState, v interface{}) *lua.LTable {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return goMapToLuaTable(L, m)
}

func goMapToLuaTable(L *lua.LState, m map[string]any) *lua.LTable {
	tbl := L.NewTable()
	for k, v := range m {
		tbl.RawSetString(k, goToLua(L, v))
	}
	return tbl
}

func goToLua(L *lua.LState, v interface{}) lua.LValue {
	switch val := v.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(val)
	case float64:
		return lua.LNumber(val)
	case int:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case uint:
		return lua.LNumber(val)
	case string:
		return lua.LString(val)
	case map[string]interface{}:
		return goMapToLuaTable(L, val)
	case []interface{}:
		arr := L.CreateTable(len(val), 0)
		for i, item := range val {
			arr.RawSetInt(i+1, goToLua(L, item))
		}
		return arr
	default:
		return lua.LNil
	}
}

func luaToGo(v lua.LValue) interface{} {
	switch val := v.(type) {
	case lua.LBool:
		return bool(val)
	case lua.LNumber:
		return float64(val)
	case lua.LString:
		return string(val)
	case *lua.LTable:
		m := make(map[string]interface{})
		val.ForEach(func(k, v lua.LValue) { m[k.String()] = luaToGo(v) })
		return m
	default:
		return nil
	}
}

func execHTTP(method, rawURL string, headers map[string]string, body string) (string, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, reqBody)
	if err != nil {
		return "", 0, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b), resp.StatusCode, nil
}