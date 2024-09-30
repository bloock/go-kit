package observability

import (
	"context"
	"github.com/bloock/go-kit/config"
	"github.com/getsentry/sentry-go"
)

const (
	RepositoryOperation = "repository"
	HttpServerOperation = "http.server"
)

type Tracer struct{}

type Span struct {
	span *sentry.Span
}

func InitTracer(ctx context.Context, connUrl, env, version string, l Logger) error {
	options := sentry.ClientOptions{
		Dsn:         connUrl,
		Environment: env,
		Release:     version,
	}
	if env == config.ProductionEnvironment {
		options.EnableTracing = true
		options.TracesSampleRate = 1.0
		options.TracesSampler = func(ctx sentry.SamplingContext) float64 {
			if ctx.Span.Op == RepositoryOperation || ctx.Span.Op == HttpServerOperation  {
				return 1.0
			}
			return 0.0
		}
	}

	if err := sentry.Init(options); err != nil {
		l.Error(ctx).Msgf("sentry initialization failed: %v\n", err.Error())
		return err
	}
	return nil
}

func NewRepositorySpan(ctx context.Context, name string) *Span {
	span := sentry.StartSpan(ctx, RepositoryOperation, sentry.WithTransactionName(name))
	return &Span{
		span: span,
	}
}

func (s *Span) Close() {
	s.span.Finish()
}