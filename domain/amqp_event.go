package domain

import (
	"encoding/json"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

type AMQPEvent struct {
	rabbitmq.Delivery
	_type EventType
}

func NewAMQPEvent(msg rabbitmq.Delivery, _type EventType) Event {
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

func (e AMQPEvent) Type() EventType {
	return e._type
}

func (e AMQPEvent) Payload() []byte {
	return e.Delivery.Body
}

func (e AMQPEvent) Unmarshall(i interface{}) error {
	return json.Unmarshal(e.Delivery.Body, &i)
}

func (e AMQPEvent) Headers() map[string]interface{} {
	return e.Delivery.Headers
}
