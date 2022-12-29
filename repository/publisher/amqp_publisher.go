package publisher

import (
	"context"

	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/domain"
	"github.com/bloock/go-kit/observability"
)

type AMQPPublisher struct {
	client *client.AMQPClient
	logger observability.Logger
}

func NewAMQPPublisher(client *client.AMQPClient, l observability.Logger) Publisher {
	l.UpdateLogger(l.With().Str("layer", "repository").Str("component", "publisher").Logger())
	return &AMQPPublisher{
		client: client,
		logger: l,
	}
}

func (p AMQPPublisher) Publish(ctx context.Context, event domain.Event, args *PublisherArgs) error {
	if args == nil {
		args = &PublisherArgs{}
	}
	err := p.client.Publish(ctx, event, args.Headers, args.Expiration)
	if err != nil {
		p.logger.Error(ctx).Msgf("error publishing event %s with ID %s: %s", event.Type().Name(), event.ID(), err.Error())
		return err
	}
	p.logger.Info(ctx).Str("type", event.Type().Name()).Str("id", event.ID()).Msgf("published new message with ID %s", event.ID())
	return nil
}
