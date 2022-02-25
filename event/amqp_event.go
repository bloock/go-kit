package event

import (
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

type AMQPEvent struct {
	amqp.Delivery
	_type Type
}

func NewAMQPEvent(msg amqp.Delivery, _type Type) Event {
	return AMQPEvent{
		Delivery: msg,
		_type:    _type,
	}
}

func (e AMQPEvent) ID() string {
	return e.Delivery.CorrelationId
}

func (e AMQPEvent) OccurredOn() time.Time {
	return e.Delivery.Timestamp
}

func (e AMQPEvent) Type() Type {
	return e._type
}

func (e AMQPEvent) Payload() []byte {
	return e.Delivery.Body
}

func (e AMQPEvent) Unmarshall(i interface{}) error {
	return json.Unmarshal(e.Delivery.Body, &i)
}
