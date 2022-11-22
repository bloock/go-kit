package featureflag

import (
	"github.com/Unleash/unleash-client-go/v3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFeatureToggle(t *testing.T) {
	t.Run("Given application name and token it should initialize ff", func(t *testing.T) {
		unleash.Initialize(
			unleash.WithListener(&unleash.DebugListener{}),
			unleash.WithAppName("default"),
			unleash.WithUrl("https://unleash.bloock.dev/api"),
			unleash.WithCustomHeaders(http.Header{"Authorization": {"default:development.3060027b6591ecbd507d45356a18dd5e6774683ab54503d4756374a5"}}),
		)

		enabled := unleash.IsEnabled("test")

		assert.True(t, enabled)
	})
}
