package query

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReturnQuery struct {
	ID string
}
type HandlerTest struct {
	err         error
	returnQuery ReturnQuery
}

func (ht HandlerTest) Handle(ctx context.Context, c Query) (interface{}, error) {
	return ht.returnQuery, ht.err
}

func (ht HandlerTest) Type() Type {
	return "test"
}

func TestInmemoryCommand(t *testing.T) {
	t.Run("Given a correct query bus should return query's result", func(t *testing.T) {
		returnQueryID := "id-no-error"
		ht := HandlerTest{
			err:         nil,
			returnQuery: ReturnQuery{ID: returnQueryID},
		}

		queryBus := NewQueryBus()
		queryBus.Register(ht.Type(), ht)

		returnQuery, err := queryBus.Dispatch(context.Background(), ht)
		assert.NoError(t, err)
		assert.Equal(t, returnQueryID, returnQuery.(ReturnQuery).ID)
	})

	t.Run("Given a correct query bus with a custom error", func(t *testing.T) {
		var CustomError = errors.New("Custom error")
		ht := HandlerTest{
			err: CustomError,
		}

		queryBus := NewQueryBus()
		queryBus.Register(ht.Type(), ht)

		_, err := queryBus.Dispatch(context.Background(), ht)
		assert.ErrorIs(t, err, CustomError)
	})

	t.Run("Given and unregister query bus should return error", func(t *testing.T) {
		ht := HandlerTest{}

		queryBus := NewQueryBus()

		result, err := queryBus.Dispatch(context.Background(), ht)
		assert.ErrorIs(t, err, ErrorQueryBus)
		assert.Nil(t, result)
	})
}
