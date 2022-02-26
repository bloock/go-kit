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
	publishChannel       *amqp.Channel
	consumeChannel       *amqp.Channel
	publishNotifyClose   chan *amqp.Error
	publishNotifyConfirm chan amqp.Confirmation
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
		time.Sleep(1 * time.Second)
	}

	err := a.consumeChannel.Qos(1, 0, false)
	if err != nil {
		return err
	}

	err = a.consumeChannel.ExchangeDeclare(
		string(t),
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q, err := a.consumeChannel.QueueDeclare(
		fmt.Sprintf("%s.%s", string(t), a.consumer),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = a.consumeChannel.QueueBind(q.Name, "", string(t), false, nil)
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
					a.logAndNack(msg, startTime, err.Error())
				} else {
					a.logger.Info().Int64("took-ms", time.Since(startTime).Milliseconds()).Str("type", string(evt.Type())).Str("id", evt.ID()).Msg("successfully consumed message")
					msg.Ack(false)
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
func (a *AMQPClient) Publish(event event.Event) error {
	for {
		if a.isConnected {
			err := a.UnsafePush(string(event.Type()), event.ID(), event.Payload())
			if err != nil {
				if err == ErrDisconnected {
					continue
				}
				return err
			}
		}
		select {
		case confirm := <-a.publishNotifyConfirm:
			if confirm.Ack {
				return nil
			}
		case <-time.After(a.resendDelay):
		}
	}
}

func (c *AMQPClient) Ping() error {
	if !c.isConnected {
		return fmt.Errorf("client not connected")
	}

	return nil
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
		case <-a.publishNotifyClose:
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
	pch, err := conn.Channel()
	if err != nil {
		a.logger.Printf("failed connecting to channel: %v", err)
		return false
	}
	pch.Confirm(false)
	cch, err := conn.Channel()
	if err != nil {
		a.logger.Printf("failed connecting to channel: %v", err)
		return false
	}
	cch.Confirm(false)
	a.changeConnection(conn, pch, cch)
	a.isConnected = true
	return true
}

// changeConnection takes a new connection to the queue,
// and updates the channel listeners to reflect this.
func (c *AMQPClient) changeConnection(connection *amqp.Connection, pch *amqp.Channel, cch *amqp.Channel) {
	c.connection = connection
	c.publishChannel = pch
	c.consumeChannel = cch
	c.publishNotifyClose = make(chan *amqp.Error)
	c.publishNotifyConfirm = make(chan amqp.Confirmation)
	c.consumeNotifyClose = make(chan *amqp.Error)
	c.consumeNotifyConfirm = make(chan amqp.Confirmation)
	c.publishChannel.NotifyClose(c.publishNotifyClose)
	c.publishChannel.NotifyPublish(c.publishNotifyConfirm)
	c.consumeChannel.NotifyClose(c.consumeNotifyClose)
	c.consumeChannel.NotifyPublish(c.consumeNotifyConfirm)
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (a *AMQPClient) UnsafePush(exchange string, id string, data []byte) error {
	if !a.isConnected {
		return ErrDisconnected
	}
	return a.publishChannel.Publish(
		exchange, // Exchange
		"",       // Routing key
		false,    // Mandatory
		false,    // Immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: id,
			Body:          data,
		},
	)
}

func (a *AMQPClient) logAndNack(msg amqp.Delivery, t time.Time, err string, args ...interface{}) {
	msg.Nack(false, false)
	a.logger.Error().Int64("took-ms", time.Since(t).Milliseconds()).Msg(fmt.Sprintf(err, args...))
}

func (a *AMQPClient) consumerName(queue string) string {
	return fmt.Sprintf("%s-%s", a.consumer, queue)
}
