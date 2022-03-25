package client

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bloock/go-kit/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/wagslane/go-rabbitmq"
)

type AMQPHandler func(context.Context, event.Event) error

var (
	ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
)

type AMQPClient struct {
	ctx context.Context

	consumerPrefix string

	addr           string
	consumers      []string
	consumer       rabbitmq.Consumer
	publisher      *rabbitmq.Publisher
	reconnectDelay time.Duration

	logger zerolog.Logger
	wg     *sync.WaitGroup
}

func NewAMQPClient(ctx context.Context, c string, user, password, host, port, vhost string, l zerolog.Logger) (*AMQPClient, error) {
	l = l.With().Str("layer", "infrastructure").Str("component", "rabbitmq").Logger()

	addr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost)

	connectionCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	for {
		select {
		case <-connectionCtx.Done():
			return &AMQPClient{}, fmt.Errorf("timeout waiting for amqp connection")
		default:
			time.Sleep(1 * time.Second)

			publisher, err := rabbitmq.NewPublisher(
				addr, amqp.Config{},
				rabbitmq.WithPublisherOptionsLogging,
				rabbitmq.WithPublisherOptionsLogger(&l),
			)
			if err != nil {
				l.Warn().Msgf("error while creating publisher with error %s", err.Error())
				continue
			}

			consumer, err := rabbitmq.NewConsumer(
				addr, amqp.Config{},
				rabbitmq.WithConsumerOptionsLogging,
				rabbitmq.WithConsumerOptionsLogger(&l),
			)
			if err != nil {
				l.Warn().Msgf("couldn't create consumer: %s", err.Error())
				continue
			}

			client := AMQPClient{
				ctx:            ctx,
				addr:           addr,
				consumerPrefix: c,
				publisher:      publisher,
				consumer:       consumer,
				reconnectDelay: 5 * time.Second,
				logger:         l,
				wg:             &sync.WaitGroup{},
			}

			return &client, nil
		}
	}
}

// Consume implements the event.Bus interface.
func (a *AMQPClient) Consume(ctx context.Context, t event.Type, handlers ...AMQPHandler) error {
	a.wg.Add(1)

	q, err := a.DeclareQueue(string(t), nil)
	if err != nil {
		return err
	}

	go func() {
		defer a.wg.Done()

		err = a.consumer.StartConsuming(
			func(d rabbitmq.Delivery) rabbitmq.Action {
				result := rabbitmq.NackDiscard

				startTime := time.Now()

				evt := event.NewAMQPEvent(d, t)

				defer func(e event.Event) {
					if err := recover(); err != nil {
						stack := make([]byte, 8096)
						stack = stack[:runtime.Stack(stack, false)]
						a.logger.Fatal().Timestamp().Bytes("stack", stack).Interface("error", err).Msg("panic recovery for rabbitMQ message")
					}
				}(evt)

				err = a.handleMessage(ctx, evt, handlers...)
				if err != nil {
					a.logger.Error().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", string(t)).Msgf("error while consuming message: %s", err.Error())
					result = rabbitmq.NackDiscard
				} else {
					a.logger.Info().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", string(evt.Type())).Str("id", evt.ID()).Msg("successfully consumed message")
					result = rabbitmq.Ack

				}
				return result
			},
			q.Name,
			[]string{""},
			rabbitmq.WithConsumeOptionsConcurrency(1),
			rabbitmq.WithConsumeOptionsQueueNoDeclare,
			rabbitmq.WithConsumeOptionsQOSPrefetch(1),
			rabbitmq.WithConsumeOptionsConsumerName(a.consumerName(string(t))),
		)
		if err != nil {
			a.logger.Error().Str("type", string(t)).Msgf("error starting consuming queue: %s", err.Error())
		} else {
			a.consumers = append(a.consumers, a.consumerName(string(t)))
		}
	}()

	return nil
}

func (a *AMQPClient) handleMessage(ctx context.Context, evt event.Event, handlers ...AMQPHandler) error {
	for _, handler := range handlers {
		err := handler(ctx, evt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Publish implements the event.Bus interface.
func (a *AMQPClient) Publish(event event.Event, headers map[string]interface{}, expiration int) error {
	exp := ""
	if expiration != 0 {
		exp = fmt.Sprintf("%d", expiration)
	}

	err := a.publisher.Publish(
		event.Payload(),
		[]string{""},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(string(event.Type())),
		rabbitmq.WithPublishOptionsHeaders(headers),
		rabbitmq.WithPublishOptionsExpiration(exp),
		rabbitmq.WithPublishOptionsCorrelationID(event.ID()),
	)
	if err != nil {
		a.logger.Warn().Str("type", string(event.Type())).Msgf("error while publishing with error %s", err.Error())
		return err
	}

	return nil
}

type DeclareQueueArgs struct {
	DeadLetterExchange string
}

func (a *AMQPClient) DeclareQueue(name string, args *DeclareQueueArgs) (*amqp.Queue, error) {
	if args == nil {
		args = &DeclareQueueArgs{}
	}

	conn, err := amqp.Dial(a.addr)
	if err != nil {
		a.logger.Error().Msgf("cannot dial: %v: %q", err, a.addr)
		return &amqp.Queue{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		a.logger.Error().Msgf("cannot create channel: %v", err)
		return &amqp.Queue{}, err
	}

	if args.DeadLetterExchange == "" {
		args.DeadLetterExchange = fmt.Sprintf("%s.dead", name)
		err := ch.ExchangeDeclare(
			args.DeadLetterExchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}

		qdlx, err := ch.QueueDeclare(
			fmt.Sprintf("%s.%s.dead", name, a.consumerPrefix),
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}

		err = ch.QueueBind(qdlx.Name, fmt.Sprintf("%s.%s", name, a.consumerPrefix), args.DeadLetterExchange, false, nil)
		if err != nil {
			return nil, err
		}
	}

	err = ch.ExchangeDeclare(
		name,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s.%s", name, a.consumerPrefix),
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    args.DeadLetterExchange,
			"x-dead-letter-routing-key": fmt.Sprintf("%s.%s", name, a.consumerPrefix),
		},
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(q.Name, "", name, false, nil)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (c *AMQPClient) Ping() error {
	return nil
}

func (a *AMQPClient) Close() error {
	a.consumer.Disconnect()

	for _, c := range a.consumers {
		a.consumer.StopConsuming(c, false)
	}

	a.publisher.StopPublishing()

	a.logger.Info().Msg("gracefully stopped rabbitMQ connection")
	return nil
}

func (a *AMQPClient) consumerName(queue string) string {
	return fmt.Sprintf("%s-%s", a.consumerPrefix, queue)
}
