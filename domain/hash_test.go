package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHash(t *testing.T) {
	t.Run("Given hash it should return error if is not valid keccak256", func(t *testing.T) {
		sha256 := "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb65"
		_, err := NewHash(sha256)

		assert.Error(t, err)
	})

	t.Run("Given hash it should not return error when is a valid keccak256", func(t *testing.T) {
		keccak256 := "9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658"
		hash, err := NewHash(keccak256)

		assert.NoError(t, err)
		assert.Equal(t, keccak256, hash.Hash())
	})
}
