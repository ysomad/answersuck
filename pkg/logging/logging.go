package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type logger struct {
	zerolog *zerolog.Logger
}

var _ Logger = (*logger)(nil)

func NewLogger(level string) *logger {
	var zl zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		zl = zerolog.ErrorLevel
	case "warn":
		zl = zerolog.WarnLevel
	case "info":
		zl = zerolog.InfoLevel
	case "debug":
		zl = zerolog.DebugLevel
	default:
		zl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(zl)
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger{
		zerolog: &l,
	}
}

func (l *logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *logger) Error(message interface{}, args ...interface{}) {
	if l.zerolog.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	} else {
		l.msg("error", message, args...)
	}
}

func (l *logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.zerolog.Info().Msg(message)
	} else {
		l.zerolog.Info().Msgf(message, args...)
	}
}

func (l *logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
