package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID_validate(t *testing.T) {
	t.Run("Given a valid ID, should return no err", func(t *testing.T) {
		id := "8a8579ab-61a1-433f-9fc8-e511c8d7a4d0"
		err := validateID(id)
		assert.NoError(t, err)
	})

	t.Run("Given an invalid ID (illegal characters), should return ErrIdInvalid", func(t *testing.T) {
		id := "Ma8579ab-61a1-433f-9fc8-e511c8d7a4dM"
		err := validateID(id)
		assert.Equal(t, ErrIDInvalid, err)
	})

	t.Run("Given an invalid ID (illegal length), should return ErrIdInvalid", func(t *testing.T) {
		id := "8a8579ab-61a1-433f-9fc8-e511c8d7a4d02"
		err := validateID(id)
		assert.Equal(t, ErrIDInvalid, err)
	})

	t.Run("Given an empty ID, should return ErrIdEmpty", func(t *testing.T) {
		id := ""
		err := validateID(id)
		assert.Equal(t, ErrIDEmpty, err)
	})

	t.Run("Given a null ID, should return ErrIdNull", func(t *testing.T) {
		id := "00000000-0000-0000-0000-000000000000"
		err := validateID(id)
		assert.Equal(t, ErrIDNull, err)
	})
}

func TestUUID_GenUUID(t *testing.T) {
	t.Run("Given a call, should work", func(t *testing.T) {
		id := GenUUID()
		assert.NotEqual(t, "", id)
	})
}
