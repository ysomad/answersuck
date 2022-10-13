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

func WithLocation(zoneinfo string) Option {
	return func(l *logger) {
		loc, err := time.LoadLocation(zoneinfo)
		if err != nil {
			l.location = defaultLocation
		}

		l.location = loc
	}
}

// WithMoscowLocation is helper function for WithLocation.
func WithMoscowLocation() Option { return WithLocation("Europe/Moscow") }

func WithSkipFrameCount(c int) Option {
	return func(l *logger) {
		l.skipFrameCount = c
	}
}
