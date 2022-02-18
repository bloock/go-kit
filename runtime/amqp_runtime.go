package runtime

import (
	"context"
	"time"

	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/event"
	"github.com/rs/zerolog"
)

type AMQPRuntime struct {
	client   *client.AMQPClient
	handlers map[string][]client.AMQPHandler

	shutdownTime time.Duration

	logger zerolog.Logger
}

func NewAMQPRuntime(c *client.AMQPClient, shutdownTime time.Duration, l zerolog.Logger) (*AMQPRuntime, error) {
	e := AMQPRuntime{
		client:       c,
		handlers:     make(map[string][]client.AMQPHandler),
		shutdownTime: shutdownTime,
		logger:       l,
	}

	return &e, nil
}

func (e *AMQPRuntime) SetHandlers(h map[string][]client.AMQPHandler) {
	e.handlers = h
}

func (e *AMQPRuntime) Run(ctx context.Context) {
	for t, h := range e.handlers {
		err := e.client.Consume(event.Type(t), h...)
		if err != nil {
			e.logger.Error().Msgf("error consuming type %s: %s", t, err.Error())
			continue
		}

		e.logger.Info().Msgf("starting consuming type %s", t)
	}

	<-ctx.Done()

	if err := e.client.Close(); err != nil {
		e.logger.Info().Msgf("error while closing amqp runtime: %s", err.Error())
	} else {
		e.logger.Info().Msg("amqp runtime closed successfully")
	}
}
