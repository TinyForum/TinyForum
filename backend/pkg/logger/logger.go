package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"tiny-forum/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

// Config 日志配置
type Config struct {
	Level      string // debug, info, warn, error, fatal
	Filename   string // 日志文件路径
	MaxSize    int    // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int    // 最多保留多少个备份
	MaxAge     int    // 文件最多保存多少天
	Compress   bool   // 是否压缩
	Console    bool   // 是否输出到控制台
	JSONFormat bool   // 是否使用 JSON 格式

	// DB 可选。非零值时自动初始化 SQLite 日志数据库并接入日志管道。
	// 若已单独调用过 InitDB，此字段可留空（不会重复初始化）。
	DB *config.DBConfig
}

// Init 初始化日志
func Init(cfg Config) error {
	// 解析日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// 创建 cores
	var cores []zapcore.Core

	// 控制台输出
	if cfg.Console {
		consoleEncoder := getConsoleEncoder(cfg.JSONFormat)
		consoleSyncer := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleSyncer, level))
	}

	// 文件输出
	if cfg.Filename != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %w", err)
		}

		// 配置日志轮转
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		}

		fileEncoder := getFileEncoder(cfg.JSONFormat)
		fileSyncer := zapcore.AddSync(fileWriter)
		cores = append(cores, zapcore.NewCore(fileEncoder, fileSyncer, level))
	}

	// 如果没有配置任何输出，默认输出到控制台
	if len(cores) == 0 {
		consoleEncoder := getConsoleEncoder(false)
		consoleSyncer := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleSyncer, level))
	}

	// ── SQLite 数据库输出（可选，非侵入式插入）─────────────────
	if cfg.DB != nil {
		if err := InitDB(cfg.DB); err != nil {
			return fmt.Errorf("初始化 SQLite 日志失败: %w", err)
		}
		cores = append(cores, newDBCore(level))
	} else {
		// 已通过 InitDB 单独初始化时，同样接入 core
		globalDBMu.Lock()
		alreadyInit := globalDB != nil
		globalDBMu.Unlock()
		if alreadyInit {
			cores = append(cores, newDBCore(level))
		}
	}
	// ──────────────────────────────────────────────────────────

	// 合并多个 cores
	core := zapcore.NewTee(cores...)

	// 创建 logger
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	Sugar = Log.Sugar()

	return nil
}

// getConsoleEncoder 获取控制台编码器
func getConsoleEncoder(jsonFormat bool) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if jsonFormat {
		return zapcore.NewJSONEncoder(encoderCfg)
	}
	return zapcore.NewConsoleEncoder(encoderCfg)
}

// getFileEncoder 获取文件编码器
func getFileEncoder(jsonFormat bool) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if jsonFormat {
		return zapcore.NewJSONEncoder(encoderCfg)
	}
	return zapcore.NewConsoleEncoder(encoderCfg)
}

// ========== 格式化日志方法（支持参数拼接）==========

// Infof 格式化 Info 日志
func Infof(format string, args ...interface{}) {
	Sugar.Infof(format, args...)
}

// Errorf 格式化 Error 日志
func Errorf(format string, args ...interface{}) {
	Sugar.Errorf(format, args...)
}

// Warnf 格式化 Warn 日志
func Warnf(format string, args ...interface{}) {
	Sugar.Warnf(format, args...)
}

// Debugf 格式化 Debug 日志
func Debugf(format string, args ...interface{}) {
	Sugar.Debugf(format, args...)
}

// Fatalf 格式化 Fatal 日志
func Fatalf(format string, args ...interface{}) {
	Sugar.Fatalf(format, args...)
}

// ========== 结构化日志方法（支持键值对）==========

// Info 结构化 Info 日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Error 结构化 Error 日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Warn 结构化 Warn 日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Debug 结构化 Debug 日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Fatal 结构化 Fatal 日志
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// ========== 键值对日志方法（更简洁）==========

// InfoKV 使用键值对记录 Info 日志
func InfoKV(msg string, keysAndValues ...interface{}) {
	Sugar.Infow(msg, keysAndValues...)
}

// ErrorKV 使用键值对记录 Error 日志
func ErrorKV(msg string, keysAndValues ...interface{}) {
	Sugar.Errorw(msg, keysAndValues...)
}

// WarnKV 使用键值对记录 Warn 日志
func WarnKV(msg string, keysAndValues ...interface{}) {
	Sugar.Warnw(msg, keysAndValues...)
}

// DebugKV 使用键值对记录 Debug 日志
func DebugKV(msg string, keysAndValues ...interface{}) {
	Sugar.Debugw(msg, keysAndValues...)
}

// ========== 辅助方法 ==========

// With 创建带字段的子 logger
func With(fields ...zap.Field) *zap.Logger {
	return Log.With(fields...)
}

// WithOptions 创建带选项的子 logger
func WithOptions(opts ...zap.Option) *zap.Logger {
	return Log.WithOptions(opts...)
}

// Sync 刷新缓冲区（包含 SQLite 缓冲区）
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}

// ========== 字段辅助函数（避免命名冲突）==========

// Str 创建字符串字段
func Str(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int 创建整数字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Int64 创建 int64 字段
func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

// Uint 创建 uint 字段
func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

// Bool 创建 bool 字段
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Err 创建 error 字段（避免与 Error 函数冲突）
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Any 创建任意类型字段
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// Dur 创建时间间隔字段
func Dur(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

// Time 创建时间字段
func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}