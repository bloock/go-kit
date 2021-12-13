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
	t.Run("given one page completely filled, the pagination values should be correct", func(t *testing.T) {
		p := NewPagination(1, 10, 10)
		assert.Equal(t, int64(1), p.Meta.CurrentPage)
		assert.Equal(t, int64(10), p.Meta.PerPage)
		assert.Equal(t, int64(1), p.Meta.From)
		assert.Equal(t, int64(10), p.Meta.To)
		assert.Equal(t, int64(10), p.Meta.Total)
		assert.Equal(t, int64(1), p.Meta.LastPage)
	})

	t.Run("given one page not filled, the pagination values should be correct", func(t *testing.T) {
		p := NewPagination(1, 10, 8)
		assert.Equal(t, int64(1), p.Meta.CurrentPage)
		assert.Equal(t, int64(10), p.Meta.PerPage)
		assert.Equal(t, int64(1), p.Meta.From)
		assert.Equal(t, int64(8), p.Meta.To)
		assert.Equal(t, int64(8), p.Meta.Total)
		assert.Equal(t, int64(1), p.Meta.LastPage)
	})

	t.Run("given multiple pages and last page is multiple of per page value, the last page value should be correct", func(t *testing.T) {
		p := NewPagination(1, 10, 20)
		assert.Equal(t, int64(1), p.Meta.CurrentPage)
		assert.Equal(t, int64(10), p.Meta.PerPage)
		assert.Equal(t, int64(1), p.Meta.From)
		assert.Equal(t, int64(10), p.Meta.To)
		assert.Equal(t, int64(20), p.Meta.Total)
		assert.Equal(t, int64(2), p.Meta.LastPage)
	})

	t.Run("given multiple pages and last page isn't multiple of per page value, the last page value should be correct", func(t *testing.T) {
		p := NewPagination(2, 10, 25)
		assert.Equal(t, int64(2), p.Meta.CurrentPage)
		assert.Equal(t, int64(10), p.Meta.PerPage)
		assert.Equal(t, int64(11), p.Meta.From)
		assert.Equal(t, int64(20), p.Meta.To)
		assert.Equal(t, int64(25), p.Meta.Total)
		assert.Equal(t, int64(3), p.Meta.LastPage)
	})

	t.Run("given last page requested and not enough records to fill the page, the to value should be correct", func(t *testing.T) {
		p := NewPagination(3, 10, 25)
		assert.Equal(t, int64(3), p.Meta.CurrentPage)
		assert.Equal(t, int64(10), p.Meta.PerPage)
		assert.Equal(t, int64(21), p.Meta.From)
		assert.Equal(t, int64(25), p.Meta.To)
		assert.Equal(t, int64(25), p.Meta.Total)
		assert.Equal(t, int64(3), p.Meta.LastPage)
	})

	t.Run("given a total of 0 rows returned, all values should be correct", func(t *testing.T) {
		p := NewPagination(1, 10, 0)
		assert.Equal(t, int64(1), p.Meta.CurrentPage)
		assert.Equal(t, int64(1), p.Meta.LastPage)
		assert.Equal(t, int64(10), p.Meta.PerPage)
		assert.Equal(t, int64(1), p.Meta.From)
		assert.Equal(t, int64(1), p.Meta.To)
		assert.Equal(t, int64(0), p.Meta.Total)
	})
}

func TestPaginationQuery(t *testing.T) {
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
