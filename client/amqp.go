package client

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bloock/go-kit/event"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type AMQPHandler func(context.Context, event.Event) error

var (
	ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
)

type AMQPClient struct {
	ctx context.Context

	consumer       string
	url            string
	reconnectDelay time.Duration
	resendDelay    time.Duration

	logger               zerolog.Logger
	connection           *amqp.Connection
	consumeChannel       *amqp.Channel
	consumeNotifyClose   chan *amqp.Error
	consumeNotifyConfirm chan amqp.Confirmation
	isConnected          bool
	alive                bool
	wg                   *sync.WaitGroup
}

func NewAMQPClient(ctx context.Context, consumer string, user, password, host, port, vhost string, l zerolog.Logger) (*AMQPClient, error) {
	l = l.With().Str("layer", "infrastructure").Str("component", "rabbitmq").Logger()

	addr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost)

	client := AMQPClient{
		ctx:            ctx,
		consumer:       consumer,
		url:            addr,
		reconnectDelay: 5 * time.Second,
		resendDelay:    5 * time.Second,
		logger:         l,
		alive:          true,
		wg:             &sync.WaitGroup{},
	}

	go client.handleReconnect(addr)

	return &client, nil
}

// Consume implements the event.Bus interface.
func (a *AMQPClient) Consume(ctx context.Context, t event.Type, handlers ...AMQPHandler) error {
	a.wg.Add(1)
	for {
		if a.isConnected {
			break
		}
		time.Sleep(a.reconnectDelay)

		select {
		case <-a.ctx.Done():
			return fmt.Errorf("could not connect to client")
		default:
		}
	}

	err := a.consumeChannel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	q, err := a.DeclareQueue(string(t), nil)
	if err != nil {
		return err
	}

	msgs, err := a.consumeChannel.Consume(
		q.Name,
		a.consumerName(string(t)), // Consumer
		false,                     // Auto-Ack
		false,                     // Exclusive
		false,                     // No-local
		false,                     // No-Wait
		nil,                       // Args
	)
	if err != nil {
		return err
	}

	go func() {
		defer a.wg.Done()

		for {
			select {
			case <-a.ctx.Done():
				return
			case err = <-a.consumeNotifyClose:
				a.logger.Warn().Str("type", string(t)).Msgf("consumer channel closed with error %s", err.Error())
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				startTime := time.Now()

				evt := event.NewAMQPEvent(msg, t)

				defer func(e event.Event, m amqp.Delivery) {
					if err := recover(); err != nil {
						stack := make([]byte, 8096)
						stack = stack[:runtime.Stack(stack, false)]
						a.logger.Fatal().Timestamp().Bytes("stack", stack).Interface("error", err).Msg("panic recovery for rabbitMQ message")
						msg.Nack(false, false)
					}
				}(evt, msg)

				err = a.handleMessage(ctx, evt, handlers...)
				if err != nil {
					a.logAndNack(t, msg, startTime, err.Error())
				} else {
					err := msg.Ack(false)
					if err != nil {
						a.logger.Error().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", string(evt.Type())).Str("id", evt.ID()).Msgf("error while sending ack: %s", err.Error())
					} else {
						a.logger.Info().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", string(evt.Type())).Str("id", evt.ID()).Msg("successfully consumed message")
					}
				}
			}
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
	for {
		publishConfirm := make(chan amqp.Confirmation)

		var ch *amqp.Channel
		var err error

		if a.isConnected {
			ch, err = a.connection.Channel()
			if err != nil {
				a.logger.Warn().Str("type", string(event.Type())).Msgf("error while creating publishing channel with error %s", err.Error())
				continue
			}
			ch.Confirm(false)
			ch.NotifyPublish(publishConfirm)

			err = a.UnsafePush(ch, string(event.Type()), event.ID(), event.Payload(), headers, expiration)
			if err != nil {
				a.logger.Warn().Str("type", string(event.Type())).Msgf("error while publishing with error %s", err.Error())
				if err == ErrDisconnected {
					continue
				}
				return err
			}
		}
		select {
		case confirm := <-publishConfirm:
			if confirm.Ack {
				err := ch.Close()
				if err != nil {
					a.logger.Warn().Str("type", string(event.Type())).Str("id", event.ID()).Msgf("error while closing publish channel: %s", err.Error())
				}
				return nil
			}
		case <-a.ctx.Done():
			return fmt.Errorf("could not connect to client")
		case <-time.After(a.resendDelay):
		}
	}
}

type DeclareQueueArgs struct {
	DeadLetterExchange string
}

func (a *AMQPClient) DeclareQueue(name string, args *DeclareQueueArgs) (*amqp.Queue, error) {
	for {
		if a.isConnected {
			break
		}
		time.Sleep(a.reconnectDelay)

		select {
		case <-a.ctx.Done():
			return nil, fmt.Errorf("could not connect to client")
		default:
		}
	}

	if args == nil {
		args = &DeclareQueueArgs{}
	}

	if args.DeadLetterExchange == "" {
		args.DeadLetterExchange = fmt.Sprintf("%s.dead", name)
		err := a.consumeChannel.ExchangeDeclare(
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

		qdlx, err := a.consumeChannel.QueueDeclare(
			fmt.Sprintf("%s.%s.dead", name, a.consumer),
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}

		err = a.consumeChannel.QueueBind(qdlx.Name, fmt.Sprintf("%s.%s", name, a.consumer), args.DeadLetterExchange, false, nil)
		if err != nil {
			return nil, err
		}
	}

	err := a.consumeChannel.ExchangeDeclare(
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

	q, err := a.consumeChannel.QueueDeclare(
		fmt.Sprintf("%s.%s", name, a.consumer),
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    args.DeadLetterExchange,
			"x-dead-letter-routing-key": fmt.Sprintf("%s.%s", name, a.consumer),
		},
	)
	if err != nil {
		return nil, err
	}

	err = a.consumeChannel.QueueBind(q.Name, "", name, false, nil)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (c *AMQPClient) Ping() error {
	if !c.isConnected {
		return fmt.Errorf("client not connected")
	}

	return nil
}

func (c *AMQPClient) NotifyConsumeClose() chan *amqp.Error {
	return c.consumeNotifyClose
}

func (c *AMQPClient) Close() error {
	if !c.isConnected {
		return nil
	}
	c.alive = false

	c.wg.Wait()

	c.isConnected = false
	err := c.connection.Close()
	if err != nil {
		return err
	}
	c.logger.Info().Msg("gracefully stopped rabbitMQ connection")
	return nil
}

// handleReconnect will wait for a connection error on
// notifyClose, and then continuously attempt to reconnect.
func (a *AMQPClient) handleReconnect(addr string) {
	for a.alive {
		a.isConnected = false
		t := time.Now()

		var retryCount int
		for !a.connect(addr) {
			if !a.alive {
				return
			}
			select {
			case <-a.ctx.Done():
				return
			case <-time.After(a.reconnectDelay + time.Duration(retryCount)*time.Second):
				a.logger.Printf("disconnected from rabbitMQ and failed to connect")
				retryCount++
			}
		}
		a.logger.Printf("Connected to rabbitMQ in: %vms", time.Since(t).Milliseconds())
		select {
		case <-a.ctx.Done():
			return
		case <-a.consumeNotifyClose:
		}

	}
}

// connect will make a single attempt to connect to
// RabbitMq. It returns the success of the attempt.
func (a *AMQPClient) connect(addr string) bool {
	a.logger.Info().Msgf("Attempting to connect to rabbitMQ: %s", addr)
	conn, err := amqp.Dial(addr)
	if err != nil {
		a.logger.Printf("failed to dial rabbitMQ server: %v", err)
		return false
	}
	cch, err := conn.Channel()
	if err != nil {
		a.logger.Printf("failed connecting to channel: %v", err)
		return false
	}
	cch.Confirm(false)
	a.changeConnection(conn, cch)
	a.isConnected = true
	return true
}

// changeConnection takes a new connection to the queue,
// and updates the channel listeners to reflect this.
func (c *AMQPClient) changeConnection(connection *amqp.Connection, cch *amqp.Channel) {
	c.connection = connection
	c.consumeChannel = cch
	c.consumeNotifyClose = make(chan *amqp.Error)
	c.consumeNotifyConfirm = make(chan amqp.Confirmation)
	c.consumeChannel.NotifyClose(c.consumeNotifyClose)
	c.consumeChannel.NotifyPublish(c.consumeNotifyConfirm)
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (a *AMQPClient) UnsafePush(ch *amqp.Channel, exchange string, id string, data []byte, headers map[string]interface{}, expiration int) error {
	if !a.isConnected {
		return ErrDisconnected
	}

	exp := ""
	if expiration != 0 {
		exp = fmt.Sprintf("%d", expiration)
	}
	return ch.Publish(
		exchange, // Exchange
		"",       // Routing key
		false,    // Mandatory
		false,    // Immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: id,
			Body:          data,
			Headers:       headers,
			Expiration:    exp,
		},
	)
}

func (a *AMQPClient) logAndNack(_type event.Type, msg amqp.Delivery, t time.Time, err string, args ...interface{}) {
	nErr := msg.Nack(false, false)
	if nErr != nil {
		a.logger.Error().Int64("took-ms", time.Since(t).Milliseconds()).Str("type", string(_type)).Msgf("error while sending nack with error %s: %s", nErr.Error(), fmt.Sprintf(err, args...))
	} else {
		a.logger.Error().Int64("took-ms", time.Since(t).Milliseconds()).Str("type", string(_type)).Msg(fmt.Sprintf(err, args...))
	}
}

func (a *AMQPClient) consumerName(queue string) string {
	return fmt.Sprintf("%s-%s", a.consumer, queue)
}
