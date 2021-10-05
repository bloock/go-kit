package command

import (
	"context"
	"errors"
)

var ErrorCommandBus = errors.New("CommandBus: Command bus handler not found")

type CommandBus struct {
	handlers map[Type]Handler
}

func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[Type]Handler),
	}
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return ErrorCommandBus
	}
	errChannel := make(chan error)
	go func() {
		err := handler.Handle(ctx, cmd)
		if err != nil {
			errChannel <- err
			return
		}
		errChannel <- nil
	}()

	return <-errChannel
}

func (b *CommandBus) Register(cmdType Type, handler Handler) {
	b.handlers[cmdType] = handler
}
