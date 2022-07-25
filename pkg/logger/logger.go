package logger

import (
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func levelFromString(l string) zapcore.Level {
	l = strings.ToLower(l)
	switch l {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func New(w io.Writer, level string) *zap.Logger {
	if w == nil {
		panic("writer is nil")
	}
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(w),
		zapcore.Level(levelFromString(level)),
	)
	return zap.New(core)
}
