package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeExternalService struct {
}

func NewFakeExternalService() FakeExternalService {
	return FakeExternalService{}
}

func (f FakeExternalService) HealthCheck() ExternalServiceDetails {
	return ExternalServiceDetails{Error: "someError"}
}

func TestHandler_Check(t *testing.T) {

	health := NewHealth("v0.0.1", "v0.0.1-release", "Error health check", nil, nil, nil)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", health.CheckGinHandler())

	t.Run("handler health returns 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("handler health returns 500 when is not healthy", func(t *testing.T) {
		r := gin.New()
		health = NewHealth("v0.0.1", "v0.0.1-release", "Error health check", nil, nil, []ExternalService{NewFakeExternalService()})
		r.GET("/health", health.CheckGinHandler())
		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
	})
}
