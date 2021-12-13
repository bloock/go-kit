package pagination

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	t.Run("pagination query should bind from gin context", func(t *testing.T) {
		queryParams := make(map[string]string)
		queryParams["page"] = "1"
		queryParams["per_page"] = "10"
		c := createTestRequest(queryParams)

		pq, err := NewPaginationQuery(c)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), pq.Page)
		assert.Equal(t, int64(10), pq.PerPage)
	})
}

func createTestRequest(queryParams map[string]string) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := &http.Request{
		URL: &url.URL{},
	}

	q := req.URL.Query()
	for i, s := range queryParams {
		q.Add(i, s)
	}

	req.URL.RawQuery = q.Encode()

	c.Request = req

	return c
}
