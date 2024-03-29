package runtime

import (
	"context"
	"time"

	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/domain"
	"github.com/bloock/go-kit/observability"
)

type AMQPRuntime struct {
	client   *client.AMQPClient
	handlers map[domain.EventType][]client.AMQPHandler

	shutdownTime time.Duration

	logger observability.Logger
}

func NewAMQPRuntime(c *client.AMQPClient, shutdownTime time.Duration, l observability.Logger) (*AMQPRuntime, error) {
	e := AMQPRuntime{
		client:       c,
		handlers:     make(map[domain.EventType][]client.AMQPHandler),
		shutdownTime: shutdownTime,
		logger:       l,
	}

	return &e, nil
}

func (e *AMQPRuntime) SetHandlers(h map[domain.EventType][]client.AMQPHandler) {
	e.handlers = h
}

func (e *AMQPRuntime) Run(ctx context.Context) {
out:
	for {
		for t, h := range e.handlers {
			err := e.client.Consume(ctx, t, h...)
			if err != nil {
				e.logger.Error(ctx).Msgf("error consuming type %s: %s", t.Name(), err.Error())
				continue
			}

			e.logger.Info(ctx).Msgf("starting consuming type %s", t.Name())
		}

		select {
		case <-ctx.Done():
			break out
		}
	}

	if err := e.client.Close(ctx); err != nil {
		e.logger.Info(ctx).Msgf("error while closing amqp runtime: %s", err.Error())
	} else {
		e.logger.Info(ctx).Msg("amqp runtime closed successfully")
	}
}
