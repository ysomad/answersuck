package logger

import (
	"io"
	"time"
)

type Option func(*logger)

func WithLevel(level string) Option {
	return func(l *logger) {
		if level == "" {
			l.level = defaultLogLevel
			return
		}

		l.level = level
	}
}

func WithWriter(w io.Writer) Option {
	return func(l *logger) {
		if w == nil {
			l.writer = defaultWriter
			return
		}

		l.writer = w
	}
}

func WithLocation(loc *time.Location) Option {
	return func(l *logger) {
		if loc == nil {
			l.location = defaultLocation
			return
		}

		l.location = loc
	}
}

func WithSkipFrameCount(c int) Option {
	return func(l *logger) {
		l.skipFrameCount = c
	}
}
