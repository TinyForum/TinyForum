package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(level, filename string) error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleSyncer := zapcore.AddSync(os.Stdout)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		consoleSyncer,
		zapLevel,
	)

	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

func Info(msg string, fields ...zap.Field)  { Log.Info(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Log.Error(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { Log.Warn(msg, fields...) }
func Debug(msg string, fields ...zap.Field) { Log.Debug(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Log.Fatal(msg, fields...) }
