package client

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	bloockContext "github.com/bloock/go-kit/context"
	"github.com/bloock/go-kit/domain"
	"github.com/bloock/go-kit/observability"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

type AMQPHandler func(context.Context, domain.Event) error

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

	logger observability.Logger
	wg     *sync.WaitGroup
}

func NewAMQPClient(ctx context.Context, c string, user, password, host, port, vhost string, l observability.Logger) (*AMQPClient, error) {
	l.UpdateLogger(l.With().Str("layer", "infrastructure").Str("component", "rabbitmq").Logger())

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
				l.Warn(connectionCtx).Msgf("error while creating publisher with error %s", err.Error())
				continue
			}

			consumer, err := rabbitmq.NewConsumer(
				addr, amqp.Config{},
				rabbitmq.WithConsumerOptionsLogging,
				rabbitmq.WithConsumerOptionsLogger(&l),
			)
			if err != nil {
				l.Warn(connectionCtx).Msgf("couldn't create consumer: %s", err.Error())
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
func (a *AMQPClient) Consume(ctx context.Context, t domain.EventType, handlers ...AMQPHandler) error {
	a.wg.Add(1)

	q, err := a.DeclareQueue(ctx, t, nil)
	if err != nil {
		return err
	}

	go func() {
		defer a.wg.Done()

		err = a.consumer.StartConsuming(
			func(d rabbitmq.Delivery) rabbitmq.Action {
				result := rabbitmq.Ack

				evt := domain.NewAMQPEvent(d, t)

				err = a.handleMessage(ctx, evt, handlers...)
				if err != nil {
					result = rabbitmq.NackDiscard
				}
				return result
			},
			q.Name,
			[]string{""},
			rabbitmq.WithConsumeOptionsConcurrency(1),
			rabbitmq.WithConsumeOptionsQueueNoDeclare,
			rabbitmq.WithConsumeOptionsQOSPrefetch(1),
			rabbitmq.WithConsumeOptionsConsumerName(a.consumerName(t.Name())),
		)
		if err != nil {
			a.logger.Error(ctx).Str("type", t.Name()).Msgf("error starting consuming queue: %s", err.Error())
		} else {
			a.consumers = append(a.consumers, a.consumerName(t.Name()))
		}
	}()

	return nil
}

func (a *AMQPClient) handleMessage(ctx context.Context, evt domain.Event, handlers ...AMQPHandler) error {
	startTime := time.Now()

	ctx = context.WithValue(ctx, bloockContext.UserIDKey, "")
	ctx = context.WithValue(ctx, bloockContext.RequestIDKey, evt.ID())

	defer func(ctx context.Context, e domain.Event) {
		if err := recover(); err != nil {
			stack := make([]byte, 8096)
			stack = stack[:runtime.Stack(stack, false)]
			a.logger.Fatal(ctx).Timestamp().Bytes("stack", stack).Interface("error", err).Msg("panic recovery for rabbitMQ message")
		}
	}(ctx, evt)

	for _, handler := range handlers {
		err := handler(ctx, evt)
		if err != nil {
			a.logger.Error(ctx).Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", evt.Type().Name()).Msgf("error while consuming message: %s", err.Error())
			return err
		}
	}

	a.logger.Info(ctx).Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", evt.Type().Name()).Str("id", evt.ID()).Msg("successfully consumed message")
	return nil
}

func (a *AMQPClient) Publish(ctx context.Context, event domain.Event, headers map[string]interface{}, expiration int) error {
	return a.publish(ctx, event, event.Type().Name(), headers, expiration)
}

func (a *AMQPClient) PublishRetry(ctx context.Context, event domain.Event, headers map[string]interface{}, expiration int) error {
	return a.publish(ctx, event, event.Type().GetRetryName(), headers, expiration)
}

// publish implements the event.Bus interface.
func (a *AMQPClient) publish(ctx context.Context, event domain.Event, name string, headers map[string]interface{}, expiration int) error {
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
		rabbitmq.WithPublishOptionsExchange(name),
		rabbitmq.WithPublishOptionsHeaders(headers),
		rabbitmq.WithPublishOptionsExpiration(exp),
		rabbitmq.WithPublishOptionsCorrelationID(event.ID()),
		rabbitmq.WithPublishOptionsTimestamp(event.OccurredOn()),
	)
	if err != nil {
		a.logger.Warn(ctx).Str("type", name).Msgf("error while publishing with error %s", err.Error())
		return err
	}

	return nil
}

type DeclareQueueArgs struct {
	DeadLetterExchange string
}

func (a *AMQPClient) DeclareQueue(ctx context.Context, t domain.EventType, args *DeclareQueueArgs) (*amqp.Queue, error) {
	if args == nil {
		args = &DeclareQueueArgs{}
	}

	conn, err := amqp.Dial(a.addr)
	if err != nil {
		a.logger.Error(ctx).Msgf("cannot dial: %v: %q", err, a.addr)
		return &amqp.Queue{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		a.logger.Error(ctx).Msgf("cannot create channel: %v", err)
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

	if t.HasRetry() {
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

func (a *AMQPClient) Close(ctx context.Context) error {
	a.consumer.Disconnect()

	for _, c := range a.consumers {
		a.consumer.StopConsuming(c, false)
	}

	a.publisher.StopPublishing()

	a.logger.Info(ctx).Msg("gracefully stopped rabbitMQ connection")
	return nil
}

func (a *AMQPClient) consumerName(queue string) string {
	return fmt.Sprintf("%s-%s", a.consumerPrefix, queue)
}
