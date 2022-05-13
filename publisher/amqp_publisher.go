package publisher

import (
	"github.com/bloock/go-kit/client"
	"github.com/bloock/go-kit/event"
	"github.com/rs/zerolog"
)

type AMQPPublisher struct {
	client *client.AMQPClient
	logger zerolog.Logger
}

func NewAMQPPublisher(client *client.AMQPClient, l zerolog.Logger) Publisher {
	l = l.With().Str("layer", "repository").Str("service", "publisher").Logger()
	return &AMQPPublisher{
		client: client,
		logger: l,
	}
}

func (p AMQPPublisher) Publish(event event.Event, args *PublisherArgs, retry bool) error {
	if args == nil {
		args = &PublisherArgs{}
	}
	err := p.client.Publish(event, args.Headers, args.Expiration, retry)
	if err != nil {
		p.logger.Error().Msgf("error publishing event %s with ID %s: %s", event.Type().Name(), event.ID(), err.Error())
		return err
	}
	p.logger.Info().Str("type", event.Type().Name()).Str("id", event.ID()).Msgf("published new message with ID %s", event.ID())
	return nil
}
