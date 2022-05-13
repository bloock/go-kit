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

type AMQPHandler func(context.Context, event.Event) (Action, error)

var (
	ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
	ErrPendingTransactions = errors.New("there are already pending transactions")
)

type Action = int

const (
	// Ack default ack this msg after you have successfully processed this delivery.
	Ack Action = iota
	// NackDiscard the message will be dropped or delivered to a server configured dead-letter queue.
	NackDiscard
	// NackRequeue deliver this message to a different consumer.
	NackRequeue
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

	q, err := a.DeclareQueue(t, nil)
	if err != nil {
		return err
	}

	go func() {
		defer a.wg.Done()

		err = a.consumer.StartConsuming(
			func(d rabbitmq.Delivery) rabbitmq.Action {
				result := NackDiscard

				startTime := time.Now()

				evt := event.NewAMQPEvent(d, t)

				defer func(e event.Event) {
					if err := recover(); err != nil {
						stack := make([]byte, 8096)
						stack = stack[:runtime.Stack(stack, false)]
						a.logger.Fatal().Timestamp().Bytes("stack", stack).Interface("error", err).Msg("panic recovery for rabbitMQ message")
					}
				}(evt)

				result, err = a.handleMessage(ctx, evt, handlers...)
				if err != nil {
					a.logger.Error().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", t.Name()).Msgf("error while consuming message: %s", err.Error())
				} else if errors.Is(err, ErrPendingTransactions) {
					a.logger.Warn().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", t.Name()).Msgf("wait: %s", err.Error())
				} else {
					a.logger.Info().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", evt.Type().Name()).Str("id", evt.ID()).Msg("successfully consumed message")
				}
				return rabbitmq.Action(result)
			},
			q.Name,
			[]string{""},
			rabbitmq.WithConsumeOptionsConcurrency(1),
			rabbitmq.WithConsumeOptionsQueueNoDeclare,
			rabbitmq.WithConsumeOptionsQOSPrefetch(1),
			rabbitmq.WithConsumeOptionsConsumerName(a.consumerName(t.Name())),
		)
		if err != nil {
			a.logger.Error().Str("type", t.Name()).Msgf("error starting consuming queue: %s", err.Error())
		} else {
			a.consumers = append(a.consumers, a.consumerName(t.Name()))
		}
	}()

	return nil
}

func (a *AMQPClient) handleMessage(ctx context.Context, evt event.Event, handlers ...AMQPHandler) (Action, error) {
	var action Action
	for _, handler := range handlers {
		act, err := handler(ctx, evt)
		if err != nil {
			if evt.Type().HasRetry() {
				if err := a.Publish(evt, evt.Headers(), evt.Type().RetryExpiration()); err != nil {
					return NackRequeue, err
				}
			}
			return act, err
		}
		action = act
	}
	return action, nil
}

// Publish implements the event.Bus interface.
func (a *AMQPClient) Publish(event event.Event, headers map[string]interface{}, expiration int) error {
	exp := ""
	name := event.Type().Name()
	if expiration != 0 {
		exp = fmt.Sprintf("%d", expiration)
	}
	if event.Type().HasRetry() {
		name = event.Type().GetRetryName()
	}

	err := a.publisher.Publish(
		event.Payload(),
		[]string{""},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(name),
		rabbitmq.WithPublishOptionsHeaders(headers),
		rabbitmq.WithPublishOptionsExpiration(exp),
		rabbitmq.WithPublishOptionsCorrelationID(event.ID()),
	)
	if err != nil {
		a.logger.Warn().Str("type", event.Type().Name()).Msgf("error while publishing with error %s", err.Error())
		return err
	}

	return nil
}

type DeclareQueueArgs struct {
	DeadLetterExchange string
}

func (a *AMQPClient) DeclareQueue(t event.Type, args *DeclareQueueArgs) (*amqp.Queue, error) {
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
		args.DeadLetterExchange = fmt.Sprintf("%s.dead", t.Name())
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
			fmt.Sprintf("%s.%s.dead", t.Name(), a.consumerPrefix),
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}

		err = ch.QueueBind(qdlx.Name, fmt.Sprintf("%s.%s", t.Name(), a.consumerPrefix), args.DeadLetterExchange, false, nil)
		if err != nil {
			return nil, err
		}
	}

	if t.HasRetry() == true {
		retryName := t.GetRetryName()
		err = ch.ExchangeDeclare(
			retryName,
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

		qr, err := ch.QueueDeclare(
			fmt.Sprintf("%s.%s.retry", t.Name(), a.consumerPrefix),
			true,
			false,
			false,
			false,
			amqp.Table{
				"x-dead-letter-exchange":    t.Name(),
				"x-dead-letter-routing-key": fmt.Sprintf("%s.%s.retry", t.Name(), a.consumerPrefix),
			},
		)
		if err != nil {
			return nil, err
		}

		err = ch.QueueBind(qr.Name, "", retryName, false, nil)
		if err != nil {
			return nil, err
		}
	}

	err = ch.ExchangeDeclare(
		t.Name(),
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
		fmt.Sprintf("%s.%s", t.Name(), a.consumerPrefix),
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    args.DeadLetterExchange,
			"x-dead-letter-routing-key": fmt.Sprintf("%s.%s", t.Name(), a.consumerPrefix),
		},
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(q.Name, "", t.Name(), false, nil)
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
