package observability

import (
	"context"
	"os"

	bloockContext "github.com/bloock/go-kit/context"
	"github.com/rs/zerolog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Logger struct {
	underlying zerolog.Logger
}

func InitLogger(env, service, version string, debug bool) Logger {
	var l zerolog.Logger

	if debug {
		l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		l = zerolog.New(os.Stderr)
	}

	l = l.With().Timestamp().Str("dd.env", env).Str("dd.service", service).Str("dd.version", version).Logger()

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
	l.underlying.Printf(format, args...)
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
		event = event.Str("user_id", userID)
	}
	requestID := bloockContext.GetRequestID(ctx)
	if requestID != "" {
		event = event.Str("request_id", requestID)
	}
	trace, ok := tracer.SpanFromContext(ctx)
	if ok {
		event = event.Uint64("dd.trace_id", trace.Context().TraceID())
		event = event.Uint64("dd.span_id", trace.Context().SpanID())
	}
	return event
}
