package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(msg string)
	Debugf(format string, args ...any)
	Info(msg string)
	Infof(format string, args ...any)
	Warn(msg string)
	Warnf(format string, args ...any)
	Error(msg string)
	Errorf(format string, args ...any)
	Fatal(msg string)
	Fatalf(format string, args ...any)
}

var _ Logger = &logger{}

const defaultLogLevel = "info"

var (
	defaultWriter   = os.Stdout
	defaultLocation = time.UTC
)

type logger struct {
	appVer         string
	level          string
	skipFrameCount int

	location *time.Location
	writer   io.Writer

	logger *zerolog.Logger
}

func New(appVersion string, opts ...Option) *logger {
	l := &logger{
		appVer:   appVersion,
		level:    defaultLogLevel,
		location: defaultLocation,
		writer:   defaultWriter,
	}

	for _, opt := range opts {
		opt(l)
	}

	zerolog.SetGlobalLevel(l.zerologLevel())

	zerologger := zerolog.
		New(l.writer).
		With().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + l.skipFrameCount).
		Logger()

	l.logger = &zerologger

	return l
}

func (l *logger) zerologLevel() zerolog.Level {
	switch strings.ToLower(l.level) {
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *logger) nowUnix() int64 {
	return time.Now().In(l.location).Unix()
}

func (l *logger) log(e *zerolog.Event) *zerolog.Event {
	return e.Str("ver", l.appVer).Int64("time", l.nowUnix())
}

func (l *logger) Debug(msg string) {
	l.log(l.logger.Debug()).Msg(msg)
}

func (l *logger) Debugf(format string, args ...any) {
	l.log(l.logger.Debug()).Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Info(msg string) {
	l.log(l.logger.Info()).Msg(msg)
}

func (l *logger) Infof(format string, args ...any) {
	l.log(l.logger.Info()).Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Warn(msg string) {
	l.log(l.logger.Warn()).Msg(msg)
}

func (l *logger) Warnf(format string, args ...any) {
	l.log(l.logger.Warn()).Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Error(msg string) {
	l.log(l.logger.Error()).Msg(msg)
}

func (l *logger) Errorf(format string, args ...any) {
	l.log(l.logger.Error()).Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Fatal(msg string) {
	l.log(l.logger.Fatal()).Msg(msg)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.log(l.logger.Fatal()).Msg(fmt.Sprintf(format, args...))
}
