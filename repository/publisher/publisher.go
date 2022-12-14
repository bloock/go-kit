package publisher

import (
	"context"

	"github.com/bloock/go-kit/domain"
)

type PublisherArgs struct {
	Expiration int
	Headers    map[string]interface{}
}

type Publisher interface {
	Publish(ctx context.Context, event domain.Event, args *PublisherArgs) error
}
