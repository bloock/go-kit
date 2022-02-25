package publisher

import "github.com/bloock/go-kit/event"

type Publisher interface {
	Publish(event event.Event) error
}
