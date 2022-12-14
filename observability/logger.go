package observability

import (
	"context"
	"os"

	bloockContext "github.com/bloock/go-kit/context"
	"github.com/rs/zerolog"
)

type Logger struct {
	underlying zerolog.Logger
}

func InitLogger(app string, debug bool) Logger {
	var l zerolog.Logger

	if debug {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Str("service", app).Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		l = zerolog.New(os.Stderr).With().Timestamp().Str("service", app).Logger()
	}

	return Logger{
		underlying: l,
	}
}

func (l *Logger) Debug(ctx context.Context) *zerolog.Event {
	e := l.underlying.Debug()
	return populateLogEvent(e, ctx)
}

func (l *Logger) Error(ctx context.Context) *zerolog.Event {
	e := l.underlying.Error()
	return populateLogEvent(e, ctx)
}

func (l *Logger) Fatal(ctx context.Context) *zerolog.Event {
	e := l.underlying.Fatal()
	return populateLogEvent(e, ctx)
}

func (l *Logger) Info(ctx context.Context) *zerolog.Event {
	e := l.underlying.Info()
	return populateLogEvent(e, ctx)
}

func (l *Logger) Log(ctx context.Context) *zerolog.Event {
	e := l.underlying.Log()
	return populateLogEvent(e, ctx)
}

func (l *Logger) Panic(ctx context.Context) *zerolog.Event {
	e := l.underlying.Panic()
	return populateLogEvent(e, ctx)
}

func (l *Logger) Warn(ctx context.Context) *zerolog.Event {
	e := l.underlying.Warn()
	return populateLogEvent(e, ctx)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.underlying.Printf(format, args)
}

func (l *Logger) With() zerolog.Context {
	return l.underlying.With()
}

func (l *Logger) Logger() zerolog.Logger {
	return l.underlying
}

func (l *Logger) UpdateLogger(logger zerolog.Logger) {
	l.underlying = logger
}

func populateLogEvent(l *zerolog.Event, ctx context.Context) *zerolog.Event {
	event := l
	userID := bloockContext.GetUserID(ctx)
	if userID != "" {
		event = event.Str("user-id", userID)
	}
	requestID := bloockContext.GetRequestID(ctx)
	if requestID != "" {
		event = event.Str("request-id", requestID)
	}
	return event
}
