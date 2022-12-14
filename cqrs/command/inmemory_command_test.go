package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HandlerTest struct {
	err error
}

func (ht HandlerTest) Handle(ctx context.Context, q Command) error {
	return ht.err
}

func (ht HandlerTest) Type() Type {
	return "test"
}

func TestInmemoryCommand(t *testing.T) {

	t.Run("Given a correct command bus should not return err", func(t *testing.T) {
		ht := HandlerTest{}

		commandBus := NewCommandBus()
		commandBus.Register(ht.Type(), ht)

		err := commandBus.Dispatch(context.Background(), ht)
		assert.NoError(t, err)

	})

	t.Run("Given a correct command bus with a custom error", func(t *testing.T) {
		var CustomError = errors.New("custom error")

		ht := HandlerTest{err: CustomError}

		commandBus := NewCommandBus()
		commandBus.Register(ht.Type(), ht)

		err := commandBus.Dispatch(context.Background(), ht)
		assert.ErrorIs(t, err, CustomError)
	})

	t.Run("Given and unregister command bus should return error", func(t *testing.T) {
		ht := HandlerTest{}

		commandBus := NewCommandBus()

		err := commandBus.Dispatch(context.Background(), ht)
		assert.ErrorIs(t, err, ErrorCommandBus)
	})
}
