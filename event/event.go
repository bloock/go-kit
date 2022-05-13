package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Type struct {
	name       string
	retry      bool
	expiration int
}

type Event interface {
	ID() string
	OccurredOn() time.Time
	Type() Type
	Payload() []byte
	Unmarshall(interface{}) error
	Headers() map[string]interface{}
}

type EntityEvent struct {
	eventID    string
	_type      Type
	payload    []byte
	occurredOn time.Time
}

func NewEntityEvent(_type Type, payload interface{}) (Event, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return EntityEvent{
		eventID:    uuid.New().String(),
		_type:      _type,
		payload:    b,
		occurredOn: time.Now(),
	}, nil
}

func (b EntityEvent) ID() string {
	return b.eventID
}

func (b EntityEvent) OccurredOn() time.Time {
	return b.occurredOn
}

func (b EntityEvent) Type() Type {
	return b._type
}

func (b EntityEvent) Unmarshall(i interface{}) error {
	return json.Unmarshal(b.payload, &i)
}

func (e EntityEvent) Payload() []byte {
	return e.payload
}

func (e EntityEvent) Headers() map[string]interface{} {
	return map[string]interface{}{}
}

func (t Type) Name() string {
	return t.name
}

func (t Type) HasRetry() bool {
	return t.retry
}

func (t Type) RetryExpiration() int {
	return t.expiration
}
