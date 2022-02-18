package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Type string

type Event interface {
	ID() string
	OccurredOn() time.Time
	Type() Type
	Payload() []byte
	Unmarshall(interface{}) error
}

type BaseEvent struct {
	eventID    string
	_type      Type
	occurredOn time.Time
}

func NewBaseEvent(aggregateID, _type string) Event {
	return BaseEvent{
		eventID:    uuid.New().String(),
		occurredOn: time.Now(),
		_type:      Type(_type),
	}
}

func (b BaseEvent) ID() string {
	return b.eventID
}

func (b BaseEvent) Type() Type {
	return b._type
}

func (b BaseEvent) OccurredOn() time.Time {
	return b.occurredOn
}

func (b BaseEvent) Unmarshall(i interface{}) error {
	return nil
}

func (e BaseEvent) Payload() []byte {
	return []byte{}
}

type EntityEvent struct {
	eventID    string
	_type      Type
	payload    []byte
	occurredOn time.Time
}

func NewEntityEvent(_type string, payload interface{}) (Event, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return EntityEvent{
		eventID:    uuid.New().String(),
		_type:      Type(_type),
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
