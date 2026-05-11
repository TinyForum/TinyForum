package luaengine

// Package lua 提供安全的 Lua 脚本执行沙箱。

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tiny-forum/internal/infra/lua/sdk"
	"tiny-forum/internal/model/do"

	lua "github.com/yuin/gopher-lua"
)

// ExecResult 单次执行结果
type ExecResult struct {
	Output   interface{}
	Logs     []string
	Duration time.Duration
	Err      error
}

// LuaSandbox 无状态沙箱，每次 Execute 创建新 LState。
type LuaSandbox struct {
	defaultTimeoutSec int
	allowedDomains    []string
}

func NewLuaSandbox(defaultTimeoutSec int, allowedDomains []string) *LuaSandbox {
	return &LuaSandbox{
		defaultTimeoutSec: defaultTimeoutSec,
		allowedDomains:    allowedDomains,
	}
}

// Execute 执行机器人 Lua 脚本。
//   - ctx：上层 context，传递给 ForumAPI 调用
//   - bot：机器人实例，提供权限和超时配置
//   - api：论坛能力实现
//   - event：触发事件数据
func (s *LuaSandbox) Execute(
	ctx context.Context,
	bot *do.Bot,
	api sdk.ForumAPI,
	event map[string]any,
) *ExecResult {
	start := time.Now()
	botSDK := sdk.NewBotSDK(bot, api, s.allowedDomains)

	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()

	s.openSafeLibs(L)
	botSDK.InjectAll(L, ctx, event)

	timeoutSec := bot.TimeoutSec
	if timeoutSec <= 0 {
		timeoutSec = s.defaultTimeoutSec
	}
	execCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	type result struct {
		output interface{}
		err    error
	}
	ch := make(chan result, 1)

	go func() {
		// 加载脚本
		if err := L.DoString(bot.ScriptCode); err != nil {
			ch <- result{err: err}
			return
		}
		// 调用 main()
		fn := L.GetGlobal("main")
		if fn.Type() != lua.LTFunction {
			ch <- result{err: errors.New("main() function not found in script")}
			return
		}
		if err := L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
			ch <- result{err: err}
			return
		}
		ret := L.Get(-1)
		L.Pop(1)
		ch <- result{output: luaToGo(ret)}
	}()

	select {
	case res := <-ch:
		return &ExecResult{
			Output:   res.output,
			Logs:     botSDK.Logs(),
			Duration: time.Since(start),
			Err:      res.err,
		}
	case <-execCtx.Done():
		L.Close() // 强制终止 goroutine 中的 LState
		return &ExecResult{
			Logs:     botSDK.Logs(),
			Duration: time.Since(start),
			Err:      fmt.Errorf("execution timeout after %ds", timeoutSec),
		}
	}
}

// openSafeLibs 加载安全库，禁用危险函数
func (s *LuaSandbox) openSafeLibs(L *lua.LState) {
	lua.OpenBase(L)
	lua.OpenTable(L)
	lua.OpenString(L)
	lua.OpenMath(L)

	// 移除危险全局函数
	for _, name := range []string{
		"dofile", "loadfile", "load", "loadstring", "require",
		"rawget", "rawset", "rawequal", "rawlen",
		"debug", "coroutine", "package", "io",
	} {
		L.SetGlobal(name, lua.LNil)
	}

	// 替换 os 为只读安全子集
	osTable := L.NewTable()
	osTable.RawSetString("time", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(time.Now().Unix()))
		return 1
	}))
	osTable.RawSetString("clock", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(float64(time.Now().UnixMicro()) / 1e6))
		return 1
	}))
	L.SetGlobal("os", osTable)

	// 移除 string.dump
	if str := L.GetGlobal("string"); str.Type() == lua.LTTable {
		str.(*lua.LTable).RawSetString("dump", lua.LNil)
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
