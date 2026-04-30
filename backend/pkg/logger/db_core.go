package logger

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap/zapcore"
)

// dbCore 实现 zapcore.Core 接口，将日志桥接到 SQLite sink
// 设计原则：只做转换，不持有状态；写入通过 globalDB 完成。
type dbCore struct {
	level  zapcore.LevelEnabler
	fields []zapcore.Field // With() 附加的字段
}

func newDBCore(level zapcore.LevelEnabler) *dbCore {
	return &dbCore{level: level}
}

// ── zapcore.Core 接口 ─────────────────────────────────────────

func (c *dbCore) Enabled(lvl zapcore.Level) bool {
	return c.level.Enabled(lvl)
}

func (c *dbCore) With(fields []zapcore.Field) zapcore.Core {
	clone := &dbCore{level: c.level}
	clone.fields = make([]zapcore.Field, len(c.fields)+len(fields))
	copy(clone.fields, c.fields)
	copy(clone.fields[len(c.fields):], fields)
	return clone
}

func (c *dbCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *dbCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	globalDBMu.Lock()
	sink := globalDB
	globalDBMu.Unlock()

	if sink == nil {
		return nil // 未初始化时静默跳过，不影响其他 core
	}

	// 合并 With 字段 + Write 字段
	all := append(c.fields, fields...)
	fieldsJSON := marshalFields(all)

	sink.write(DBEntry{
		Time:       entry.Time,
		Level:      entry.Level.CapitalString(),
		Caller:     fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line),
		// Msg:        entry.Message,
		Stacktrace: entry.Stack,
		Fields:     fieldsJSON,
	})
	return nil
}

func (c *dbCore) Sync() error { return nil }

// ── 字段序列化 ───────────────────────────────────────────────

// marshalFields 将 zap.Field 列表序列化为 JSON 字符串
// 利用 zapcore.NewMapObjectEncoder 避免手动 type-switch
func marshalFields(fields []zapcore.Field) string {
	if len(fields) == 0 {
		return ""
	}
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}
	b, err := json.Marshal(enc.Fields)
	if err != nil {
		return "{}"
	}
	return string(b)
}