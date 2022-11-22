package featureflag

import (
	"github.com/Unleash/unleash-client-go/v3"
	"net/http"
)

type FeatureToggle struct {
}

func InitFeatureToggleClient(application, token string) (*FeatureToggle, error) {
	err := unleash.Initialize(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithAppName(application),
		unleash.WithUrl("https://unleash.bloock.dev/api"),
		unleash.WithCustomHeaders(http.Header{"Authorization": {token}}),
	)

	if err != nil {
		return nil, err
	}

	return &FeatureToggle{}, nil
}

func IsEnabled(featureName string) bool {
	return unleash.IsEnabled(featureName)
}
