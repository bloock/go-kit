package publisher

import "github.com/bloock/go-kit/event"

type PublisherArgs struct {
	Expiration int
	Headers    map[string]interface{}
}

type Publisher interface {
	Publish(event event.Event, args *PublisherArgs, retry bool) error
}
