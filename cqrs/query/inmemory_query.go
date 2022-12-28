package query

import (
	"context"
	"errors"

	"github.com/bloock/go-kit/observability"
)

var ErrorQueryBus = errors.New("QueryBus: Query bus handler not found")

type QueryBus struct {
	handlers map[Type]Handler
}

func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[Type]Handler),
	}
}

func (b *QueryBus) Dispatch(ctx context.Context, q Query) (interface{}, error) {
	s, ctx := observability.NewSpan(ctx, string(q.Type()))
	defer s.Finish()

	handler, ok := b.handlers[q.Type()]
	if !ok {
		return nil, ErrorQueryBus
	}

	result, err := handler.Handle(ctx, q)
	return result, err
}

func (b *QueryBus) Register(queryType Type, handler Handler) {
	b.handlers[queryType] = handler
}
