package observability

import (
	"context"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Tracer struct{}

func InitTracer(app string, debug bool) Tracer {
	tracer.Start(
		tracer.WithService(app),
		tracer.WithDebugMode(debug),
	)

	return Tracer{}
}

func NewSpan(ctx context.Context, name string) (ddtrace.Span, context.Context) {
	return tracer.StartSpanFromContext(ctx, name)
}

func (*Tracer) Stop() {
	tracer.Stop()
}
