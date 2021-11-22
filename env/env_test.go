package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {

	os.Setenv("KIT_ENV", "ok")

	t.Run("read env should be ok", func(t *testing.T) {
		var cfg struct{ Env string }
		err := ReadEnv("KIT", &cfg)

		assert.NoError(t, err)
		assert.Equal(t, "ok", cfg.Env)
	})
}
