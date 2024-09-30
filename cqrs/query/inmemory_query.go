package query

import (
	"context"
	"errors"
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
