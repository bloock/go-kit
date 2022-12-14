package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventType struct {
	name       string
	retry      bool
	expiration int
}

type Event interface {
	ID() string
	OccurredOn() time.Time
	Type() EventType
	Payload() []byte
	Unmarshall(interface{}) error
	Headers() map[string]interface{}
}

type EntityEvent struct {
	eventID    string
	_type      EventType
	payload    []byte
	occurredOn time.Time
}

func NewEntityEvent(_type EventType, payload interface{}) (Event, error) {
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

func (b EntityEvent) Type() EventType {
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

func (t EventType) Name() string {
	return t.name
}

func (t EventType) GetRetryName() string {
	return fmt.Sprintf("%s.retry", t.name)
}

func (t EventType) HasRetry() bool {
	return t.retry
}

func (t EventType) RetryExpiration() int {
	return t.expiration
}
